package loadbalancing

import "math/rand"

type RandomSelector struct{}

func (r *RandomSelector) Select(instances []string, key string) string {
	if len(instances) == 0 {
		return ""
	}
	return instances[rand.Intn(len(instances))]
}
