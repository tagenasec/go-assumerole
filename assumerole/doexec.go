package assumerole

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/apex/log"
)

func doExec(executable string, args []string, environ map[string]string) error {
	absoluteExecutable, err := exec.LookPath(executable)
	if err != nil {
		log.WithField("args", args).WithError(err).WithField("executable", executable).Fatal("Unable to find executable")
		return err
	}
	realArgs := []string{executable}
	realArgs = append(realArgs, args...)
	finalEnviron := append([]string{}, os.Environ()...)
	for key, value := range environ {
		finalEnviron = append(finalEnviron, key+"="+value)
	}
	err = syscall.Exec(absoluteExecutable, realArgs, finalEnviron)
	if err != nil {
		log.WithField("args", args).WithError(err).WithField("executable", absoluteExecutable).Fatal("Unable to exec")
		return err
	}
	return nil
}
