package main

import (
	"encoding/json"
	"github.com/DmitryOdintsov/Level_0/internal/config"
	"github.com/DmitryOdintsov/Level_0/internal/models"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.MustLoad("PRODUCER") // загрузка конфига для producerа
	var order models.Order_client      // структура заказа

	// подключение к Nats-Streaming
	sc, err := stan.Connect(cfg.NatsConfig.ClusterID, cfg.NatsConfig.ClientID)
	if err != nil {
		log.Fatalf("Unable to connect %s", err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 2)

			// Создание новых заказов
			err = gofakeit.Struct(&order)
			if err != nil {
				log.Printf("Unable to generate json due to %s", err)
			}

			// Order_client -> []byte
			jsonToSend, err := json.MarshalIndent(order, "", " ")
			if err != nil {
				log.Printf("Unable to marshal JSON due to %s", err)
			}

			// Публикация в канал
			err = sc.Publish("JsonPipe", jsonToSend)
			if err != nil {
				log.Printf("Json wasn't published due to: %s", err)
			} else {
				log.Printf("Успешно отправлено: %s", order.Order_uid)
			}
		}
	}()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	<-sigch

	_ = sc.Close()
}
