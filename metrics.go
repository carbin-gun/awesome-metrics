package metrics

import (
	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/carbin-gun/awesome-metrics/registry"
)

//helper file.all the functions can be done within the  packages in this repo

var (
	DEFAULT_REGISTRY = NewRegistry()
)

func NewRegistry() registry.Registry {
	return registry.NewRegistry()
}

func Timer() mechanism.Timer {

}
func Counter() mechanism.Counter {

}
func Meter() mechanism.Meter {

}
