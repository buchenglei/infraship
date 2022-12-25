package config

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
)

func NewCliReader(rootCmd ...*cobra.Command) (io.ReadCloser, error) {
	cli := &cliArg{}
	if len(rootCmd) == 1 && rootCmd[0] != nil {
		cli.cmd = rootCmd[0]
	} else {
		cli.cmd = &cobra.Command{
			Use: "app",
		}
	}

	err := cli.execute()
	if err != nil {
		return nil, err
	}

	return cli.getReader()
}

type cliArg struct {
	cmd      *cobra.Command
	provider string

	filePath string
}

func (c *cliArg) execute() error {
	configCmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf", "setting"},
	}
	configCmd.PersistentFlags().StringVarP(&c.provider, "provider", "p", "file", "set config reader provider")
	configCmd.PersistentFlags().StringVarP(&c.filePath, "conf", "c", "./config.toml", "config file path")

	c.cmd.AddCommand(configCmd)
	return c.cmd.Execute()
}

func (c *cliArg) getReader() (io.ReadCloser, error) {
	switch c.provider {
	case "file":
		return NewFileReader(c.filePath)
	}

	return nil, errors.New("read empty config source")
}
