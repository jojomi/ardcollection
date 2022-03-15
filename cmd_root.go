package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jojomi/feeddownload/feeddownload"
	"github.com/spf13/cobra"
)

func getRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ardcollection",
		Run: handleRootCmd,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("feed ID and target folder arguments are required")
			}
			if len(args) < 2 {
				return errors.New("a target folder argument is required")
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolP("dry-run", "d", false, "just simulate, no downloads are executed")
	f.BoolP("verbose", "v", false, "print more status")
	f.IntP("max-count", "m", 10, "max items checked")

	return cmd
}

func handleRootCmd(cmd *cobra.Command, args []string) {
	env := EnvRoot{}
	err := env.ParseFrom(cmd, args)
	if err != nil {
		log.Fatal("could not parse command: " + err.Error())
	}
	err = handleRoot(env)
	if err != nil {
		panic(err)
	}
}

func handleRoot(env EnvRoot) error {
	var (
		existing       bool
		targetFilename string
	)

	targetDir := env.TargetDir
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Target dir %s does not exist. Please create it and restart.", targetDir)
		os.Exit(1)
	}

	episodes, err := getEpisodes(env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get episodes: %s", err.Error())
		os.Exit(1)
	}
	for _, episode := range episodes {
		targetFilename = filepath.Join(targetDir, feeddownload.FilenameFromTitle(episode.Title)+filepath.Ext(episode.RemoteURL))

		_, err := os.Stat(targetFilename)
		existing = !os.IsNotExist(err)

		if env.Verbose || !existing {
			fmt.Println("Episode title:", episode.Title)
		}

		err = feeddownload.HandleFile(episode.RemoteURL, targetFilename, env.DryRun)
		if err != nil {
			panic(err)
		}

		if !existing {
			fmt.Println("")
		}
	}

	fmt.Printf("checked %d episodes\n", len(episodes))

	return nil
}
