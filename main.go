package main

import (
	// Standard
	"flag"
	"github.com/e-valente/tullius/pkg/logging"
	"os"

	// 3rd Party
	"github.com/fatih/color"

	// Tullius
	"github.com/e-valente/tullius/pkg/banner"
	"github.com/e-valente/tullius/pkg/cli"
)

// Global Variables
var build = "NewHorizons01"

func main() {
	logging.Server("Starting Tullius Server ")

	flag.Usage = func() {
		color.Blue("#################################################")
		color.Blue("#\t\tWelcome to Tullius\t\t\t#")
		color.Blue("#################################################")
		color.Blue("Version: 0.1 " )
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	color.Blue(banner.TulliusBanner1)
	color.Blue("\t\t   Version: 0.1")
	color.Blue("\t\t   Build: %s", build)

	// Start Tullius Command Line Interface
	cli.Shell()
}

