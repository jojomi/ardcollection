package main

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
)

// EnvRoot encapsulates the environment for the CLI root handler.
type EnvRoot struct {
	ID        int
	TargetDir string

	MaxCount int

	DryRun  bool
	Verbose bool
}

// ParseFrom reads the state from a given cobra command and its args.
func (e *EnvRoot) ParseFrom(command *cobra.Command, args []string) error {
	var err error

	if len(args) < 1 {
		return errors.New("need collection ID as first arg")
	}
	e.ID, err = strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	if len(args) < 2 {
		return errors.New("need target dir as second arg")
	}
	e.TargetDir = args[1]

	e.MaxCount, err = command.Flags().GetInt("max-count")
	if err != nil {
		return err
	}

	e.DryRun, err = command.Flags().GetBool("dry-run")
	if err != nil {
		return err
	}
	e.Verbose, err = command.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	return nil
}
