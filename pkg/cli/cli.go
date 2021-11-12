package cli

import (
	"bufio"
	"fmt"
	"github.com/e-valente/tullius/pkg/api/messages"
	"github.com/e-valente/tullius/pkg/banner"
	"github.com/e-valente/tullius/pkg/logging"
	awsscan "github.com/e-valente/tullius/pkg/modules/aws"
	"github.com/e-valente/tullius/pkg/modules/network"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	// 3rd Party
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/mattn/go-shellwords"
	"github.com/olekukonko/tablewriter"
)

// Global Variables
var prompt *readline.Instance
var shellCompleter *readline.PrefixCompleter
var shellMenuContext = "main"

// MessageChannel is used to input user messages that are eventually written to STDOUT on the CLI application
var MessageChannel = make(chan messages.UserMessage)
var clientID = uuid.NewV4()

// Shell is the exported function to start the command line interface
func Shell() {

	printUserMessage()
	registerMessageChannel()
	getUserMessages()

	p, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31mTullius»\033[0m ",
		HistoryFile:         "/tmp/readline.tmp",
		AutoComplete:        shellCompleter,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})

	if err != nil {
		MessageChannel <- messages.UserMessage{
			Level:   messages.Warn,
			Message: fmt.Sprintf("There was an error with the provided input: %s", err.Error()),
			Time:    time.Now().UTC(),
			Error:   true,
		}
	}
	prompt = p

	defer func() {
		err := prompt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.SetOutput(prompt.Stderr())

	for {
		line, err := prompt.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			exit()
		}

		line = strings.TrimSpace(line)
		//cmd := strings.Fields(line)
		cmd, err := shellwords.Parse(line)
		if err != nil {
			MessageChannel <- messages.UserMessage{
				Level:   messages.Warn,
				Message: fmt.Sprintf("error parsing command line arguments:\r\n%s", err),
				Time:    time.Now().UTC(),
				Error:   false,
			}
		}

		if len(cmd) > 0 {
			switch shellMenuContext {
			case "main":
				switch cmd[0] {
				case "banner":
					m := "\n"
					m += color.BlueString(banner.TulliusBanner1)
					m += color.BlueString("\r\n\t\t   Version: 1.0")
					m += color.BlueString("\r\n\t\t   Build: 1.0.0\n")
					MessageChannel <- messages.UserMessage{
						Level:   messages.Plain,
						Message: m,
						Time:    time.Now().UTC(),
						Error:   false,
					}
				case "network":
					menuNetwork()
				case "aws":
					menuAWS()
				case "use":
					menuUse(cmd[1:])
				case "help":
					menuHelpMain()
				case "?":
					menuHelpMain()
				case "exit", "quit":
					if len(cmd) > 1 {
						if strings.ToLower(cmd[1]) == "-y" {
							exit()
						}
					}
					if confirm("Are you sure you want to exit?") {
						exit()
					}
				default:
					if len(cmd) > 1 {
						executeCommand(cmd[0], cmd[1:])
					} else {
						executeCommand(cmd[0], []string{})
					}
				}


			case "network":
				switch cmd[0] {
					case "info", "help", "list":
						menuNetwork()
					case "back":
						menuSetMain()
					case "net-scan":
						network.NetworkScan()
					default:
						if len(cmd) > 1 {
							executeCommand(cmd[0], cmd[1:])
						} else {
							executeCommand(cmd[0], []string{})
						}
				}
			case "aws":
				switch cmd[0] {
					case "info", "help", "list":
						menuAWS()
					case "back":
						menuSetMain()
					case "s3-bucket-scan":
						awsscan.AWSS3ScanBucket()
					case "s3-object-scan":
						awsscan.AWSS3ScanObjects()
					default:
						if len(cmd) > 1 {
							executeCommand(cmd[0], cmd[1:])
						} else {
							executeCommand(cmd[0], []string{})
						}
				}

			}

		}

	}
}

// printUserMessage is used to print all messages to STDOUT for command line clients
func printUserMessage() {
	go func() {
		for {
			m := <-MessageChannel
			switch m.Level {
			case messages.Info:
				fmt.Println(color.CyanString("\n[i] %s", m.Message))
			case messages.Note:
				fmt.Println(color.YellowString("\n[-] %s", m.Message))
			case messages.Warn:
				fmt.Println(color.RedString("\n[!] %s", m.Message))
			case messages.Success:
				fmt.Println(color.GreenString("\n[+] %s", m.Message))
			case messages.Plain:
				fmt.Println("\n" + m.Message)
			default:
				fmt.Println(color.RedString("\n[_-_] Invalid message level: %d\r\n%s", m.Level, m.Message))
			}
		}
	}()
}

func registerMessageChannel() {
	um := messages.Register(clientID)
	if um.Error {
		MessageChannel <- um
		return
	}

}

func getUserMessages() {
	go func() {
		for {
			MessageChannel <- messages.GetMessageForClient(clientID)
		}
	}()
}


func menuHelpMain() {
	MessageChannel <- messages.UserMessage{
		Level:   messages.Plain,
		Message: color.YellowString("Tullius (version 0.0.1)\n"),
		Time:    time.Now().UTC(),
		Error:   false,
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCaption(true, "Main Menu Help")
	table.SetHeader([]string{"Command", "Description", "Options"})

	data := [][]string{
		{"use module", "Use a tullius module", "network|aws|gcp|k8s|vault"},
		{"banner", "Print the tullius banner", ""},
		{"exit", "Exit and close the tullius server", ""},
		{"quit", "Exit and close the tullius server", ""},
		{"version", "Print the tullius server version", ""},
		{"*", "Anything else is executed on the host operating system", ""},
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

func menuSetMain() {
	prompt.SetPrompt("\033[31mTullius»\033[0m ")
	shellMenuContext = "main"
}

func menuNetwork() {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCaption(true, "Network Module")
	table.SetHeader([]string{"Command", "Description", "Options"})

	prompt.SetPrompt("\033[31mTullius[\033[32mmodule\033[31m][\033[33m" + "network" + "\033[31m]»\033[0m ")
	shellMenuContext = "network"

	data := [][]string{
		{"net-scan", "Network Scan", ""},
		{"*", "Anything else is executed on the host operating system", ""},
		{"back", "Back to the main menu", ""},
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

func menuAWS() {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCaption(true, "AWS Module")
	table.SetHeader([]string{"Command", "Description", "Options"})

	prompt.SetPrompt("\033[31mTullius[\033[32mmodule\033[31m][\033[33m" + "aws" + "\033[31m]»\033[0m ")
	shellMenuContext = "aws"

	data := [][]string{
		{"s3-bucket-scan", "Scan S3 service for buckets", "--wordlist <path|rockyou"},
		{"s3-object-scan", "Scan S3 bucket for objects", "bucket <bucketName> --wordlist <path|gobuster"},
		{"*", "Anything else is executed on the host operating system", ""},
		{"back", "Back to the main menu", ""},
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}



// confirm reads in a string and returns true if the string is y or yes but does not provide the prompt question
func confirm(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	//fmt.Print(color.RedString(fmt.Sprintf("%s [yes/NO]: ", question)))
	MessageChannel <- messages.UserMessage{
		Level:   messages.Plain,
		Message: color.RedString(fmt.Sprintf("%s [yes/NO]: ", question)),
		Time:    time.Now().UTC(),
		Error:   false,
	}
	response, err := reader.ReadString('\n')
	if err != nil {
		MessageChannel <- messages.UserMessage{
			Level:   messages.Warn,
			Message: fmt.Sprintf("There was an error reading the input:\r\n%s", err.Error()),
			Time:    time.Now().UTC(),
			Error:   true,
		}
	}
	response = strings.ToLower(response)
	response = strings.Trim(response, "\r\n")
	yes := []string{"y", "yes", "-y", "-Y"}

	for _, match := range yes {
		if response == match {
			return true
		}
	}
	return false
}

// exit will prompt the user to confirm if they want to exit
func exit() {
	color.Red("[!]Quitting...")
	logging.Server("Shutting down tullius due to user input")
	os.Exit(0)
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func executeCommand(name string, arg []string) {

	cmd := exec.Command(name, arg...) // #nosec G204 Users can execute any arbitrary command by design

	out, err := cmd.CombinedOutput()

	MessageChannel <- messages.UserMessage{
		Level:   messages.Info,
		Message: "Executing system command...",
		Time:    time.Time{},
		Error:   false,
	}
	if err != nil {
		MessageChannel <- messages.UserMessage{
			Level:   messages.Warn,
			Message: err.Error(),
			Time:    time.Time{},
			Error:   true,
		}
	} else {
		MessageChannel <- messages.UserMessage{
			Level:   messages.Success,
			Message: fmt.Sprintf("%s", out),
			Time:    time.Time{},
			Error:   false,
		}
	}
}

func menuUse(cmd []string) {
	if len(cmd) > 0 {
		switch cmd[0] {
		case "module":
			if len(cmd) > 1 {
				menuSetModule(cmd[1])
			} else {
				MessageChannel <- messages.UserMessage{
					Level:   messages.Warn,
					Message: "Invalid module",
					Time:    time.Now().UTC(),
					Error:   false,
				}
			}
		case "":
		default:
			MessageChannel <- messages.UserMessage{
				Level:   messages.Note,
				Message: "Invalid 'use' command",
				Time:    time.Now().UTC(),
				Error:   false,
			}
		}
	} else {
		MessageChannel <- messages.UserMessage{
			Level:   messages.Note,
			Message: "Invalid 'use' command",
			Time:    time.Now().UTC(),
			Error:   false,
		}
	}
}

func menuSetModule(cmd string) {
	if len(cmd) > 0 {
		switch cmd {
		case "network":
			prompt.SetPrompt("\033[31mTullius[\033[32mmodule\033[31m][\033[33m" + "network" + "\033[31m]»\033[0m ")
			shellMenuContext = "network"
			menuNetwork()
		case "aws":
			prompt.SetPrompt("\033[31mTullius[\033[32mmodule\033[31m][\033[33m" + "aws" + "\033[31m]»\033[0m ")
			shellMenuContext = "aws"
			menuAWS()
		case "k8s":
			prompt.SetPrompt("\033[31mTullius[\033[32mmodule\033[31m][\033[33m" + "k8s" + "\033[31m]»\033[0m ")
			shellMenuContext = "k8s"
		}

	}
}


