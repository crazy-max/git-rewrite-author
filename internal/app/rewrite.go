package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"

	"github.com/crazy-max/git-rewrite-author/internal/utl"
	"github.com/rs/zerolog/log"
)

type rewriteAuthor struct {
	Old         []string `json:"old"`
	CorrectName string   `json:"correct_name"`
	CorrectMail string   `json:"correct_mail"`
}

// Rewrite rewrites an author/committer in Git history
func (gra *GitRewriteAuthor) RewriteOne() {
	correctName, correctMail, err := utl.ParseAddress(gra.fl.Correct)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot parse new Git name and email")
	}

	gra.rewrite([]*rewriteAuthor{
		{
			Old:         []string{gra.fl.Old},
			CorrectName: correctName,
			CorrectMail: correctMail,
		},
	})
}

// RewriteList rewrites a list of authors/committers in Git history
func (gra *GitRewriteAuthor) RewriteList() {
	authorsFile, err := ioutil.ReadFile(gra.fl.AuthorsJSON)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot read authors JSON file")
	}

	var rewriteAuthors []*rewriteAuthor
	if err = json.Unmarshal(authorsFile, &rewriteAuthors); err != nil {
		log.Fatal().Err(err).Msg("Cannot unmarshal authors JSON")
	}

	gra.rewrite(rewriteAuthors)
}

func (gra *GitRewriteAuthor) rewrite(rewriteAuthors []*rewriteAuthor) {
	b, _ := json.MarshalIndent(rewriteAuthors, "", "  ")
	log.Debug().Msg(string(b))

	tpl, err := template.New("rewrite").Parse(`{{range $i, $rewriteAuthor := .}}
OLD_EMAILS_{{$i}}=({{range $j, $old := .Old}}{{if $j}} {{end}}"{{$old}}"{{end}})
CORRECT_NAME_{{$i}}="{{.CorrectName}}"
CORRECT_EMAIL_{{$i}}="{{.CorrectMail}}"
for OLD_EMAIL_{{$i}} in ${OLD_EMAILS_{{$i}}[@]}; do
	if [ "$GIT_COMMITTER_EMAIL" == "$OLD_EMAIL_{{$i}}" ]; then
		export GIT_COMMITTER_NAME="$CORRECT_NAME_{{$i}}"
		export GIT_COMMITTER_EMAIL="$CORRECT_EMAIL_{{$i}}"
	fi
	if [ "$GIT_AUTHOR_EMAIL" == "$OLD_EMAIL_{{$i}}" ]; then
		export GIT_AUTHOR_NAME="$CORRECT_NAME_{{$i}}"
		export GIT_AUTHOR_EMAIL="$CORRECT_EMAIL_{{$i}}"
	fi
done
{{end}}`)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot parse rewrite template")
	}

	var rewriteScript bytes.Buffer
	if err := tpl.Execute(&rewriteScript, rewriteAuthors); err != nil {
		log.Fatal().Err(err).Msg("Cannot execute rewrite template")
	}
	log.Debug().Msgf("Rewrite script: %s", rewriteScript.String())

	log.Info().Msg("Following authors/committers will be rewritten")
	for _, rewriteAuthor := range rewriteAuthors {
		log.Info().Msgf("%s => '%s <%s>'", rewriteAuthor.Old, rewriteAuthor.CorrectName, rewriteAuthor.CorrectMail)
	}

	err = gra.repo.FilterBranch("--env-filter", rewriteScript.String(), "--tag-name-filter", "cat", "-f", "--", "--all")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot rewrite authors")
	}
}
