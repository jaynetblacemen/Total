package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Username string `json:"username"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

type Friends map[string]string

func main() {
	baseDir := getBaseDir()
	configPath := filepath.Join(baseDir, "config.json")
	friendsPath := filepath.Join(baseDir, "friends.json")

	os.MkdirAll(baseDir, 0700)

	config := loadOrCreateConfig(configPath)
	startServer(config)
	showBanner(config)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("total> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "help":
			showHelp()
		case "friends":
			friends := loadFriends(friendsPath)
			showFriends(friends)
		case "friendadd":
			if len(args) < 1 {
				fmt.Println("Usage: friendadd username@ip:port")
				continue
			}
			addFriend(friendsPath, args[0])
		case "exit", "quit":
			fmt.Println("bye 👋")
			return
		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}
	}
}

func getBaseDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".total")
}

func loadOrCreateConfig(path string) Config {
	if _, err := os.Stat(path); err == nil {
		data, _ := os.ReadFile(path)
		var cfg Config
		json.Unmarshal(data, &cfg)
		return cfg
	}

	fmt.Println("Welcome to Total.")
	fmt.Print("Create your account.\nEnter your username: ")
	var username string
	fmt.Scanln(&username)

	ip := getLocalIP()
	cfg := Config{
		Username: username,
		IP:       ip,
		Port:     4040,
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(path, data, 0600)

	fmt.Println("Account created.")
	return cfg
}

func getLocalIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			ip := ipnet.IP.To4()
			if ip != nil && !ip.IsLoopback() {
				return ip.String()
			}
		}
	}
	return "127.0.0.1"
}

func loadFriends(path string) Friends {
	if _, err := os.Stat(path); err != nil {
		return Friends{}
	}
	data, _ := os.ReadFile(path)
	var f Friends
	json.Unmarshal(data, &f)
	return f
}

func addFriend(path string, input string) {
	parts := strings.Split(input, "@")
	if len(parts) != 2 {
		fmt.Println("Invalid format. Use username@ip:port")
		return
	}

	username := parts[0]
	address := parts[1]

	friends := loadFriends(path)
	friends[username] = address

	data, _ := json.MarshalIndent(friends, "", "  ")
	os.WriteFile(path, data, 0600)

	fmt.Println("Friend added:", username)
}

func showBanner(cfg Config) {
	fmt.Println("TOTAL")
	fmt.Printf("│ You: %s@%s:%d\n", cfg.Username, cfg.IP, cfg.Port)
}

func showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                - Show this help message")
	fmt.Println("  friends             - List all friends")
	fmt.Println("  friendadd <user@ip> - Add a new friend")
	fmt.Println("  exit / quit         - Exit the program")
}

func showFriends(friends Friends) {
	if len(friends) == 0 {
		fmt.Println("No friends yet.")
		return
	}
	fmt.Println("Friends:")
	i := 1
	for name, addr := range friends {
		fmt.Printf("  [%d] %s (%s)\n", i, name, addr)
		i++
	}
}

func startServer(cfg Config) {
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Printf("Error starting TCP server: %v\n", err)
			return
		}
		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err == nil {
				// In the future, handle connections here.
				conn.Close()
			}
		}
	}()
}
