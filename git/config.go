package git

import (
	"strings"

	. "github.com/crazy-max/git-rewrite-author/utils"
)

func (r *Repo) readConfig() {
	if r.cfg != nil {
		return
	}
	cmd, stdout, stderr := r.Git("config", "-l", "-z")
	if err := cmd.Run(); err != nil {
		Error(stderr.String())
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
	return
}

// ReloadConfig will force the config for this git repo to be lazily reloaded.
func (r *Repo) ReloadConfig() {
	r.cfg = nil
}

// Get a specific config value.
func (r *Repo) Get(key string) (val string, found bool) {
	r.readConfig()
	val, found = r.cfg[key]
	return
}

// Set a config variable.
func (r *Repo) Set(key, val string) {
	r.readConfig()
	cmd, _, _ := r.Git("config", "--add", key, val)
	if err := cmd.Run(); err != nil {
		Error(err.Error())
	}
	r.cfg[key] = val
}

// Find all config variables with a specific prefix.
func (r *Repo) Find(prefix string) (res map[string]string) {
	r.readConfig()
	res = make(map[string]string)
	for k, v := range r.cfg {
		if strings.HasPrefix(k, prefix) {
			res[k] = v
		}
	}
	return res
}
