package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/crazy-max/git-rewrite-author/internal/app"
	"github.com/crazy-max/git-rewrite-author/internal/logging"
	"github.com/crazy-max/git-rewrite-author/internal/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	gra     *app.GitRewriteAuthor
	flags   model.Flags
	version = "dev"
)

func main() {
	var err error

	cmd := kingpin.New("git-rewrite-author", "Rewrite authors history of a Git repository with ease. Moreinfo: https://github.com/crazy-max/git-rewrite-author")
	cmd.Flag("repo", "Git repository path.").Default(".").StringVar(&flags.Repo)
	cmd.Flag("log-level", "Set log level.").Default(zerolog.InfoLevel.String()).StringVar(&flags.LogLevel)
	cmd.Flag("log-caller", "Enable to add file:line of the caller.").Default("false").BoolVar(&flags.LogCaller)
	cmd.UsageTemplate(kingpin.CompactUsageTemplate).Version(version).Author("CrazyMax")

	configGet := cmd.Command("config-get", "Get current user name and email from Git config.")

	configSet := cmd.Command("config-set", "Set user name and email to Git config.")
	configSet.Arg("name", "Git username").Required().StringVar(&flags.GitUsername)
	configSet.Arg("email", "Git email").Required().StringVar(&flags.GitEmail)

	list := cmd.Command("list", "Display all authors/committers.")

	rewrite := cmd.Command("rewrite", "Rewrite an author/committer in Git history.")
	rewrite.Arg("old", "Current email linked to Git author to rewrite").Required().StringVar(&flags.Old)
	rewrite.Arg("correct", "New Git name and email to set").Required().StringVar(&flags.Correct)

	rewriteList := cmd.Command("rewrite-list", "Rewrite a list of authors/committers in Git history.")
	rewriteList.Arg("file", "Authors JSON file").Required().StringVar(&flags.AuthorsJSON)

	_, _ = cmd.Parse(os.Args[1:])

	// Logger
	logging.Configure(flags)

	// Init
	if gra, err = app.New(flags); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize")
	}

	// Handle os signals
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-channel
		log.Warn().Msgf("Caught signal %v", sig)
		os.Exit(0)
	}()

	switch kingpin.MustParse(cmd.Parse(os.Args[1:])) {
	case configGet.FullCommand():
		gra.ConfigGet()
	case configSet.FullCommand():
		gra.ConfigSet()
	case list.FullCommand():
		gra.List()
	case rewrite.FullCommand():
		gra.RewriteOne()
	case rewriteList.FullCommand():
		gra.RewriteList()
	default:
		log.Fatal().Err(err).Msg("Unknown command")
	}
}
