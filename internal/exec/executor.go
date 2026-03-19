package exec

type Executor interface {
	Run(name string, args ...string) error
	RunInteractive(name string, args ...string) error
	RunDetached(name string, args ...string) error
	Output(name string, args ...string) (string, error)
	CombinedOutput(name string, args ...string) (string, error)
}
