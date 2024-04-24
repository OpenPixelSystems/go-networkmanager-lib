package interfaces

type Cmd interface {
	Run() error
	Output() ([]byte, error)
	CombinedOutput() ([]byte, error)
}

type Exec interface {
	Command(name string, arg ...string) Cmd
}
