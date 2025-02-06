package cmd

import (
	"flag"
	"fmt"
	"os"
)

func Execute() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	connectCmd := flag.NewFlagSet("connect", flag.ExitOnError)

	port := startCmd.String("port", "8080", "Port to run the server")

	serverAddr := connectCmd.String("server", "localhost:8080", "Server address")
	username := connectCmd.String("username", "", "Your username")
	token := connectCmd.String("token", "", "Authentication token")

	if len(os.Args) < 1 {
		fmt.Println("Usage: broadcast-server [start|connect]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		StartServer(*port)
	case "connect":
		connectCmd.Parse(os.Args[2:])
		if *username == "" || *token == "" {
			fmt.Println("Error: Username and token are required for authentication")
			os.Exit(1)
		}
		ConnectClient(*serverAddr, *username, *token)
	default:
		fmt.Println("Unknown command. Use 'start' or 'connect'")
		os.Exit(1)
	}
}
