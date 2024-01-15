package cashe

import (
	"errors"
	"github.com/DmitryOdintsov/Level_0/internal/models"
)

type Cashe struct {
	Data map[string]models.Order_client
}

func NewCashe() Cashe {
	return Cashe{Data: make(map[string]models.Order_client)}
}

func (c *Cashe) InsertToCashe(jsonToInsert models.Order_client) error {
	_, ok := c.Data[jsonToInsert.Order_uid]
	if ok {
		return errors.New("уже существует в кэш")
	}
	c.Data[jsonToInsert.Order_uid] = jsonToInsert
	return nil
}

func (c *Cashe) GetFromCashe(order_uid string) (models.Order_client, error) {
	order, ok := c.Data[order_uid]
	if !ok {
		return models.Order_client{}, errors.New("нет такого order_uid в кэше")
	}
	return order, nil
}
