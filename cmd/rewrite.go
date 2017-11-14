package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/mail"
	"strings"

	"github.com/crazy-max/git-rewrite-author/git"
	. "github.com/crazy-max/git-rewrite-author/utils"
	"github.com/spf13/cobra"
)

type rewriteAuthor struct {
	Old     []string `json:"old"`
	Correct string   `json:"correct"`
}

type rewriteAuthors []rewriteAuthor

// rewriteCmd to rewrite authors / committers in Git history
var rewriteCmd = &cobra.Command{
	Use:     `rewrite "old.email" "correct.name <correct.email>"`,
	Short:   "Rewrite an author / committer in Git history",
	Example: AppName + ` rewrite "root@localhost" "John Smith <john.smith@domain.com>"`,
	Args:    cobra.ExactArgs(2),
	Run:     rewriteCmdRun,
}

var rewriteListCmd = &cobra.Command{
	Use:   `rewrite-list "path/to/authors.json"`,
	Short: "Rewrite a list of authors / committers in Git history",
	Long: `Example of authors.json file :

[
  {
    "old": [ "root@localhost", "noreply@github.com" ],
    "correct": "John Smith <john.smith@domain.com>"
  },
  {
    "old": [ "ohcrap@bad.com" ],
    "correct": "Foo Bar <foobar@users.noreply.github.com>"
  }
]`,
	Example: AppName + ` rewrite-list "~/authors.json"`,
	Args:    cobra.ExactArgs(1),
	Run:     rewriteListCmdRun,
}

func init() {
	RootCmd.AddCommand(rewriteCmd, rewriteListCmd)
}

func rewriteCmdRun(cmd *cobra.Command, args []string) {
	rewrite([]string{args[0]}, args[1])
}

func rewriteListCmdRun(cmd *cobra.Command, args []string) {
	var rewriteAuthors []rewriteAuthor

	authorsRaw, err := ioutil.ReadFile(args[0])
	CheckIfError(err)

	err = json.Unmarshal(authorsRaw, &rewriteAuthors)
	CheckIfError(err)

	for _, rewriteAuthor := range rewriteAuthors {
		rewrite(rewriteAuthor.Old, rewriteAuthor.Correct)
	}
}

func rewrite(olds []string, correct string) {
	for _, old := range olds {
		_, err := mail.ParseAddress(old)
		CheckIfError(err)
	}

	correctName, correctMail, err := ParseAddress(correct)
	CheckIfError(err)

	fmt.Printf("\nRewritting %s to '%s <%s>'...\n", strings.Join(olds, ", "), correctName, correctMail)

	repo, err := git.Open(Dir)
	CheckIfError(err)

	err = repo.FilterBranch("--env-filter",
		fmt.Sprintf(`OLD_EMAILS=(%s)
CORRECT_NAME="%s"
CORRECT_EMAIL="%s"
for OLD_EMAIL in ${OLD_EMAILS[@]}; do
	if [ "$GIT_COMMITTER_EMAIL" == "$OLD_EMAIL" ]; then
		export GIT_COMMITTER_NAME="$CORRECT_NAME"
		export GIT_COMMITTER_EMAIL="$CORRECT_EMAIL"
	fi
	if [ "$GIT_AUTHOR_EMAIL" == "$OLD_EMAIL" ]; then
		export GIT_AUTHOR_NAME="$CORRECT_NAME"
		export GIT_AUTHOR_EMAIL="$CORRECT_EMAIL"
	fi
done`, strings.Join(olds, " "), correctName, correctMail),
		"--tag-name-filter", "cat", "-f", "--", "--all",
	)
	CheckIfError(err)
}
