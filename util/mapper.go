package util

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

func BindingStruct(src interface{}, desc interface{}) error {
	// convert to byte
	byteSrc, err := json.Marshal(src)
	if err != nil {
		return err
	}
	// binding to desc
	err = json.Unmarshal(byteSrc, &desc)
	if err != nil {
		return err
	}

	return nil
}

type MapStructOption func(*mapstructure.DecoderConfig)

// Option for string to time conversion.
func WithDecodeTimeFormat(fm string) MapStructOption {
	return func(c *mapstructure.DecoderConfig) {
		if c.DecodeHook != nil {
			c.DecodeHook = mapstructure.ComposeDecodeHookFunc(
				c.DecodeHook,
				mapstructure.StringToTimeHookFunc(fm),
			)
		} else {
			c.DecodeHook = mapstructure.StringToTimeHookFunc(fm)
		}
	}
}

// Option for weakly type mapping. Default true
func WithWeaklyTypedInput(enabled bool) MapStructOption {
	return func(c *mapstructure.DecoderConfig) {
		c.WeaklyTypedInput = true
	}
}

func MapStruct(input interface{}, output interface{}, opts ...MapStructOption) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &output,
	}

	for _, opt := range opts {
		opt(config)
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(input)
	return err
}
