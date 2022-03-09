package worker

import (
	"os"
	"os/exec"
)

type fakeWriter struct{}

func (_ fakeWriter) Write(_ []byte) (_ int, _ error) { return }

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Stdout = fakeWriter{}
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
