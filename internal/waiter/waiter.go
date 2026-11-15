package waiter

import (
	"sync"
)

type WaiterService struct {
	Mu      sync.Mutex
	Waiters map[string]chan error
}

func NewWaiterService() *WaiterService {
	return &WaiterService{
		Waiters: make(map[string]chan error),
	}
}
