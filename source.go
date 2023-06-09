package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/* ANON FILES API STRUCTS */
type APIResponse struct {
	Status bool `json:"status"`
	Data   struct {
		File struct {
			URL struct {
				Full  string `json:"full"`
				Short string `json:"short"`
			} `json:"url"`
			Metadata struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Size struct {
					Bytes    int    `json:"bytes"`
					Readable string `json:"readable"`
				} `json:"size"`
			} `json:"metadata"`
		} `json:"file"`
	} `json:"data"`
	Error *APIError `json:"error"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}


func main() {
	fmt.Println("\033[1;32mWelcome to the SkyNet Terminal!\033[0m")
	fmt.Println("\033[1;32mVersion: \033[90m1.0.1\033[0m")
	/* TODO: Add version checker & auto updater. */

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
		case "cd":
			cd(args)
		case "ls", ".":
			ls()
		case "help":
			if len(args) == 1 {
				printHelp()
			} else {
				printCommandHelp(args[1])
			}
		case "hostfile":
			if len(args) > 1 {
				hostFile(args[1])
			} else {
				fmt.Println("\033[1;31mCommand > hostfile: Please provide a file name.\033[0m")
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

func cd(args []string) {
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
	fmt.Println("- help <command>\tDisplays information about a specified command.")
	fmt.Println("- cd <directory>\tChanges or displays the current directory.")
	fmt.Println("- ping <host>\t\tPings a specified host.")
	fmt.Println("- echo <message>\tPrints anything.")
	fmt.Println("- ls\t\t\tDisplays a list of files/folders inside the current working directory.")
	fmt.Println("- self\t\t\tDisplays information about your IP Address and Operating System.")
	fmt.Println("- clear / cls\t\tClears the screen from previous commands.")
	fmt.Println("- hostfile <file>\tHosts the given file from the current directory on Anonfiles and retrieves the link.")
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
		"ping":    "The 'ping' command is used to ping a specified host. \nTakes (1) argument, 'host' \n\nExample: \n>>> ping example.com \n\nReturns: \nPinging example.com ...",
		"cd":      "The 'cd' command is used to change or display the current working directory. \nTakes (1) argument, optional: 'directory' \n\nExample: \n>>> cd Desktop/ \n>>> cd \n\nReturns: Desktop/",
		"ls":      "The 'ls' command is used to display a list of files/folders that are inside the current working directory. \nTakes no additional arguments.",
		"hostfile": "The 'hostfile' command is used to host the given file from the current directory on Anonfiles and retrieves the link. \nTakes (1) argument, 'file' \n\nExample: \n>>> hostfile myfile.txt",
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

func hostFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to open file '%s'. Error: %s\033[0m\n", fileName, err.Error())
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(fileName))
	if err != nil {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to create form file. Error: %s\033[0m\n", err.Error())
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to copy file data. Error: %s\033[0m\n", err.Error())
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to close multipart writer. Error: %s\033[0m\n", err.Error())
		return
	}

	req, err := http.NewRequest("POST", "https://api.anonfiles.com/upload", body)
	if err != nil {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to create request. Error: %s\033[0m\n", err.Error())
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to upload file. Error: %s\033[0m\n", err.Error())
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to read response body. Error: %s\033[0m\n", err.Error())
		return
	}

	var responseData APIResponse
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to parse response body. Error: %s\033[0m\n", err.Error())
		return
	}

	if responseData.Status {
		fmt.Println("File uploaded.")
		fmt.Printf("URL: %s\n", responseData.Data.File.URL.Full)
	} else {
		fmt.Printf("\033[1;31mCommand > hostfile: Failed to host file. Error: %s\033[0m\n", responseData.Error.Message)
	}
}



