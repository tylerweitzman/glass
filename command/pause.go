package command

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/errwrap"

	"github.com/tylerweitzman/glass/vcs"
)

type Pause struct {
	*command
}

func NewPause() *Pause {
	return &Pause{newCommand()}
}

func (c *Pause) Name() string {
	return "pause"
}

func (c *Pause) Description() string {
	return fmt.Sprintf("Pauses the timer, running 'glass start' or editing a file in the repository resumes the timer")
}

func (c *Pause) Usage() string {
	return "Manually Pause the timer, go get some coffee"
}

func (c *Pause) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (c *Pause) Action() func(ctx *cli.Context) {
	return c.command.Action(c.Run)
}

func (c *Pause) Run(ctx *cli.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return errwrap.Wrapf("Failed to fetch current working dir: {{err}}", err)
	}

	vc, err := vcs.GetVCS(dir)
	if err != nil {
		return errwrap.Wrapf("Failed to setup VCS: {{err}}", err)
	}

	c.Printf("Pausing timer...")

	client := NewClient()
	err = client.PauseTimer(vc.Root())
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to pause timer: {{err}}"), err)
	}

	c.Printf("Done!")
	return nil
}
