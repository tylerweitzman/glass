package command

import (
	"fmt"
	"github.com/tylerweitzman/cli"
	daemon "github.com/tylerweitzman/glass/glass-daemon"
)

type Daemon struct {
	*command
}

func NewDaemon() *Daemon {
	return &Daemon{newCommand()}
}

func (c *Daemon) Name() string {
	return "daemon"
}

func (c *Daemon) Description() string {
	return fmt.Sprintf("Runs the glass-daemon server.")
}

func (c *Daemon) Usage() string {
	return "Start the background server"
}//.

func (c *Daemon) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (c *Daemon) Action() func(ctx *cli.Context) {
	return c.command.Action(c.Run)
}

func (c *Daemon) Run(ctx *cli.Context) error {
	c.Println("Running Daemon Server")
	daemon.SimulateMain("", ctx)
	return nil
}
