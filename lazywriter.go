package lazywriter

import (
	"time"
)

type state int

const (
	notStarted state = iota
	running
)

// LazyWriter -
type LazyWriter struct {
	Interval time.Duration
	Name     string
	store    *Cache
	ticker   *time.Ticker
	state    state
}

// LazyObject -
type LazyObject struct {
	ID       string
	Data     interface{}
	WriterFn WriterFn
	updated  bool
	locked   bool
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
		state:    notStarted,
	}
}

// NewLazyObject -
func NewLazyObject(id string, data interface{}, writerFn WriterFn) LazyObject {
	return LazyObject{
		ID:       id,
		Data:     data,
		WriterFn: writerFn,
	}
}

// Start - Start Lazy Writer task in new go routine
func (lw *LazyWriter) Start() {
	if lw.state == running {
		return
	}

	lw.state = running

	go func() {
		for {
			select {
			case <-lw.ticker.C:
				lw.processWriter()
			}
		}
	}()
}

func (lw *LazyWriter) processWriter() {
	for _, id := range lw.store.Keys() {
		lo := lw.store.MustGet(id)
		if lo == nil {
			continue
		}

		if !lo.updated {
			continue
		}

		lo.locked = true

		if err := lo.WriterFn(lo.ID, lo.Data); err != nil {
			// handle error
		}

		lo.locked = false
	}
}

// Add -
func (lw *LazyWriter) Add(lazyObj LazyObject) {
	lo, ok := lw.store.Get(lazyObj.ID)
	for lo.locked {
		// TODO: mechanism to lock lo when the data is being written
		// better one, not feasible with go routines.
	}
	if ok {
		lo.updated = true
		lo.Data = lazyObj.Data
		lo.WriterFn = lazyObj.WriterFn
		return
	}
	lazyObj.updated = true
	lw.store.Add(lazyObj.ID, &lazyObj)
}

// Get -
func (lw *LazyWriter) Get(key string) (LazyObject, bool) {
	lo, ok := lw.store.Get(key)
	return *lo, ok
}
