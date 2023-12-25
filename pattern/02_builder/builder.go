package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Builder Pattern  позволяет разбить строительство сложного объекта,
// кроме того, он также позволяет скрыть от пользователя процесс строительства.

// Плюсы использования:
// Позволяет варьировать внутреннее представление продукта.
// Инкапсулирует код для построения и представления.
// Обеспечивает контроль над этапами строительного процесса.
// Недостатки:
// Для каждого типа продукта должен быть создан отдельный ConcreteBuilder.
// Классы построителей должны быть изменяемыми.
// Может затруднить/усложнить внедрение зависимостей.


type Data struct {
	Text   []byte
	Format string
}

type DataBuilder interface {
	SetRecipient(recipient string)
	SetText(text string)
	Data() (*Data, error)
}

// Конкретная реализация интерфейса. Используется пакет json для кодировки данных.
type JSONDataBuilder struct {
	dataRecipient string
	dataText      string
}

func (b *JSONDataBuilder) SetRecipient(recipient string) {
	b.dataRecipient = recipient
}

func (b *JSONDataBuilder) SetText(text string) {
	b.dataText = text
}

func (b *JSONDataBuilder) Data() (*Data, error) {
	m := make(map[string]string)
	m["recipient"] = b.dataRecipient
	m["message"] = b.dataText

	text, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &Data{Text: text, Format: "JSON"}, nil
}

// Реализация интерфейса с использованием пакета XML для кодировки данных.
type XMLDataBuilder struct {
	dataRecipient string
	dataText      string
}

func (b *XMLDataBuilder) SetRecipient(recipient string) {
	b.dataRecipient = recipient
}

func (b *XMLDataBuilder) SetText(text string) {
	b.dataText = text
}

func (b *XMLDataBuilder) Data() (*Data, error) {
	type XMLMessage struct {
		Recipient string `xml:"recipient"`
		Text      string `xml:"body"`
	}

	m := XMLMessage{
		Recipient: b.dataRecipient,
		Text:      b.dataText,
	}

	text, err := xml.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &Data{Text: text, Format: "XML"}, nil
}

type Publisher struct{}

func (p *Publisher) BuildMessage(builder DataBuilder) (*Data, error) {
	builder.SetRecipient("Publisher")
	builder.SetText("Text")
	return builder.Data()
}

// С помощью спроектированной архитектуры можно создавать сообщения в XML или JSON формате

func main() {
	publisher := Publisher{}
	jsonMsg, err := publisher.BuildMessage(&JSONDataBuilder{})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonMsg.Text))

	xmlMsg, err := publisher.BuildMessage(&XMLDataBuilder{})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(xmlMsg.Text))
}
