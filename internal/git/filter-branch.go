package git

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func (r *Repo) FilterBranch(args ...string) (err error) {
	var path string
	if r.WorkDir == "" {
		path = r.GitDir
	} else {
		path = r.WorkDir
	}

	cmdArgs := []string{"filter-branch"}
	cmdArgs = append(cmdArgs, args...)

	stdout, stderr := os.Stdout, new(bytes.Buffer)

	cmd := exec.Command(gitCmd, cmdArgs...)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	cmd.Dir = path
	cmd.Env = []string{
		"FILTER_BRANCH_SQUELCH_WARNING=1",
	}

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
