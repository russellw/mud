package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn, game *Game) {
	defer conn.Close()
	
	fmt.Fprintf(conn, "Welcome to the MUD!\r\n")
	fmt.Fprintf(conn, "What is your name? ")
	
	scanner := bufio.NewScanner(conn)
	
	if !scanner.Scan() {
		return
	}
	
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Fprintf(conn, "Invalid name. Goodbye!\r\n")
		return
	}
	
	player := &Player{
		conn:      conn,
		name:      name,
		scanner:   scanner,
		inventory: make([]*Item, 0),
	}
	
	game.AddPlayer(player)
	defer game.RemovePlayer(player)
	
	player.SendMessage(fmt.Sprintf("Hello, %s!", name))
	player.HandleCommand(game, "look")
	
	player.location.Broadcast(fmt.Sprintf("%s has entered the game.", name), player)
	
	for scanner.Scan() {
		command := scanner.Text()
		player.HandleCommand(game, command)
	}
	
	player.location.Broadcast(fmt.Sprintf("%s has left the game.", name), player)
}

func main() {
	game := NewGame()
	
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer listener.Close()
	
	fmt.Println("MUD server listening on port 4000")
	fmt.Println("Connect with: telnet localhost 4000")
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		
		go handleConnection(conn, game)
	}
}