**_ArkEvents_**

Реализуем события,

- Умеем создавать события,
- Подписываться на события,
- Посылать события...

--

    import "github.com/ark-go/arkEvents/pkg/watch"

    test := watch.NewWatch("test")

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

    test.Emit("enter", []string{"11111", "111"})

пример:  
github.com/ark-go/arkEvents/cmd/main/main.go
