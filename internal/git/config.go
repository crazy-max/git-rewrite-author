package git

import (
	"strings"

	"github.com/pkg/errors"
)

func (r *Repo) readConfig() error {
	if r.cfg != nil {
		return nil
	}

	cmd, stdout, stderr := r.Git("config", "-l", "-z")
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, stderr.String())
	}
	r.cfg = make(ConfigMap)

	for _, line := range strings.Split(stdout.String(), "\x00") {
		parts := strings.SplitN(line, "\n", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" {
			continue
		}
		r.cfg[k] = v
	}

	return nil
}

// ReloadConfig will force the config for this git repo to be lazily reloaded.
func (r *Repo) ReloadConfig() {
	r.cfg = nil
}

// Get a specific config value.
func (r *Repo) Get(key string) (val string, found bool, err error) {
	err = r.readConfig()
	val, found = r.cfg[key]
	return
}

// Set a config variable.
func (r *Repo) Set(key, val string) error {
	if err := r.readConfig(); err != nil {
		return err
	}

	cmd, _, _ := r.Git("config", "--add", key, val)
	if err := cmd.Run(); err != nil {
		return err
	}

	r.cfg[key] = val
	return nil
}

// Find all config variables with a specific prefix.
func (r *Repo) Find(prefix string) (res map[string]string, err error) {
	if err := r.readConfig(); err != nil {
		return nil, err
	}

	res = make(map[string]string)
	for k, v := range r.cfg {
		if strings.HasPrefix(k, prefix) {
			res[k] = v
		}
	}

	return res, nil
}
