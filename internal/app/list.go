package app

import (
	"github.com/rs/zerolog/log"
)

// List displays all authors/committers
func (gra *GitRewriteAuthor) List() {
	logs, err := gra.repo.Logs()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot fetch git logs")
	}

	log.Debug().Msg("Seeking authors...")
	for _, author := range logs.GetAuthors() {
		log.Info().Msgf("%s <%s>", author.Name, author.Email)
	}
}
