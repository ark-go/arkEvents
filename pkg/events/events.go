package events

import (
	"sync"
)

// interface для передачи полезной нагрузки в Emit
type iPayload interface {
}

// watch - хранит все прослущиватели / подписки
type watch struct {
	name      string                     // имя наблюдателя
	listeners map[string][]chan iPayload // слушатели
	mu        sync.RWMutex
}

// Создает Watch
func NewWatch() *watch {
	return &watch{
		listeners: make(map[string][]chan iPayload),
	}
}

// вернет количество зарегистрированных событий
func (w *watch) Count() int {
	return len(w.listeners)

}
func (w *watch) GetListenerName() []string {
	m := []string{}
	for key := range w.listeners {
		m = append(m, key)
	}
	return m
}

// name - название/тип события
//	вернет количество подписчиков на событие
func (w *watch) CountListener(name string) int {
	return len(w.listeners[name])
}

// Удаляет все регистрации по имени - name
//	возвращает сколько было зарегистрировано слушателей
//	name - название/тип события
func (w *watch) DeleteAllListener(name string) (count int) {
	if _, ok := w.listeners[name]; ok {
		count = len(w.listeners[name])
		delete(w.listeners, name)
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
				ch <- message
			}(ch)
			count++
		}
	}
	return
}
