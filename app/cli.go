package app

import (
	"context"

	"github.com/spf13/cobra"
)

var _ Runner = &CliRunner{}

// cli 是一种特殊的runner
type CliRunner struct {
	app  *Application
	root *cobra.Command
}

func NewCliRunner(app *Application) *CliRunner {
	return &CliRunner{
		app: app,
	}
}

// WithCommand 添加Command
// 命令层级最多允许一层，也就是只允许一个子命令
func (c *CliRunner) WithCommand(cmd *cobra.Command) *CliRunner {
	if c.root == nil {
		c.root = cmd
	} else {
		c.root.AddCommand(cmd)
	}

	return c
}

func (c *CliRunner) Name() string {
	return c.app.name
}

func (c *CliRunner) Run(ctx context.Context) error {
	return c.root.Execute()
}

func (c *CliRunner) Reload(ctx context.Context) error {
	return nil
}

func (c *CliRunner) Exit(ctx context.Context) error {
	return nil
}
