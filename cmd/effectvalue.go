package cmd

import (
	"errors"

	"github.com/pugkong/ylc/yeelight"
)

type effectValue yeelight.Effect

func newEffectValue(value *yeelight.Effect) *effectValue {
	return (*effectValue)(value)
}

func (e *effectValue) String() string {
	return string(*e)
}

var ErrUnknownEffect = errors.New("unknown effect")

func (e *effectValue) Set(value string) error {
	switch value {
	case "sudden":
		*e = effectValue(yeelight.EffectSudden)
	case "smooth":
		*e = effectValue(yeelight.EffectSmooth)
	default:
		return ErrUnknownEffect
	}

	return nil
}

func (e *effectValue) Type() string {
	return "effect"
}
