package lazywriter

import (
	"fmt"
	"time"
)

// LazyWriter -
type LazyWriter struct {
	Interval time.Duration
	Name     string
	store    *Cache
	ticker   *time.Ticker
}

// LazyObject -
type LazyObject struct {
	ID       string
	Data     interface{}
	WriterFn WriterFn
}

// WriterFn -
type WriterFn func(key string, value interface{}) error

var (
	// DefaultInterval - Default interval to execute write function, 5 sec
	DefaultInterval = time.Second * 5
)

// NewLazyWriter -
func NewLazyWriter(name string, interval time.Duration) *LazyWriter {
	if interval.Seconds() < DefaultInterval.Seconds() {
		interval = DefaultInterval
	}
	return &LazyWriter{
		Name:     name,
		Interval: interval,
		store:    NewStore(),
		ticker:   time.NewTicker(interval),
	}
}

// NewLazyObject -
func NewLazyObject(id string, data interface{}, writerFn WriterFn) *LazyObject {
	return &LazyObject{
		ID:       id,
		Data:     data,
		WriterFn: writerFn,
	}
}

// Start - Start Lazy Writer task in new go routine
func (lw *LazyWriter) Start() {
	go func() {
		for {
			select {
			case <-lw.ticker.C:
				// lw.store.Iterate()
				fmt.Println("Processing lazy objects")
				lw.processWriter()
			}
		}
	}()
}

func (lw *LazyWriter) processWriter() {
	fmt.Println("In process Function")
	// lw.store.m.Lock()
	for _, id := range lw.store.Keys() {
		lo := lw.store.MustGet(id)
		if err := lo.WriterFn(lo.ID, lo.Data); err != nil {
			fmt.Println("ERROR >> ", id, " >> ", err)
		}
	}
	// lw.store.m.Unlock()
}

// Add -
func (lw *LazyWriter) Add(lazyObj *LazyObject) {
	// lw.store.SetDefault(lazyObj.ID, lazyObj)
	fmt.Println("ADD >> ", lazyObj.ID)
	lw.store.Add(lazyObj.ID, lazyObj)
}

// Get -
func (lw *LazyWriter) Get(key string) (*LazyObject, bool) {
	return lw.store.Get(key)
}
