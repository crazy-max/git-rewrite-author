package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/crazy-max/git-rewrite-author/git"
	. "github.com/crazy-max/git-rewrite-author/utils"
	"github.com/spf13/cobra"
)

type rewriteAuthor struct {
	Old         []string `json:"old"`
	CorrectName string   `json:"correct_name"`
	CorrectMail string   `json:"correct_mail"`
}

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
		"correct_name": "John Smith",
		"correct_mail": "john.smith@domain.com"
	},
 	{
		"old": [ "ohcrap@bad.com" ],
		"correct_name": "Good Sir",
		"correct_mail": "goodsir@users.noreply.github.com"
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
	var rewriteAuthors []*rewriteAuthor

	correctName, correctMail, err := ParseAddress(args[1])
	CheckIfError(err)

	rewriteAuthors = append(rewriteAuthors, &rewriteAuthor{
		Old:         []string{args[0]},
		CorrectName: correctName,
		CorrectMail: correctMail,
	})

	rewrite(rewriteAuthors)
}

func rewriteListCmdRun(cmd *cobra.Command, args []string) {
	var rewriteAuthors []*rewriteAuthor

	authorsRaw, err := ioutil.ReadFile(args[0])
	CheckIfError(err)

	err = json.Unmarshal(authorsRaw, &rewriteAuthors)
	CheckIfError(err)

	rewrite(rewriteAuthors)
}

func rewrite(rewriteAuthors []*rewriteAuthor) {
	if DebugEnabled {
		prettyRewriteAuthors, _ := json.MarshalIndent(rewriteAuthors, "", "  ")
		Debug("\n### Author(s) to rewrite:\n")
		Debug(string(prettyRewriteAuthors))
	}

	t := template.New("rewrite")
	t, _ = t.Parse(`{{range $i, $rewriteAuthor := .}}
OLD_EMAILS_{{$i}}=({{range $j, $old := .Old}}{{if $j}} {{end}}"{{$old}}"{{end}})
CORRECT_NAME_{{$i}}="{{.CorrectName}}"
CORRECT_EMAIL_{{$i}}="{{.CorrectMail}}"
for OLD_EMAIL_{{$i}} in ${OLD_EMAILS_{{$i}}[@]}; do
	if [ "$GIT_COMMITTER_EMAIL" = "$OLD_EMAIL_{{$i}}" ]; then
		GIT_COMMITTER_NAME="$CORRECT_NAME_{{$i}}"
		GIT_COMMITTER_EMAIL="$CORRECT_EMAIL_{{$i}}"
	fi
	if [ "$GIT_AUTHOR_EMAIL" = "$OLD_EMAIL_{{$i}}" ]; then
		GIT_AUTHOR_NAME="$CORRECT_NAME_{{$i}}"
		GIT_AUTHOR_EMAIL="$CORRECT_EMAIL_{{$i}}"
	fi
done
{{end}}`)

	var tpl bytes.Buffer
	err := t.Execute(&tpl, rewriteAuthors)
	CheckIfError(err)

	if DebugEnabled {
		Debug("\n### Template used:\n%s", tpl.String())
	}

	fmt.Printf("\nFollowing authors / committers will be rewritten :")
	for _, rewriteAuthor := range rewriteAuthors {
		fmt.Printf("\n- %s => '%s <%s>'",
			strings.Join(rewriteAuthor.Old, ", "),
			rewriteAuthor.CorrectName,
			rewriteAuthor.CorrectMail,
		)
	}
	fmt.Printf("\n\n")

	repo, err := git.Open(Dir)
	CheckIfError(err)

	err = repo.FilterBranch("--env-filter", tpl.String(), "--tag-name-filter", "cat", "--", "--all")
	CheckIfError(err)
}
