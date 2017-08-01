package meander

import (
	"strings"
	"github.com/pkg/errors"
)

type Cost int8

const (
	_ Cost = iota
	Cost1
	Cost2
	Cost3
	Cost4
	Cost5
)

var costStrings = map[string]Cost{
	"$": Cost1,
	"$$": Cost2,
	"$$$": Cost3,
	"$$$$": Cost4,
	"$$$$$": Cost5,
}

func (c Cost) String() string {
	return strings.Repeat("$", int(c))
}

func ParseCost(s string) Cost {
	return costStrings[s]
}

type CostRange struct {
	From Cost
	To Cost
}

func (r CostRange) String() string {
	return r.From.String() + "..." + r.To.String()
}

func ParseCostRange(s string) (CostRange, error) {
	var r CostRange
	segs := strings.Split(s, "...")
	if len(segs) != 2 {
		return r, errors.New("Invalid cost range")
	}
	r.From = ParseCost(segs[0])
	r.To = ParseCost(segs[1])
	return r, nil
}