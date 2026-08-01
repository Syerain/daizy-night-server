package confx

type options struct {
	searchPaths []string
	fileType    string
	fileName    string
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		searchPaths: []string{".", "./config"},
		fileType:    "yaml",
		fileName:    "config",
	}
}

// WithSearchPaths
func WithSearchPaths(paths ...string) Option {
	return func(o *options) {
		o.searchPaths = paths
	}
}

// WithFileType yaml/json/toml
func WithFileType(t string) Option {
	return func(o *options) {
		o.fileType = t
	}
}

func WithFileName(n string) Option {
	return func(o *options) {
		o.fileName = n
	}
}
