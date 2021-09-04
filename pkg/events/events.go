package events

import (
	"sync"
)

// interface для передачи полезной нагрузки в Emit
type iPayload interface {
}

// watch - хранит все прослущиватели / подписки
type watch struct {
	listeners map[string][]chan iPayload // слушатели
	mu        sync.RWMutex
	handlers  map[chan iPayload]func(interface{})
}

// type watchFunc struct {
// 	listeners map[string][]chan iPayload // слушатели
// 	mu        sync.RWMutex
// 	fun       func(interface{})
// }

// Создает Watch, для хранения зарегистрированных Listeners т.е. созданных потоков
func NewWatch() *watch {
	return &watch{
		listeners: make(map[string][]chan iPayload),
		handlers:  make(map[chan iPayload]func(interface{})),
	}
}

// вернет количество зарегистрированных событий
func (b *watch) Count() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.listeners)

}

// возвращает все имена зарегистрированных событий
func (b *watch) GetListenerNames() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	m := []string{}
	for key := range b.listeners {
		m = append(m, key)
	}
	return m
}

// name - название/тип события
//	вернет количество подписчиков на событие
func (b *watch) CountListener(name string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.listeners[name])
}

// Количество зарегистрированных функций
func (b *watch) CountRegFunc() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers)
}

// Удаляет все регистрации по имени - name
//	возвращает сколько было зарегистрировано слушателей
//	name - название/тип события
func (b *watch) DeleteAllListener(name string) (count int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.listeners[name]; ok {
		count = len(b.listeners[name])
		delete(b.listeners, name)
	}
	return
}

// добавляет прослушиватель события к экземпляру структуры watch.
//	возвращает канал для события
//	name - название/тип события
func (b *watch) AddListener(name string) chan iPayload {
	ch := make(chan iPayload)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners[name] = append(b.listeners[name], ch)
	return ch
}

// удаляет listener из экземпляра структуры watch.
//	по каналу chan
//	name - название/тип события
//	ch - канал созданный AddListener
func (b *watch) RemoveListener(name string, ch chan iPayload) (count int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.listeners[name]; ok {
		for i := range b.listeners[name] {
			// перебираем всех подписавшихся
			if b.listeners[name][i] == ch {
				b.listeners[name] = append(b.listeners[name][:i], b.listeners[name][i+1:]...)
				count++
				break
			}
		}
		delete(b.handlers, ch)

	}
	return
}

// Генерируем событие
//	возвращаем количество получателей, зарегистрированных слушателей
//	name - название/тип события
func (b *watch) Emit(name string, message iPayload) (count int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if _, ok := b.listeners[name]; ok {
		for _, ch := range b.listeners[name] {
			go func(ch chan iPayload) {
				if hndl, ok := b.handlers[ch]; ok {
					hndl(message)
				} else {
					ch <- message
				}
			}(ch)
			count++
		}
	}
	return
}

//---------------- test
func (b *watch) AddListenerFunc(name string, handler func(interface{})) chan iPayload {
	ch := make(chan iPayload)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners[name] = append(b.listeners[name], ch)
	b.handlers[ch] = handler
	return ch
}
