// Автор: Leonchik Aleksandr

package safemap

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	year       = 1939 // начало Второй мировой войны
	goroutines = 4
	perKey     = goroutines - 1 // каждый ключ обрабатывается всеми горутинами кроме одной
)

func TestSafeMapConcurrent(t *testing.T) {
	t.Parallel()

	m := New()
	var wg sync.WaitGroup

	for g := range goroutines {
		// Каждая горутина обрабатывает ключи, для которых (key-1)%goroutines != g.
		// Таким образом каждый ключ обрабатывается ровно perKey горутинами из goroutines.
		keys := make([]int, 0, year*perKey/goroutines+1) // +1 гарантирует что capacity хватит без реаллокации
		// Диапазон ключей от 1 до выбранного года
		for key := 1; key <= year; key++ {
			if (key-1)%goroutines != g {
				keys = append(keys, key)
			}
		}

		// Перемешиваем, чтобы горутины не обращались к ключам последовательно.
		rand.Shuffle(len(keys), func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
		})

		wg.Go(func() {
			for _, key := range keys {
				m.Lock()
				val := m.Get(key)
				m.Set(key, val+1)
				m.Unlock()
			}
		})
	}

	wg.Wait()

	// Каждый ключ должен иметь значение perKey.
	for key := 1; key <= year; key++ {
		require.Equal(t, perKey, m.data[key], "key %d", key)
	}

	// Счётчик обращений = год * perKey.
	require.Equal(t, year*perKey, m.Accesses)

	// Счётчик добавлений = год.
	require.Equal(t, year, m.Additions)
}
