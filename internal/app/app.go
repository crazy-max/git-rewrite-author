package app

import (
	"github.com/crazy-max/git-rewrite-author/internal/git"
	"github.com/crazy-max/git-rewrite-author/internal/model"
)

// GitRewriteAuthor represents an active git-rewrite-author object
type GitRewriteAuthor struct {
	cli  model.Cli
	repo *git.Repo
}

// New creates new git-rewrite-author instance
func New(cli model.Cli) (*GitRewriteAuthor, error) {
	repo, err := git.Open(cli.Repo)
	return &GitRewriteAuthor{
		cli:  cli,
		repo: repo,
	}, err
}
