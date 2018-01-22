package tachikoma

import (
	"github.com/mitchellh/mapstructure"
)

func Unmap(v map[string]interface{}, t interface{}) {
	config := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: t}
	decoder, err := mapstructure.NewDecoder(&config)
	Fatal(err)

	err = decoder.Decode(v)
	Fatal(err)
}
