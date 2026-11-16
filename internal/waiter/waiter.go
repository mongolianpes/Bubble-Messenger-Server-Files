package waiter

import (
	"sync"
)

type WaiterService struct {
	Mu      sync.Mutex
	Waiters map[string]chan WaiterInfo
}

type WaiterInfo struct {
	Path string
	Err  error
}

func NewWaiterService() *WaiterService {
	return &WaiterService{
		Waiters: make(map[string]chan WaiterInfo),
	}
}
