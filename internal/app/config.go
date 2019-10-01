package app

import (
	"github.com/rs/zerolog/log"
)

// ConfigGet gets current user name and email from Git config
func (gra *GitRewriteAuthor) ConfigGet() {
	gra.repo.ReloadConfig()

	if name, found, err := gra.repo.Get("user.name"); found {
		log.Info().Msgf("user.name:  %s", name)
	} else if err != nil {
		log.Fatal().Err(err).Msgf("Cannot retrieve user.email in Git config")
	} else {
		log.Warn().Msg("user.name key not found in Git config")
	}

	if email, found, err := gra.repo.Get("user.email"); found {
		log.Info().Msgf("user.email: %s", email)
	} else if err != nil {
		log.Fatal().Err(err).Msgf("Cannot retrieve user.email in Git config")
	} else {
		log.Warn().Msg("user.email key not found in Git config")
	}
}

// ConfigSet sets user name and email to Git config
func (gra *GitRewriteAuthor) ConfigSet() {
	if err := gra.repo.Set("user.name", gra.fl.GitUsername); err != nil {
		log.Fatal().Err(err).Msg("Cannot set user.name")
	}
	if err := gra.repo.Set("user.email", gra.fl.GitEmail); err != nil {
		log.Fatal().Err(err).Msg("Cannot set user.email")
	}
}
