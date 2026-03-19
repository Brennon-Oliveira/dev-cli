package exec

type realExecutor struct{}

func NewExecutor() *realExecutor {
	return &realExecutor{}
}

// func (e *realExecutor) Run(name string, args ...string) error {
// 	cmd := exec.Command(name, args...)
// 	cmd.Stdout = os.Stdout
// }
