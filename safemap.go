// Автор: Leonchik Aleksandr

package safemap

import "sync"

type SafeMap struct {
	mu        sync.Mutex
	data      map[int]int
	Accesses  int // счётчик обращений к ключам
	Additions int // счётчик добавлений новых ключей
}

// New создаёт пустой SafeMap.
func New() *SafeMap {
	return &SafeMap{data: make(map[int]int)}
}

// Lock захватывает мьютекс.
func (m *SafeMap) Lock() {
	m.mu.Lock()
}

// Unlock освобождает мьютекс.
func (m *SafeMap) Unlock() {
	m.mu.Unlock()
}

// Get возвращает значение по ключу. Если ключа нет - создаёт запись со значением 0.
// Вызывать только при захваченном мьютексе.
func (m *SafeMap) Get(key int) int {
	m.Accesses++
	val, ok := m.data[key]
	if !ok {
		m.Additions++
		m.data[key] = 0

		return 0
	}

	return val
}

// Set устанавливает значение по ключу.
// Вызывать только при захваченном мьютексе.
func (m *SafeMap) Set(key int, val int) {
	m.data[key] = val
}
