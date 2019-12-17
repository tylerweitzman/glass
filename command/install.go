package command

import (
	"fmt"
	"log"
	"github.com/tylerweitzman/cli"
	daemon "github.com/tylerweitzman/glass/glass-daemon"
)

type Install struct {
	*command
}

func NewInstall() *Install {
	return &Install{newCommand()}
}

func (c *Install) Name() string {
	return "install"
}

func (c *Install) Description() string {
	return fmt.Sprintf("Runs the glass-daemon executable with both install and start. It requires admin privileges on windows and linux.")
}

func (c *Install) Usage() string {
	return "Install and start the background service"
}

func (c *Install) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (c *Install) Action() func(ctx *cli.Context) {
	return c.command.Action(c.Run)
}

func (c *Install) Run(ctx *cli.Context) error {
	c.Println("Installing the Timeglass background service...")

	//attempt to install
	log.Printf("Installing Daemon")
	daemon.SimulateMain("install");

	//attempt to start
	log.Printf("Loadiong Daemon")
	daemon.SimulateMain("start");

	// cmd := exec.Command("glass-daemon", "install")
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout

	// err := cmd.Run()
	// if err != nil {
	// 	return errwrap.Wrapf(fmt.Sprintf("Failed to install Daemon: {{err}}"), err)
	// }

	// c.Println("Starting the Timeglass background service...")

	// //attempt to start
	// cmd = exec.Command("glass-daemon", "start")
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout
	// err = cmd.Run()
	// if err != nil {
	// 	return errwrap.Wrapf(fmt.Sprintf("Failed to start Daemon: {{err}}"), err)
	// }

	c.Println("Done!")
	return nil
}
