package app

import (
	"github.com/crazy-max/git-rewrite-author/internal/git"
	"github.com/crazy-max/git-rewrite-author/internal/model"
)

// GitRewriteAuthor represents an active git-rewrite-author object
type GitRewriteAuthor struct {
	fl   model.Flags
	repo *git.Repo
}

// New creates new git-rewrite-author instance
func New(fl model.Flags) (*GitRewriteAuthor, error) {
	repo, err := git.Open(fl.Repo)
	return &GitRewriteAuthor{
		fl:   fl,
		repo: repo,
	}, err
}
