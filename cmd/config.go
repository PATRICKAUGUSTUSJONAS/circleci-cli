package cmd

import (
	"context"

	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/filetree"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

const defaultConfigPath = ".circleci/config.yml"

func newConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Operate on build config files",
	}

	collapseCommand := &cobra.Command{
		Use:   "collapse [path]",
		Short: "Collapse your CircleCI configuration to a single file",
		RunE:  collapseConfig,
		Args:  cobra.MaximumNArgs(1),
	}

	validateCommand := &cobra.Command{
		Use:     "validate [config.yml]",
		Aliases: []string{"check"},
		Short:   "Check that the config file is well formed.",
		RunE:    validateConfig,
		Args:    cobra.MaximumNArgs(1),
	}

	expandCommand := &cobra.Command{
		Use:   "expand [config.yml]",
		Short: "Expand the config.",
		RunE:  expandConfig,
		Args:  cobra.MaximumNArgs(1),
	}

	configCmd.AddCommand(collapseCommand)
	configCmd.AddCommand(validateCommand)
	configCmd.AddCommand(expandCommand)

	return configCmd
}

func validateConfig(cmd *cobra.Command, args []string) error {
	configPath := defaultConfigPath
	if len(args) == 1 {
		configPath = args[0]
	}
	ctx := context.Background()
	response, err := api.ConfigQuery(ctx, Logger, configPath)

	if err != nil {
		return err
	}

	if !response.Valid {
		return response.ToError()
	}

	Logger.Infof("Config file at %s is valid", configPath)
	return nil
}

func expandConfig(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	configPath := defaultConfigPath
	if len(args) == 1 {
		configPath = args[0]
	}
	response, err := api.ConfigQuery(ctx, Logger, configPath)

	if err != nil {
		return err
	}

	if !response.Valid {
		return response.ToError()
	}

	Logger.Info(response.OutputYaml)
	return nil
}

func collapseConfig(cmd *cobra.Command, args []string) error {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	tree, err := filetree.NewTree(root)
	if err != nil {
		return errors.Wrap(err, "An error occurred trying to build the tree")
	}

	y, err := yaml.Marshal(&tree)
	if err != nil {
		return errors.Wrap(err, "Failed trying to marshal the tree to YAML ")
	}
	Logger.Infof("%s\n", string(y))
	return nil
}
