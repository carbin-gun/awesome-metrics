package metrics

import (
	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/carbin-gun/awesome-metrics/metrics"
	"github.com/carbin-gun/awesome-metrics/registry"
)

/***
It's a helper file of the common entrance of monitor.
You can new a registry and the register all kinds of metrics supported for now for monitoring usage.
*/
type RegistryWrapper struct {
	Registry registry.Registry
}

func NewRegistry() RegistryWrapper {
	return &RegistryWrapper{
		Registry: registry.NewRegistry(),
	}
}

func NewPrefixRegistry(prefix string) RegistryWrapper {
	return &RegistryWrapper{
		Registry: registry.NewPrefixRegistry(prefix),
	}
}

func (r *RegistryWrapper) Timer(name string) mechanism.Timer {
	return r.Registry.GetOrRegister(name, metrics.NewTimer())
}
func (r *RegistryWrapper) Counter(name string) mechanism.Counter {
	return r.Registry.GetOrRegister(name, metrics.NewCounter())
}
func (r *RegistryWrapper) Meter(name string) mechanism.Meter {
	return r.Registry.GetOrRegister(name, metrics.NewMeter())
}
