package loadbalancing

import "sync"

type RoundRobinSelector struct {
	mu       sync.Mutex
	counters map[string]int
}

func NewRoundRobinSelector() *RoundRobinSelector {
	return &RoundRobinSelector{counters: make(map[string]int)}
}

func (r *RoundRobinSelector) Select(instances []string, key string) string {
	r.mu.Lock()

	defer r.mu.Unlock()

	if len(instances) == 0 {
		return ""
	}

	idx := r.counters[key] % len(instances)

	r.counters[key]++

	return instances[idx]
}
