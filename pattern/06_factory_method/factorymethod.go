package main

import "fmt"

// Factory method pattern позволяет скрыть логику создания создаваемых экземпляров.
// Клиент взаимодействует только с фабричной структурой и сообщает 
// тип экземпляров, которые необходимо создать. 
// Фабричный класс взаимодействует с соответствующими конкретными 
// структурами и возвращает правильный экземпляр обратно.

// интерфейс iDevice, который определяет все методы, которые должны быть у устройства
type iDevice interface {
	setName(name string)
	setSize(power int)
	getName() string
	getSize() int
}

//  структура device, реализующая интерфейс iDevice.
type device struct {
	name string
	size int
}

func (d *device) setName(name string) {
	d.name = name
}

func (d *device) getName() string {
	return d.name
}

func (d *device) setSize(power int) {
	d.size = power
}

func (d *device) getSize() int {
	return d.size
}

// Phone встраивает структуру device и, следовательно, также косвенно 
// реализуют все методы iDevice и, следовательно, относятся к типу iDevice
type phone struct {
	device
}

func newPhone() iDevice {
	return &phone{
		device: device{
			name: "phone",
			size: 10,
		},
	}
}

// Computer встраивает структуру device и, следовательно, также косвенно 
// реализуют все методы iDevice и, следовательно, относятся к типу iDevice
type computer struct {
	device
}

func newComputer() iDevice {
	return &computer{
		device: device{
			name: "computer",
			size: 100,
		},
	}
}

// Factory, которая создает устройство типа phone или computer.
func getDevice(deviceType string) (iDevice) {
	if deviceType == "phone" {
		return newPhone()
	}
	if deviceType == "computer" {
		return newComputer()
	}
	return nil
}

func main(){
	phone:= getDevice("phone")
	fmt.Println("device:", phone.getName(), "size:", phone.getSize())

	computer:=getDevice("computer")
	fmt.Println("device:", computer.getName(), "size:", computer.getSize())
}