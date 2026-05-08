package cli

// ExitCodeFunc maps an application execution error to a process exit code.
type ExitCodeFunc func(error) int

// ApplicationOption configures an Application constructed by NewApplication.
//
// Options are applied in the order provided to NewApplication. If multiple options configure
// the same field, the last one wins.
type ApplicationOption interface {
	apply(opts *applicationOpts)
}

type applicationOpts struct {
	exitCode ExitCodeFunc
}

type applicationOptionFunc func(*applicationOpts)

func (f applicationOptionFunc) apply(o *applicationOpts) {
	f(o)
}

// WithExitCodeFunc configures how RunCode chooses a process exit code for an error.
//
// If f returns zero or a negative value for a non-nil error, the application falls back to exit code 1.
func WithExitCodeFunc(f ExitCodeFunc) ApplicationOption {
	return applicationOptionFunc(func(o *applicationOpts) {
		if f != nil {
			o.exitCode = f
		}
	})
}

func newApplicationOpts(opts ...ApplicationOption) *applicationOpts {
	options := &applicationOpts{exitCode: defaultExitCode}
	for _, opt := range opts {
		opt.apply(options)
	}

	return options
}

func defaultExitCode(error) int {
	return 1
}
