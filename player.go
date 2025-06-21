package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Player struct {
	conn     net.Conn
	name     string
	location *Room
	scanner  *bufio.Scanner
}

func (p *Player) SendMessage(message string) {
	if p.conn != nil {
		fmt.Fprintf(p.conn, "%s\r\n", message)
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
		if p.conn != nil {
			p.conn.Close()
		}
		
	default:
		p.SendMessage("Unknown command. Try: look, say, quit")
	}
}