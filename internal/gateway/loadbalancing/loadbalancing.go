// Package loadbalancing provides interfaces for distributing requests across multiple AI resource instances.
package loadbalancing

type InstanceSelector interface {
	Select(instances []string, key string) string
}
