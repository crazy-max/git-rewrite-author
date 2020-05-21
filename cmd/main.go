package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/git-rewrite-author/internal/app"
	"github.com/crazy-max/git-rewrite-author/internal/logging"
	"github.com/crazy-max/git-rewrite-author/internal/model"
	"github.com/rs/zerolog/log"
)

var (
	gra     *app.GitRewriteAuthor
	cli     model.Cli
	version = "dev"
)

func main() {
	var err error

	// Parse command line
	kctx := kong.Parse(&cli,
		kong.Name("git-rewrite-author"),
		kong.Description(`Rewrite authors history of a Git repository with ease. More info: https://github.com/crazy-max/git-rewrite-author`),
		kong.UsageOnError(),
		kong.Vars{
			"version": fmt.Sprintf("%s", version),
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	// Logger
	logging.Configure(cli)

	// Init
	if gra, err = app.New(cli); err != nil {
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

	switch kctx.Command() {
	case "config-get":
		gra.ConfigGet()
	case "config-set <name> <email>":
		gra.ConfigSet()
	case "list":
		gra.List()
	case "rewrite <old> <correct>":
		gra.RewriteOne()
	case "rewrite-list <file>":
		gra.RewriteList()
	default:
		log.Fatal().Err(err).Msg("Unknown command")
	}
}
