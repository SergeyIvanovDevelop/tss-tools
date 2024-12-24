package interprocmutex

import (
	"sync"
	"time"

	"github.com/gofrs/flock"
)

type FileRWMutex struct {
	fileLock  *flock.Flock
	readCount int
	mu        sync.Mutex
}

// Конструктор для FileRWMutex
func NewFileRWMutex(filePath string) (*FileRWMutex, error) {
	return &FileRWMutex{
		fileLock: flock.New(filePath),
	}, nil
}

// Lock устанавливает эксклюзивную блокировку
func (m *FileRWMutex) Lock() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Ждем, пока блокировка не станет доступной
	for {
		locked, err := m.fileLock.TryLock()
		if err != nil {
			panic(err)
		}
		if locked {
			return
		}
		time.Sleep(10 * time.Millisecond) // Пауза перед повторной попыткой
	}
}

// Unlock снимает эксклюзивную блокировку
func (m *FileRWMutex) Unlock() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.fileLock.Unlock(); err != nil {
		panic(err)
	}
}

// RLock устанавливает разделяемую (shared) блокировку
func (m *FileRWMutex) RLock() {
	m.mu.Lock()
	m.readCount++
	if m.readCount == 1 {
		// Первый читатель должен захватить блокировку
		for {
			locked, err := m.fileLock.TryRLock()
			if err != nil {
				panic(err)
			}
			if locked {
				break
			}
			time.Sleep(10 * time.Millisecond) // Пауза перед повторной попыткой
		}
	}
	m.mu.Unlock()
}

// RUnlock снимает разделяемую (shared) блокировку
func (m *FileRWMutex) RUnlock() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.readCount--
	if m.readCount == 0 {
		// Последний читатель снимает блокировку
		if err := m.fileLock.Unlock(); err != nil {
			panic(err)
		}
	}
}
