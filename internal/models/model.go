package models

import "time"

type Order_client struct {
	Order_uid          string    `json:"order_uid"`
	Track_number       string    `json:"track_number"`
	Entry              string    `json:"entry" fake:"{lettern: 5}"`
	Delivery           Delivery  `json:"delivery"`
	Payment            Payment   `json:"payment"`
	Items              []Item    `json:"items"`
	Locale             string    `json:"locale"`
	Internal_signature string    `json:"internal_signature"`
	Customer_id        string    `json:"customer_id"`
	Delivery_service   string    `json:"delivery_service"`
	Shardkey           string    `json:"shardkey" fake:"{digitn: 5}"`
	Sm_id              int       `json:"sm_id"`
	Date_created       time.Time `json:"date_created"`
	Oof_shard          string    `json:"oof_shard" fake:"{digitn: 5}"`
}

type Delivery struct {
	Name    string `json:"name" fake:"{name}"`
	Phone   string `json:"phone" fake:"{phone}"`
	Zip     string `json:"zip" fake:"{digitn: 5}"`
	City    string `json:"city" fake:"{city}"`
	Address string `json:"address" fake:"{streetname}"`
	Region  string `json:"region" fake:"{state}"`
	Email   string `json:"email" fake:"{email}"`
}

type Payment struct {
	Transaction   string `json:"transaction"`
	Request_id    string `json:"request_id"`
	Currency      string `json:"currency" fake:"{currencyshort}"`
	Provider      string `json:"provider"`
	Amount        int    `json:"amount" fake:"{number: 1, 100}"`
	Payment_dt    int    `json:"payment_dt" fake:"{number: 1, 10000}"`
	Bank          string `json:"bank" fake:"{company}"`
	Delivery_cost int    `json:"delivery_cost" fake:"{number: 1, 10000}"`
	Goods_total   int    `json:"goods_total" fake:"{number: 1, 100}"`
	Custom_fee    int    `json:"custom_fee" fake:"{number: 1, 100}"`
}

type Item struct {
	Chrt_id      int    `json:"chrt_Id" fake:"{number: 1, 10000000}"`
	Track_number string `json:"track_number" `
	Price        int    `json:"price" fake:"{number: 1, 100000}"`
	Rid          string `json:"rid"`
	Name         string `json:"name" fake:"{nouncountable}"`
	Sale         int    `json:"sale" fake:"{number: 1, 100}"`
	Size         string `json:"size" fake:"{digit}"`
	Total_price  int    `json:"total_price" fake:"{number: 1, 10000000}"`
	Nm_id        int    `json:"nm_id" fake:"{number: 1, 100000}"`
	Brand        string `json:"brand" fake:"{company}"`
	Status       int    `json:"status" fake:"{number: 1, 500}"`
}
