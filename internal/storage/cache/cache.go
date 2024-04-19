package storage

import (
	"github.com/mishaRomanov/wb-l0/internal/entities"
	"sync"
)

// Cache struct
type OrdersCache struct {
	Orders map[string]entities.Order
	mu     sync.RWMutex
}

// returns a new OrdersCache instance
func New() *OrdersCache {
	return &OrdersCache{
		Orders: make(map[string]entities.Order),
	}
}

// Add method adds a new order to the cache
func (c *OrdersCache) Add(value entities.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := value.OrderUID
	c.Orders[key] = value
}

// Get method returns the order, if exists
func (c *OrdersCache) Get(key string) (entities.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.Orders[key]
	if !ok {
		return entities.Order{}, false
	}
	return value, ok
}
