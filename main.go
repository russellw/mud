package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func handleConnection(conn net.Conn, game *Game) {
	defer conn.Close()
	
	GlobalTelemetry.IncrementConnections()
	defer GlobalTelemetry.DecrementActiveConnections()
	
	fmt.Fprintf(conn, "%sWelcome to the MUD!%s\r\n", ColorBold+ColorBrightGreen, ColorReset)
	fmt.Fprintf(conn, "%sWhat is your name?%s ", ColorBrightCyan, ColorReset)
	
	scanner := bufio.NewScanner(conn)
	
	if !scanner.Scan() {
		return
	}
	
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Fprintf(conn, "%sInvalid name. Goodbye!%s\r\n", ColorError(""), ColorReset)
		return
	}
	
	player := &Player{
		conn:      conn,
		name:      name,
		scanner:   scanner,
		inventory: make([]*Item, 0),
		health:    30,
		maxHealth: 30,
		damage:    5,
	}
	
	GlobalTelemetry.IncrementPlayersCreated()
	game.AddPlayer(player)
	defer game.RemovePlayer(player)
	
	player.SendMessage(fmt.Sprintf("%sHello, %s!%s", ColorBrightGreen, ColorName(name), ColorReset))
	player.HandleCommand(game, "look")
	
	player.location.Broadcast(fmt.Sprintf("%s has entered the game.", ColorName(name)), player)
	
	for scanner.Scan() {
		command := scanner.Text()
		player.HandleCommand(game, command)
	}
	
	player.location.Broadcast(fmt.Sprintf("%s has left the game.", ColorName(name)), player)
}

func main() {
	game := NewGame()
	
	GlobalTelemetry.StartPeriodicLogging(5 * time.Minute)
	
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer listener.Close()
	
	fmt.Println("MUD server listening on port 4000")
	fmt.Println("Connect with: telnet localhost 4000")
	fmt.Printf("Telemetry: %s\n", GlobalTelemetry.GetSummary())
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		
		go handleConnection(conn, game)
	}
}