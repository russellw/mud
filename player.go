package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Player struct {
	conn      net.Conn
	name      string
	location  *Room
	scanner   *bufio.Scanner
	inventory []*Item
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
		
		if len(p.location.items) > 0 {
			p.SendMessage("\nItems here:")
			for _, item := range p.location.items {
				p.SendMessage(fmt.Sprintf("  %s", item.name))
			}
		}
		
		if len(p.location.exits) > 0 {
			p.SendMessage("\nExits:")
			for direction := range p.location.exits {
				p.SendMessage(fmt.Sprintf("  %s", direction))
			}
		}
		
		if len(p.location.players) > 1 {
			p.SendMessage("\nOther players here:")
			for _, player := range p.location.players {
				if player != p {
					p.SendMessage(fmt.Sprintf("  %s", player.name))
				}
			}
		}
		
	case "go", "north", "n", "south", "s", "east", "e", "west", "w":
		direction := cmd
		if cmd == "go" {
			if len(parts) < 2 {
				p.SendMessage("Go where?")
				return
			}
			direction = strings.ToLower(parts[1])
		}
		
		if direction == "n" {
			direction = "north"
		} else if direction == "s" {
			direction = "south"
		} else if direction == "e" {
			direction = "east"
		} else if direction == "w" {
			direction = "west"
		}
		
		nextRoom, exists := p.location.exits[direction]
		if !exists {
			p.SendMessage("You can't go that way.")
			return
		}
		
		p.location.Broadcast(fmt.Sprintf("%s leaves %s.", p.name, direction), p)
		
		for i, player := range p.location.players {
			if player == p {
				p.location.players = append(p.location.players[:i], p.location.players[i+1:]...)
				break
			}
		}
		
		p.location = nextRoom
		nextRoom.players = append(nextRoom.players, p)
		
		nextRoom.Broadcast(fmt.Sprintf("%s arrives.", p.name), p)
		p.HandleCommand(game, "look")
		
	case "get", "take":
		if len(parts) < 2 {
			p.SendMessage("Get what?")
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		for i, item := range p.location.items {
			if strings.ToLower(item.name) == itemName {
				p.location.items = append(p.location.items[:i], p.location.items[i+1:]...)
				p.inventory = append(p.inventory, item)
				p.SendMessage(fmt.Sprintf("You take the %s.", item.name))
				p.location.Broadcast(fmt.Sprintf("%s takes the %s.", p.name, item.name), p)
				return
			}
		}
		p.SendMessage("That item is not here.")
		
	case "drop":
		if len(parts) < 2 {
			p.SendMessage("Drop what?")
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		for i, item := range p.inventory {
			if strings.ToLower(item.name) == itemName {
				p.inventory = append(p.inventory[:i], p.inventory[i+1:]...)
				p.location.items = append(p.location.items, item)
				p.SendMessage(fmt.Sprintf("You drop the %s.", item.name))
				p.location.Broadcast(fmt.Sprintf("%s drops the %s.", p.name, item.name), p)
				return
			}
		}
		p.SendMessage("You don't have that item.")
		
	case "inventory", "inv", "i":
		if len(p.inventory) == 0 {
			p.SendMessage("You are not carrying anything.")
		} else {
			p.SendMessage("You are carrying:")
			for _, item := range p.inventory {
				p.SendMessage(fmt.Sprintf("  %s", item.name))
			}
		}
		
	case "examine", "ex":
		if len(parts) < 2 {
			p.SendMessage("Examine what?")
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		for _, item := range p.location.items {
			if strings.ToLower(item.name) == itemName {
				p.SendMessage(fmt.Sprintf("%s: %s", item.name, item.description))
				return
			}
		}
		
		for _, item := range p.inventory {
			if strings.ToLower(item.name) == itemName {
				p.SendMessage(fmt.Sprintf("%s: %s", item.name, item.description))
				return
			}
		}
		p.SendMessage("You don't see that here.")
		
	case "who":
		p.SendMessage("Players online:")
		for _, player := range game.players {
			p.SendMessage(fmt.Sprintf("  %s (%s)", player.name, player.location.name))
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
		p.SendMessage("Unknown command. Try: look, go <direction>, get <item>, drop <item>, inventory, examine <item>, who, say, quit")
	}
}