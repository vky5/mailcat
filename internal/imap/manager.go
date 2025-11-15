package imap

import (
	"sync"

	"github.com/emersion/go-imap/client"
	"github.com/vky5/mailcat/internal/db/models"
)

var (
	connPool = make(map[uint]*client.Client)
	mu       sync.RWMutex
)

func GetConnection(acc models.Account) (*client.Client, error) {
	mu.RLock()
	conn, exists := connPool[acc.ID]
	mu.RUnlock()

	if exists {
		return conn, nil
	}

	// Make new connection
	newConn, err := ConnectIMAP(acc)
	if err != nil {
		return nil, err
	}

	// Store it
	mu.Lock()
	connPool[acc.ID] = newConn
	mu.Unlock()

	return newConn, nil
}

func CloseAll() {
	mu.Lock()
	defer mu.Unlock()

	for id, conn := range connPool {
		conn.Logout()
		delete(connPool, id)
	}
}
