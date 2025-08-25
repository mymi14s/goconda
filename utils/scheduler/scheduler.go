package scheduler

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	c   *cron.Cron
	mu  sync.Mutex
	reg = map[string]cron.EntryID{}
)

// Start starts the global scheduler (idempotent).
func Start() {
	mu.Lock()
	defer mu.Unlock()
	if c == nil {
		c = cron.New(cron.WithSeconds())
		c.Start()
	}
}

// Stop stops the global scheduler.
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if c != nil {
		c.Stop()
		c = nil
		reg = map[string]cron.EntryID{}
	}
}

// Register registers a function with a cron spec and a unique name.
// Example spec: "*/5 * * * * *" (every 5 seconds)
func Register(name, spec string, fn func()) (cron.EntryID, error) {
	mu.Lock()
	defer mu.Unlock()
	if c == nil { Start() }
	if id, ok := reg[name]; ok {
		// remove previous before adding again
		c.Remove(id)
	}
	id, err := c.AddFunc(spec, fn)
	if err != nil { return 0, err }
	reg[name] = id
	return id, nil
}
