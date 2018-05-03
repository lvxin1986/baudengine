package server

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/urfave/cli.v2"

	"github.com/tiglabs/baudengine/util/build"
	"github.com/tiglabs/baudengine/util/log"
	"github.com/tiglabs/baudengine/util/multierror"
	"github.com/tiglabs/baudengine/util/routine"
)

var goFlags []*flag.Flag

type stopHook func() error

// VersionCommand return version sub command define
func VersionCommand() *cli.Command {
	return &cli.Command{
		Name:        "version",
		Usage:       "do the version",
		Description: "Prints out build version information",
		Action: func(c *cli.Context) error {
			fmt.Print(build.GetInfo())
			return nil
		},
	}
}

// AppendFlags append flag to command
func AppendFlags(cmd *cli.Command, flags ...cli.Flag) {
	cmd.Flags = append(cmd.Flags, flags...)
}

// AddGoFlags adds all command line flags to the app command
func AddGoFlags(cmd *cli.Command) {
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		goFlags = append(goFlags, gf)
		cmd.Flags = append(cmd.Flags, &cli.StringFlag{
			Name:        gf.Name,
			Value:       gf.Value.String(),
			Usage:       gf.Usage,
			DefaultText: gf.DefValue,
		})
	})
}

// SetGoFlagVals sets all command line flags value
func SetGoFlagVals(ctx *cli.Context) {
	for _, gf := range goFlags {
		gf.Value.Set(ctx.String(gf.Name))
	}

	goFlags = nil
}

// WaitShutdown awaits for Kill or SIGINT or SIGTERM and shutdown the server.
func WaitShutdown(stops ...stopHook) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs
	log.Info("Server Starting Shutdown...")

	merr := &multierror.MultiError{}
	for _, stop := range stops {
		merr.Append(stop())
	}
	if err := merr.ErrorOrNil(); err != nil {
		log.Error("Server Shutdown Error...", "Error", err)
	}

	// stop routine worker
	if err := routine.Stop(); err != nil {
		log.Error("Server Stop routine-worker error...", "Error", err)
	}

	log.Info("Server Shutdown And Exited...")
}

// SupressGlogWarnings is a hack to make flag.Parsed return true such that glog is happy about the flags having been parsed.
func SupressGlogWarnings() {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	_ = fs.Parse([]string{})
	flag.CommandLine = fs
}
