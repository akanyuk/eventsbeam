package options

type Option func(*Options)

type Options struct {
	SkipRemoveUnusedModels bool
	ExcludePackages        []string
}

func SkipRemoveUnusedModels() Option {
	return func(o *Options) {
		o.SkipRemoveUnusedModels = true
	}
}

func ExcludePackages(packages ...string) Option {
	return func(o *Options) {
		o.ExcludePackages = packages
	}
}

func Do(opt ...Option) Options {
	options := Options{}
	for _, callback := range opt {
		callback(&options)
	}

	return options
}
