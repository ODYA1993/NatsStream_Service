package main

import (
	"context"
	"encoding/json"
	"fmt"
	cs "github.com/DmitryOdintsov/Level_0/internal/cashe"
	"github.com/DmitryOdintsov/Level_0/internal/config"
	"github.com/DmitryOdintsov/Level_0/internal/models"
	postgres "github.com/DmitryOdintsov/Level_0/internal/storage"
	stan "github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cfg := config.MustLoad("SUBSCRIBER") // загрузка конфига subscribera

	// подключение к базе данных
	storagePath := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DataBase.Host, cfg.DataBase.Port, cfg.DataBase.User,
		cfg.DataBase.Password, cfg.DataBase.Dbname, cfg.DataBase.Sslmode,
	)
	storage, err := postgres.NewDB(storagePath)
	if err != nil {
		log.Fatalf("Subscriber: %s", err)
	}

	// создание кэша и его загрузка из бд
	cashe := cs.NewCashe()
	err = storage.UploadCashe(&cashe)
	if err != nil {
		log.Fatalf("Subscriber: %s", err)
	}

	dataRecieved := *new(models.Order_client)

	// подключение к Nats-Streaming
	sc, err := stan.Connect(cfg.NatsConfig.ClusterID, cfg.NatsConfig.ClientID)
	if err != nil {
		log.Fatalf("Subscriber: %s", err)
	}

	sub, err := sc.Subscribe("JsonPipe", func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &dataRecieved)
		if err != nil {
			log.Println(err)
		} else {
			// добавление в кэш
			err = cashe.InsertToCashe(dataRecieved)
			if err != nil {
				log.Println(err)
			} else {
				// добавление в бд
				err = storage.SaveOrder(dataRecieved)
				if err != nil {
					log.Println(err)
				}
			}
		}
	},

		stan.DeliverAllAvailable())

	// хэндлер для отправки json на сайт
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl, err := template.ParseFiles("template/ui.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case "POST":
			if order, err := cashe.GetFromCashe(r.PostFormValue("order_uid")); err == nil {
				jsonToSend, err := json.MarshalIndent(order, "", " ")
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Отправка заказа с id: %s\n", r.PostFormValue("order_uid"))
					_, _ = fmt.Fprintf(w, string(jsonToSend))
				}
			} else {
				_, _ = fmt.Fprint(w, "Error: ", err)
			}
		}
	})
	server := &http.Server{Addr: cfg.HTTPServer.Address}
	go server.ListenAndServe()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	<-sigch

	_ = sub.Unsubscribe()
	_ = sc.Close()
	_ = server.Shutdown(context.Background())
}
