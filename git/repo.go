package git

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/crazy-max/git-rewrite-author/utils"
)

// ConfigMap maps config keys to their values.
type ConfigMap map[string]string

// Repo is the main struct that we use to track Git repositories.
type Repo struct {
	GitDir  string
	WorkDir string
	cfg     ConfigMap
}

var gitCmd string

func init() {
	var err error
	if gitCmd, err = exec.LookPath("git"); err != nil {
		Error("Cannot find git command!")
	}
}

func findRepo(path string) (found bool, gitdir, workdir string) {
	stat, err := os.Stat(path)
	if err != nil {
		Error("Could not stat " + path)
	}

	if !stat.IsDir() {
		Error(path + " is not a directory!")
	}

	if strings.HasSuffix(path, ".git") {
		if stat, err = os.Stat(filepath.Join(path, "config")); err == nil {
			found = true
			gitdir = path
			workdir = ""
			return
		}
	}

	if stat, err = os.Stat(filepath.Join(path, ".git", "config")); err != nil {
		found = false
		return
	}

	found = true
	gitdir = filepath.Join(path, ".git")
	workdir = path
	return
}

// Open the first git repository that "owns" path.
func Open(path string) (repo *Repo, err error) {
	if path == "" {
		path = "."
	}

	path, err = filepath.Abs(path)
	basepath := path
	if err != nil {
		return
	}

	for {
		found, gitdir, workdir := findRepo(path)
		if found {
			repo = new(Repo)
			repo.GitDir = gitdir
			repo.WorkDir = workdir
			return
		}
		parent := filepath.Dir(path)
		if parent == path {
			break
		}
		path = parent
	}

	return nil, errors.New(fmt.Sprintf("Could not find a Git repository in %s or any of its parents!", basepath))
}

// Git is a helper for creating exec.Cmd types and arranging to capture
// the output and erro streams of the command into bytes.Buffers
func Git(cmd string, args ...string) (res *exec.Cmd, stdout, stderr *bytes.Buffer) {
	cmdArgs := make([]string, 1)
	cmdArgs[0] = cmd
	cmdArgs = append(cmdArgs, args...)
	res = exec.Command(gitCmd, cmdArgs...)
	stdout, stderr = new(bytes.Buffer), new(bytes.Buffer)
	res.Stdout, res.Stderr = stdout, stderr
	return
}

// Git is a helper for making sure that the Git command runs in the proper repository.
func (r *Repo) Git(cmd string, args ...string) (res *exec.Cmd, out, err *bytes.Buffer) {
	var path string
	if r.WorkDir == "" {
		path = r.GitDir
	} else {
		path = r.WorkDir
	}
	res, out, err = Git(cmd, args...)
	res.Dir = path
	return
}

// IsRaw checks to see if this is a raw repository.
func (r *Repo) IsRaw() (res bool) {
	return r.WorkDir == ""
}

// Path returns the best idea of the path to the repository.
// The exact value returned depends on whether this is a
// raw repository or not.
func (r *Repo) Path() (path string) {
	if r.IsRaw() {
		return r.GitDir
	}
	return r.WorkDir
}
