package impl

import (
	"os/exec"

	"github.com/openpixelsystems/go-networkmanager-lib/interfaces"
)

type Cmd struct {
	cmd *exec.Cmd
}

func (cmd *Cmd) Run() error {
	return cmd.cmd.Run()
}

func (cmd *Cmd) Output() ([]byte, error) {
	return cmd.cmd.Output()
}

func (cmd *Cmd) CombinedOutput() ([]byte, error) {
	return cmd.cmd.CombinedOutput()
}

type Exec struct {
}

func (Exec) Command(name string, arg ...string) interfaces.Cmd {
	return &Cmd{exec.Command(name, arg...)}
}
