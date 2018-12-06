package types

import (
	"strings"

	"helputils"
	"optional"
)

type ProductParameters struct {
	Boundaries Boundaries `yaml:"Boundaries"`
	Timeouts   Timeouts   `yaml:"Timeouts"`
}

func (p ProductParameters) ToString() string {
	return "{" + strings.Join(helputils.StructToStrings(p), " ") + "}"
}

type Boundaries struct {
}

type Timeouts struct {
	SystemSleep optional.String `yaml:"SystemSleep" default:"15s"`
	RemoveNode  optional.String `yaml:"RemoveNode" default:"6h"`
}

func NewProductDefaults() ProductParameters {
	return ProductParameters{Timeouts: Timeouts{
		RemoveNode:  optional.NewString("6h"),
		SystemSleep: optional.NewString("15s"),
	}}
}
