package main

import (
	// Standard
	"flag"
	"os"

	// 3rd Party
	"github.com/fatih/color"

	// Tullius
	"github.com/New-Horizons-Team/tullius/pkg/banner"
	"github.com/New-Horizons-Team/tullius/pkg/cli"
	"github.com/New-Horizons-Team/tullius/pkg/logging"
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

