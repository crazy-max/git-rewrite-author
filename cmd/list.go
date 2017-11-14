package cmd

import (
	"fmt"

	"github.com/crazy-max/git-rewrite-author/git"
	. "github.com/crazy-max/git-rewrite-author/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Display all authors / committers",
	Example: AppName + ` list`,
	Args:    cobra.ExactArgs(0),
	Run:     listRun,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func listRun(cmd *cobra.Command, args []string) {
	repo, err := git.Open(Dir)
	CheckIfError(err)

	l, err := repo.Logs()
	CheckIfError(err)

	for _, author := range l.GetAuthors() {
		fmt.Printf("%s <%s>\n", author.Name, author.Email)
	}
}
