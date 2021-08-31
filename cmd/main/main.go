package main

import (
	"fmt"
	"time"

	"github.com/ark-go/arkEvents/pkg/events"
)

func main() {
	test := events.NewWatch("test")
	// структура для передачи данных в событиях
	type PayloadCustom struct {
		ww string
	}
	//
	go func() {
		enter := test.AddListener("enter")
		for {
			msg := <-enter
			//! Желательно использовать проверку приведения типа !
			if val, ok := msg.([]string); ok {
				fmt.Println("enter 1:", val)
			}
		}
	}()
	go func() {
		enter := test.AddListener("enter")
		for {
			msg := <-enter
			val := msg.([]string)
			fmt.Println("enter 2:", val)
		}
	}()
	go func() {
		enter4 := test.AddListener("move")
		for {
			msg := <-enter4
			fmt.Println("move 1:", msg.([]string), len(msg.([]string)))
		}
	}()
	go func() {
		roll := test.AddListener("roll")
		for {
			msg := <-roll
			if val, ok := msg.([]string); ok {
				fmt.Println("roll 1:", val)
			}
		}
	}()
	go func() {
		enter := test.AddListener("enter")
		for {
			test.RemoveListener("enter", enter) // удаляем Listener
			msg := <-enter
			if val, ok := msg.([]string); ok {
				fmt.Println("Event3:", val)
			}
		}
	}()
	go func() {
		enter := test.AddListener("scroll")
		for {
			msg := <-enter
			fmt.Println("scroll 1:", msg.(string))
		}
	}()
	go func() {
		enter := test.AddListener("size")
		for {
			msg := <-enter
			if val, ok := msg.(*PayloadCustom); ok {
				fmt.Println("size 1:", val.ww)
			}
		}
	}()
	go func() {
		enter := test.AddListener("scroll")
		for {
			msg := <-enter
			fmt.Println("scroll 2:", msg.(string))
		}
	}()
	time.Sleep(3 * time.Second)
	test.Emit("scroll", "77 scroll 77")
	test.Emit("enter", []string{"11111", "111"})
	test.Emit("roll", []string{"222222"})
	time.Sleep(3 * time.Second)
	test.Emit("enter", []string{"33333", "33333"})
	test.Emit("move", []string{"44444", "44455", "444466"})
	time.Sleep(3 * time.Second)
	test.Emit("scroll", "77 scroll-2 77")
	test.Emit("size", &PayloadCustom{ww: "Привет size"})
	time.Sleep(3 * time.Second)

	fmt.Println("Количество типов:", test.Count())
	fmt.Println("Все типы:", test.GetListenerName())
	fmt.Println("Зарегистрировано enter", test.CountListener("enter"))
	fmt.Println("Было зарегистрировано enter", test.DeleteAllListener("enter"))

}
