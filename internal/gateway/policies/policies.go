// Package policies provides static and (future) dynamic policies that can be activated per AI instance and enforced by the gateway to control request handling.
package policies

import "net/http"

type Policy interface {
	Name() string
	IsActive() bool
	Check(req *http.Request) error
}

type StaticPolicy struct {
	name   string
	active bool
}

func (p *StaticPolicy) Name() string   { return p.name }
func (p *StaticPolicy) IsActive() bool { return p.active }
func (p *StaticPolicy) Check(req *http.Request) error {
	return nil // always pass, stub for future logic
}
