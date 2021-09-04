**_ArkEvents_**

Реализуем события,

- Умеем создавать события,
- Подписываться на события,
- Посылать события...

--

    import "github.com/ark-go/arkEvents/pkg/watch"

    test := watch.NewWatch()

    go func() {
    	enter := test.AddListener("enter") // создает канал
    	for {
    		msg := <-enter
    		//! Желательно использовать проверку приведения типа !
    		if val, ok := msg.([]string); ok {
    			fmt.Println("enter 1:", val)
    		}
    	}
    }()

    test.Emit("enter", []string{"11111", "111"})

смотрите пример:  
https://github.com/ark-go/arkEvents/blob/main/cmd/main/main.go

**Emit** - можно безболезненно встраивать в свой код, пока не добавлен **AddListener** он ничего не делает, но при добавлении слушателя, **Emit** будет создавать горутины для отправки данных в канал, в случае если каналы некому будет читать то горутины будут создаваться на каждый вызов **Emit("name"..**, поэтому необходимо в своем коде, при регистрации **chanX := test.AddListener**, обеспечить чтение из канала chanX

upd:
**AddListenerFunc** - возвращает канал только для того чтобы по нему можно было выполнить **RemoveListener**, в этот канал никогда не будут посланы данные.

    // функция вызывается в горутине
    mm := test.AddListenerFunc("funcTest", func(s interface{}) {
    fmt.Println("Вызов функции 111:", s)
    })


    test.Emit("funcTest", []string{"2", "44", "55"})
