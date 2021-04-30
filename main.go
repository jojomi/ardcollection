package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/jojomi/feeddownload/feeddownload"
	"github.com/spf13/cobra"
)

var (
	flagRootDryRun bool
)

func main() {
	rootCmd := &cobra.Command{
		Use: "ardcollection",
		Run: handleRootCmd,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("feed url and target folder arguments are required")
			}
			if len(args) < 2 {
				return errors.New("a target folder argument is required")
			}
			return nil
		},
	}

	pFlags := rootCmd.PersistentFlags()
	pFlags.BoolVarP(&flagRootDryRun, "dry-run", "d", false, "just simulate, no downloads are executed")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleRootCmd(cmd *cobra.Command, args []string) {
	var (
		err            error
		remoteURL      string
		targetFilename string
	)

	// parse supplied feed
	resp, err := http.Get(args[0])
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var feed Feed
	err = json.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range feed.Result.Episodes {
		remoteURL = f.Enclosure.DownloadURL

		fmt.Println("")
		fmt.Println("Episode title:", f.Title)
		fmt.Println("Remote URL:", remoteURL)

		targetFilename = filepath.Join(args[1], feeddownload.FilenameFromTitle(f.Title)+path.Ext(remoteURL))
		fmt.Println("Local file:", targetFilename)

		err = feeddownload.HandleFile(remoteURL, targetFilename, flagRootDryRun)
		if err != nil {
			panic(err)
		}
	}
}
