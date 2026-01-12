package provider

import "github.com/gookit/validate"

func NewValidationProvider() {

	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = true
	})
}
