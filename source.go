package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
    "runtime"
)

func main() {
	fmt.Println("\033[1;32mWelcome to the SkyNet Terminal!\033[0m")
    fmt.Println("\033[1;32mVersion: \033[90m1.0.0\033[0m")
    fmt.Println("")
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\033[1;36m>>> \033[0m")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")
		switch args[0] {
    case "":
		case "clear", "cls":
			clearScreen()
		case "self":
			printSelfInfo()
		case "echo":
			echo(args[1:])
		case "help":
			if len(args) == 1 {
				printHelp()
			} else {
				printCommandHelp(args[1])
			}
		case "exit":
			fmt.Println("\033[1;32mGoodbye!\033[0m")
			os.Exit(0)
		default:
			fmt.Println("\033[1;31mError: Command not found.\033[0m")
		}
	}
}

func clearScreen() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func printSelfInfo() {
	hostname, _ := os.Hostname()
	addrs, _ := net.LookupIP(hostname)
	fmt.Println("Hostname:", hostname)
	fmt.Println("IP address:", addrs[0])
	fmt.Println("Operating System:", runtime.GOOS)
}

func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func printHelp() {
	fmt.Println("List of commands:")
	fmt.Println("help\t\tDisplays a list of all commands and their descriptions.")
	fmt.Println("echo <message>\tPrints anything after the keyword 'echo'.")
	fmt.Println("clear/cls\tClears the screen from previous commands.")
	fmt.Println("self\t\tDisplays information about your IP Address and Operating System.")
	fmt.Println("exit\t\tExits the terminal.")
}

func printCommandHelp(command string) {
	helpMessages := map[string]string{
		"clear":   "The 'clear' command is used to clear the screen from previous commands. \nTakes no additional arguments.",
		"cls":     "The 'cls' command is used to clear the screen from previous commands. \nTakes no additional arguments.",
		"self":    "The 'self' command is used to display information about your IP Address and Operating System. \nTakes no additional arguments.",
		"echo":    "The 'echo' command is used to print custom messages. \nTakes (1) argument, 'message' \n\nExample: \n>>> echo This is an example \n\nReturns: \nThis is an example",
		"exit":    "The 'exit' command is used to exit the terminal. \nTakes no additional arguments.",
		"help":    "The 'help' command is used to display a list of all commands and their descriptions. \nTakes (1) argument, optional: 'command' \n\nExample: \n>>> help help \n\nReturns: \nThe 'help' command is used to displ...",
	}

	if helpMessage, ok := helpMessages[command]; ok {
		fmt.Println(helpMessage)
	} else {
		fmt.Println("\033[1;31mError: Command not found.\033[0m")
	}
}
