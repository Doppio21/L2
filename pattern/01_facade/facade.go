package main

import "fmt"

// Шаблон Facade Pattern — это структурный шаблон, который 
// предоставляет упрощенный интерфейс для большего объема кода. 
// Он используется для того, чтобы сделать сложную систему более доступной, 
// простой в использовании и более понятной для конечных пользователей.

type Market struct {
}

func (m *Market) SendOrder() {
	fmt.Println("market send the order")
}

type Client struct {
}

func (c *Client) GetOrder() {
	fmt.Println("client get the order")
}

type Delivery struct {
}

func (d *Delivery) DeliveredOrder() {
	fmt.Println("order is delivered")
}

// Интерфейс более высокого уровня, который упрощает использование подсистемы.
type OrderFacade struct{
	market *Market
	client *Client
	delivery *Delivery
}

func NewOrderFacade() *OrderFacade{
	return &OrderFacade{
		market: &Market{},
		client: &Client{},
		delivery: &Delivery{},
	}
}

func (o *OrderFacade) Start(){
	o.market.SendOrder()
	o.delivery.DeliveredOrder()
	o.client.GetOrder()
}


func main() {
	order := NewOrderFacade()
	order.Start()
}
