package cmd

import (
	"fmt"
	"strings"

	"github.com/crazy-max/git-rewrite-author/git"
	. "github.com/crazy-max/git-rewrite-author/utils"
	"github.com/spf13/cobra"
)

var configGetCmd = &cobra.Command{
	Use:     "config-get",
	Args:    cobra.ExactArgs(0),
	Short:   "Get current user name and email from Git config",
	Example: AppName + ` config-get`,
	Run:     configGetCmdRun,
}

var configSetCmd = &cobra.Command{
	Use:     "config-set user.name user.email",
	Short:   "Set user name and email to Git config",
	Example: AppName + ` config-set "John Smith" "john.smith@domain.com"`,
	Args:    cobra.ExactArgs(2),
	Run:     configSetCmdRun,
	PostRun: configGetCmdRun,
}

func init() {
	RootCmd.AddCommand(configGetCmd, configSetCmd)
}

func configGetCmdRun(cmd *cobra.Command, args []string) {
	repo, err := git.Open(Dir)
	CheckIfError(err)

	repo.ReloadConfig()
	if userName, found := repo.Get("user.name"); found {
		fmt.Printf("user.name: %s\n", userName)
	} else {
		Warning("user.name key not found in Git config")
	}
	if userEmail, found := repo.Get("user.email"); found {
		fmt.Printf("user.email: %s\n", userEmail)
	} else {
		Warning("user.email key not found in Git config")
	}
}

func configSetCmdRun(cmd *cobra.Command, args []string) {
	repo, err := git.Open(Dir)
	CheckIfError(err)

	repo.Set("user.name", strings.TrimSpace(args[0]))
	repo.Set("user.email", strings.TrimSpace(args[1]))
}
