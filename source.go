package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
    "os/exec"
	"strings"
    "runtime"
	"io/ioutil"

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
		case "ping":
		    ping(args[1:])
		case "self":
			printSelfInfo()
		case "echo":
			echo(args[1:])
		case "cd": // Needs more work.
			cd(args)
        case "ls", ".":
		    ls()
		case "hostfile":
		    if len(args) < 2 {
		    	fmt.Println("\033[1;31mCommand > HostFile: No file chosen.\033[0m")
            break
		    }

		    fileName := args[1]
		    hostfile(fileName)
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
			fmt.Println("\033[1;31mCommand Handler: Command not found.\033[0m")
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

func cd(args []string) { // NEEDS WORK
if len(args) > 1 {
err := os.Chdir(args[1])
if err != nil {
fmt.Println("\033[1;31mCommand > Cd:", err, "\033[0m")
}
} else {
currentDirectory, _ := os.Getwd()
fmt.Println(currentDirectory)
}
}

func printHelp() {
	fmt.Println("List of commands:\n")
	fmt.Println("- hostfile <file>\tHosts a specified file online on diffrent file-sharing services.")
	fmt.Println("- help <command>\tDisplays information about a specified command.")
    fmt.Println("- cd <directory>\tChanges or displays the current directory.")
    fmt.Println("- ping <host>\t\tPings a specified host.")
	fmt.Println("- echo <message>\tPrints anything.")
    fmt.Println("- ls\t\t\tDisplays a list of files/folders inside the current working directory.")
	fmt.Println("- self\t\t\tDisplays information about your IP Address and Operating System.")
	fmt.Println("- clear / cls\t\tClears the screen from previous commands.")
	fmt.Println("- exit\t\t\tExits the terminal.")
}

func printCommandHelp(command string) {
	helpMessages := map[string]string{
		"clear":   "The 'clear' command is used to clear the screen from previous commands. \nTakes no additional arguments.",
		"cls":     "The 'cls' command is used to clear the screen from previous commands. \nTakes no additional arguments.",
		"self":    "The 'self' command is used to display information about your IP Address and Operating System. \nTakes no additional arguments.",
		"echo":    "The 'echo' command is used to print custom messages. \nTakes (1) argument, 'message' \n\nExample: \n>>> echo This is an example \n\nReturns: \nThis is an example",
		"exit":    "The 'exit' command is used to exit the terminal. \nTakes no additional arguments.",
		"help":    "The 'help' command is used to display a list of all commands and their descriptions. \nTakes (1) argument, optional: 'command' \n\nExample: \n>>> help help \n\nReturns: \nThe 'help' command is used to displ...",
		"ping":    "The 'ping' command is used to ping a specified host. \nTakes (1) argument, 'host' \n\nExample: \n>>> ping example.com \n\nReturns: \nPinging example.com ...\n",
        "cd":      "The 'cd' command is used to change or display the current working directory. \nTakes (1) argument, optional: 'directory' \n\nExample: \n>>> cd Desktop/ \n>>> cd \n\nReturns: Desktop/",
        "ls":      "The 'ls' command is used to display a list of files/folders that are inside the current working directory. \nTakes no additional arguments.",
		"hostfile":    "The 'hostfile' command is used to host a file on diffrent file sharing sites. \nTakes (1) argument, 'file' \n\nExample: \n>>> hostfile example.txt \n\nReturns: \n[ API 1] (link) \n[ API 2 ] (link) \n...",
	}

	if helpMessage, ok := helpMessages[command]; ok {
		fmt.Println(helpMessage)
	} else {
		fmt.Println("\033[1;31mCommand > Help: Command not found.\033[0m")
	}
}

func ping(args []string) {
	if len(args) == 0 {
		fmt.Println("\033[1;31mCommand > Ping: Please specify a host to ping.\033[0m")
		return
	}
	host := args[0]
	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Println("\033[1;31mCommand > Ping: Host not found.\033[0m")
		return
	}
	addrs, err := net.LookupAddr(ips[0].String())
	if err != nil {
		fmt.Println("\033[1;31mCommand > Ping: Could not retrieve address information for host.\033[0m")
		return
	}
	fmt.Println("Pinged", host)
	fmt.Println("IP addresses: ", ips)
	fmt.Println("Reverse DNS: ", addrs)
}

func hostfile(file string) error {
	// Check if file exists in current directory
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("\033[1;31mCommand > HostFile: File not found. Make sure it's inside the current directory.\033[0m")
        return err
	}

    // API #1
	curl := exec.Command("curl", "--upload-file", file, "https://transfer.sh/"+file)
	curlOut, err := curl.CombinedOutput()
	if err != nil {
		return fmt.Errorf("\033[1;31mCommand > HostFile > Curl: %v\033[0m", err)
	}

	lines := strings.Split(string(curlOut), "\n")
	link := lines[len(lines)-1]

	fmt.Println("[ Transfer.sh API ]", link)

    // API #2
	curlFileio := exec.Command("curl", "-F", "file=@"+file, "https://file.io")
	curlFileioOut, err := curlFileio.CombinedOutput()
	if err != nil {
		return fmt.Errorf("\033[1;31mCommand > HostFile > Curl: %v\033[0m", err)
	}

	link = ""
	lines = strings.Split(string(curlFileioOut), "\n")
	for _, line := range lines {
		if strings.Contains(line, `"link"`) {
			parts := strings.Split(line, `":"`)
			link = strings.Trim(strings.TrimSpace(parts[1]), `","key`)
			break
		}
	}

	fmt.Println("[ File.io API ]", "https://file.io/"+link)

	return nil
}


func ls() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println("\033[1;31mCommand > ls:", err, "\033[0m")
	} else {
		for _, file := range files {
			fmt.Println(file.Name())
		}
	}
}
