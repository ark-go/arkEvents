package main

import (
	"fmt"
	"time"

	"github.com/ark-go/arkEvents/pkg/events"
)

func main() {
	// создадим хранилище для наших каналов
	test := events.NewWatch()
	// структура для передачи данных в событиях
	type PayloadCustom struct {
		ww string
	}
	// в горутине создадим канал и будем его ждать
	go func() {
		enter := test.AddListener("enter") // вернет новый канал
		for {
			msg := <-enter // emit пошлет в канал данные
			//! Желательно использовать проверку приведения типа !
			if val, ok := msg.([]string); ok {
				fmt.Println("enter 1:", val) // выполним нашу функцию
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
	// пример удаления слушателя
	go func() {
		enter := test.AddListener("enter")
		for {
			test.RemoveListener("enter", enter) //! удаляем Listener
			msg := <-enter
			if val, ok := msg.([]string); ok {
				fmt.Println("Event3:", val) // никогда не выполнится / Will never come true
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
	// функция вызывается в горутине
	mm := test.AddListenerFunc("funcTest", func(s interface{}) {
		fmt.Println("Вызов функции 111:", s)
	})
	test.AddListenerFunc("funcTest", func(s interface{}) {
		fmt.Println("Вызов функции 222:", s)
	})
	time.Sleep(1 * time.Second) // мы должны запустить все горутины
	fmt.Println("Кол-во зарегистрироанных функций", test.CountRegFunc())

	test.Emit("scroll", "scrooooooooooooooool")
	test.Emit("enter", []string{"11111", "111"})
	test.Emit("roll", []string{"222222"})

	test.Emit("enter", []string{"33333", "33333"})
	test.Emit("move", []string{"44444", "11111", "444466"})
	test.Emit("move", []string{"44444", "22222", "444466"})
	test.Emit("move", []string{"44444", "33333", "444466"})
	test.Emit("scroll", "77 scroll-2 77")
	test.Emit("size", &PayloadCustom{ww: "Привет size"})
	test.Emit("funcTest", 90)
	time.Sleep(1 * time.Second) // дадим выполнится горутине по funcTest
	test.RemoveListener("funcTest", mm)
	fmt.Println("Кол-во зарегистрироанных функций", test.CountRegFunc())
	test.Emit("funcTest", []string{"2", "44", "55"})
	fmt.Println("Количество типов:", test.Count())
	fmt.Println("Все типы:", test.GetListenerName())
	fmt.Println("Зарегистрировано enter", test.CountListener("enter"))
	fmt.Println("Было зарегистрировано enter", test.DeleteAllListener("enter"))
	time.Sleep(15 * time.Second) // ждем последний emit 1990
}

// func testfunc(str interface{}) {
// 	fmt.Println("Вызов функции из Listener:", str[0].(string))
// }
