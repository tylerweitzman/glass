package daemon

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"github.com/hashicorp/errwrap"
	"github.com/tylerweitzman/service"
	"github.com/tylerweitzman/cli"
)

var Version = "0.0.0"
var Build = "gobuild"

type daemon struct {
	keeper *Keeper
	server *Server
}

func (p *daemon) Start(s service.Service) error {
	var err error

	path, err := SystemTimeglassPathCreateIfNotExist()
	if err != nil {
		return errwrap.Wrapf("Failed to find Timeglass system path: {{err}}", err)
	}

	p.keeper, err = NewKeeper(path)
	if err != nil {
		return errwrap.Wrapf("Failed to create time keeper: {{err}}", err)
	}

	bind := "127.0.0.1:3838"
	p.server, err = NewServer(bind, p.keeper)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to create server on '%s': {{err}}, is the service already running?", bind), err)
	}

	go p.server.checkVersion()
	go p.keeper.Start()
	go p.run()
	return nil
}

func (p *daemon) Stop(s service.Service) error {
	p.keeper.Stop()
	return p.server.Stop()
}

func (p *daemon) run() error {
	return p.server.Start()
}

func SimulateMain(command string, ctx *cli.Context) {

	//setup logging to a file
	silent := ctx.GlobalBool("silent")
	var out = os.Stdout
	if silent {
		out = nil
	}
	l, err := NewLogger(out)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	log.SetOutput(l)
	defer l.Close()
	
	//initialize service
	conf := &service.Config{
		Name:        "panda-glass",
		DisplayName: "panda-glass",
		Description: "Helper for panda git commits linting metadata",
		Arguments: []string{"daemon"},
		Option:      map[string]interface{}{},
	}

	if runtime.GOOS == "darwin" {
		conf.Name = "com.tylerweitzman.panda-glass"
		//let's addssssssss
		// results in a automated installer that
		// has some serious usability issues
		conf.Option["UserService"] = true
		conf.Option["RunAtLoad"] = true
		pathOut, pathErr, logCreateError := DarwinTimeglassLogPaths();
		if logCreateError == nil {
			conf.Option["LogOutput"] = false
			conf.Option["StandardOutPath"] = pathOut
			conf.Option["StandardErrorPath"] = pathErr
			conf.Option["StdErrPath"] = pathErr
		} else {
			log.Printf("Daemon could not create log paths '%s'", err);
		}
	} else if runtime.GOOS == "windows" {

		//WATCH OUT: timeglass has a windows installer
		//that takes care of installing and starting service.
		//in addition to the command line

		conf.Name = "Panda-Glass" //windows style
	}

	d := &daemon{}
	s, err := service.New(d, conf)
	if err != nil {
		log.Fatal(err)
	}

	//handle service controls
	if len(command) > 0 {
		err = service.Control(s, command)
		if err != nil && !silent {
			ReportServiceControlErrors(err)
		}
		return
	}

	//start daemon
	log.Printf("Daemon launched, writing logs to '%s'", l.Path())
	defer func() {
		log.Printf("Daemon terminated\n\n")
	}()

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}


func main() {
	flag.Parse()

	//setup logging to a file
	l, err := NewLogger(os.Stderr)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	log.SetOutput(l)
	defer l.Close()

	//initialize service
	conf := &service.Config{
		Name:        "timeglass",
		DisplayName: "Timeglass",
		Description: "Automated time tracking daemon that monitors file changes",
		Option:      map[string]interface{}{},
	}

	if runtime.GOOS == "darwin" {
		conf.Name = "com.timeglass.glass-daemon"
		//let's addssssssss
		// results in a automated installer that
		// has some serious usability issues
		conf.Option["UserService"] = true
		conf.Option["RunAtLoad"] = true
		pathOut, pathErr, err := DarwinTimeglassLogPaths();
		if(err!=nil) {
			conf.Option["StandardOutPath"] = pathOut
			conf.Option["StandardErrorPath"] = pathErr
		}
	} else if runtime.GOOS == "windows" {

		//WATCH OUT: timeglass has a windows installer
		//that takes care of installing and starting services
		//in addition to the command line

		conf.Name = "Timeglass" //windows style
	}

	d := &daemon{}
	s, err := service.New(d, conf)
	if err != nil {
		log.Fatal(err)
	}

	//handle service controls
	if len(flag.Args()) > 0 {
		err = service.Control(s, flag.Args()[0])
		if err != nil {
			ReportServiceControlErrors(err)
		}
		return
	}

	//start daemon
	log.Printf("Daemon launched, writing logs to '%s'", l.Path())
	defer func() {
		log.Printf("Daemon terminated\n\n")
	}()

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
