package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Player struct {
	conn     net.Conn
	name     string
	location *Room
	scanner  *bufio.Scanner
}

type Room struct {
	name        string
	description string
	players     []*Player
}

type Game struct {
	rooms   map[string]*Room
	players []*Player
}

func NewGame() *Game {
	game := &Game{
		rooms:   make(map[string]*Room),
		players: make([]*Player, 0),
	}
	
	startRoom := &Room{
		name:        "Town Square",
		description: "You are standing in a bustling town square. There are paths leading in all directions.",
		players:     make([]*Player, 0),
	}
	
	game.rooms["town_square"] = startRoom
	
	return game
}

func (g *Game) AddPlayer(player *Player) {
	g.players = append(g.players, player)
	startRoom := g.rooms["town_square"]
	player.location = startRoom
	startRoom.players = append(startRoom.players, player)
}

func (g *Game) RemovePlayer(player *Player) {
	for i, p := range g.players {
		if p == player {
			g.players = append(g.players[:i], g.players[i+1:]...)
			break
		}
	}
	
	if player.location != nil {
		for i, p := range player.location.players {
			if p == player {
				player.location.players = append(player.location.players[:i], player.location.players[i+1:]...)
				break
			}
		}
	}
}

func (p *Player) SendMessage(message string) {
	fmt.Fprintf(p.conn, "%s\r\n", message)
}

func (r *Room) Broadcast(message string, except *Player) {
	for _, player := range r.players {
		if player != except {
			player.SendMessage(message)
		}
	}
}

func (p *Player) HandleCommand(game *Game, command string) {
	parts := strings.Fields(strings.TrimSpace(command))
	if len(parts) == 0 {
		return
	}
	
	cmd := strings.ToLower(parts[0])
	
	switch cmd {
	case "look", "l":
		p.SendMessage(fmt.Sprintf("=== %s ===", p.location.name))
		p.SendMessage(p.location.description)
		
		if len(p.location.players) > 1 {
			p.SendMessage("\nOther players here:")
			for _, player := range p.location.players {
				if player != p {
					p.SendMessage(fmt.Sprintf("  %s", player.name))
				}
			}
		}
		
	case "say":
		if len(parts) < 2 {
			p.SendMessage("Say what?")
			return
		}
		message := strings.Join(parts[1:], " ")
		p.SendMessage(fmt.Sprintf("You say: %s", message))
		p.location.Broadcast(fmt.Sprintf("%s says: %s", p.name, message), p)
		
	case "quit", "q":
		p.SendMessage("Goodbye!")
		p.conn.Close()
		
	default:
		p.SendMessage("Unknown command. Try: look, say, quit")
	}
}

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
		conn:    conn,
		name:    name,
		scanner: scanner,
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