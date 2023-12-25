package main

import "fmt"

// Chain of Responsibility Design Pattern это поведенческий шаблон проектирования. Он позволяет создать цепочку 
// обработчиков запросов. Для каждого входящего запроса он проходит через цепочку и каждый из обработчиков:
// Обрабатывает запрос или пропускает обработку.
// Решает, передавать ли запрос следующему обработчику в цепочке или нет

type orderDelivery interface {
	execute(*order)
	setNext(orderDelivery)
}

type market struct {
	next orderDelivery
}

func (m *market) execute(o *order) {
	if o.sendingDone {
		fmt.Println("market sent the order")
		m.next.execute(o)
		return
	}
	fmt.Println("market sends the order")
	o.sendingDone = true
	m.next.execute(o)
}

func (m *market) setNext(next orderDelivery) {
	m.next = next
}

type delivery struct {
	next orderDelivery
}

func (d *delivery) execute(o *order) {
	if o.deliveryDone {
		fmt.Println("order delivered")
		d.next.execute(o)
		return
	}
	fmt.Println("order is delivered")
	o.deliveryDone = true
	d.next.execute(o)
}

func (d *delivery) setNext(next orderDelivery) {
	d.next = next
}

type client struct {
	next orderDelivery
}

func (c *client) execute(o *order) {
	if o.getDone {
		fmt.Println("client got the order")
		
	}
	fmt.Println("client get the order")
}

func (c *client) setNext(next orderDelivery) {
	c.next = next
}

type order struct {
	sendingDone  bool
	deliveryDone bool
	getDone      bool
}

func main(){
	client:= &client{}

	// Set next for delivery
	delivery:= &delivery{}
	delivery.setNext(client)

	// Set next for market
	market:=&market{}
	market.setNext(delivery)

	order:=&order{}
	market.execute(order)
}
