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
	health    int
	maxHealth int
	damage    int
	weapon    *Item
	armor     *Item
}

func (p *Player) SendMessage(message string) {
	if p.conn != nil {
		fmt.Fprintf(p.conn, "%s\r\n", message)
	}
}

func (p *Player) TakeDamage(damage int) bool {
	p.health -= damage
	if p.health <= 0 {
		p.health = 0
		return true
	}
	return false
}

func (p *Player) GetHealthStatus() string {
	healthPercent := float64(p.health) / float64(p.maxHealth)
	
	if healthPercent > 0.8 {
		return "You feel great!"
	} else if healthPercent > 0.6 {
		return "You have some minor cuts and bruises."
	} else if healthPercent > 0.4 {
		return "You are wounded."
	} else if healthPercent > 0.2 {
		return "You are badly wounded."
	} else {
		return "You are near death!"
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
		p.SendMessage(fmt.Sprintf("=== %s ===", ColorRoomName(p.location.name)))
		p.SendMessage(ColorDescription(p.location.description))
		
		if len(p.location.items) > 0 {
			p.SendMessage(fmt.Sprintf("\n%sItems here:%s", ColorBold, ColorReset))
			for _, item := range p.location.items {
				p.SendMessage(fmt.Sprintf("  %s", ColorItem(item.name)))
			}
		}
		
		if len(p.location.exits) > 0 {
			p.SendMessage(fmt.Sprintf("\n%sExits:%s", ColorBold, ColorReset))
			for direction := range p.location.exits {
				p.SendMessage(fmt.Sprintf("  %s", ColorExit(direction)))
			}
		}
		
		if len(p.location.monsters) > 0 {
			p.SendMessage(fmt.Sprintf("\n%sMonsters here:%s", ColorBold+ColorRed, ColorReset))
			for _, monster := range p.location.monsters {
				if monster.alive {
					p.SendMessage(fmt.Sprintf("  %s", ColorMonster(monster.GetStatus())))
				}
			}
		}
		
		if len(p.location.players) > 1 {
			p.SendMessage(fmt.Sprintf("\n%sOther players here:%s", ColorBold+ColorBlue, ColorReset))
			for _, player := range p.location.players {
				if player != p {
					p.SendMessage(fmt.Sprintf("  %s", ColorPlayer(player.name)))
				}
			}
		}
		
	case "go", "north", "n", "south", "s", "east", "e", "west", "w", "up", "u", "down", "d":
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
		} else if direction == "u" {
			direction = "up"
		} else if direction == "d" {
			direction = "down"
		}
		
		nextRoom, exists := p.location.exits[direction]
		if !exists {
			p.SendMessage(ColorError("You can't go that way."))
			return
		}
		
		p.location.Broadcast(fmt.Sprintf("%s leaves %s.", ColorName(p.name), ColorExit(direction)), p)
		
		for i, player := range p.location.players {
			if player == p {
				p.location.players = append(p.location.players[:i], p.location.players[i+1:]...)
				break
			}
		}
		
		p.location = nextRoom
		nextRoom.players = append(nextRoom.players, p)
		
		nextRoom.Broadcast(fmt.Sprintf("%s arrives.", ColorName(p.name)), p)
		p.HandleCommand(game, "look")
		
	case "get", "take":
		if len(parts) < 2 {
			p.SendMessage(ColorWarning("Get what?"))
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		for i, item := range p.location.items {
			if strings.ToLower(item.name) == itemName {
				p.location.items = append(p.location.items[:i], p.location.items[i+1:]...)
				p.inventory = append(p.inventory, item)
				p.SendMessage(fmt.Sprintf("You take the %s.", ColorItem(item.name)))
				p.location.Broadcast(fmt.Sprintf("%s takes the %s.", ColorName(p.name), ColorItem(item.name)), p)
				return
			}
		}
		p.SendMessage(ColorError("That item is not here."))
		
	case "drop":
		if len(parts) < 2 {
			p.SendMessage(ColorWarning("Drop what?"))
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		for i, item := range p.inventory {
			if strings.ToLower(item.name) == itemName {
				p.inventory = append(p.inventory[:i], p.inventory[i+1:]...)
				p.location.items = append(p.location.items, item)
				p.SendMessage(fmt.Sprintf("You drop the %s.", ColorItem(item.name)))
				p.location.Broadcast(fmt.Sprintf("%s drops the %s.", ColorName(p.name), ColorItem(item.name)), p)
				return
			}
		}
		p.SendMessage(ColorError("You don't have that item."))
		
	case "inventory", "inv", "i":
		if len(p.inventory) == 0 {
			p.SendMessage(ColorInfo("You are not carrying anything."))
		} else {
			p.SendMessage(fmt.Sprintf("%sYou are carrying:%s", ColorBold, ColorReset))
			for _, item := range p.inventory {
				p.SendMessage(fmt.Sprintf("  %s", ColorItem(item.name)))
			}
		}
		
	case "examine", "ex":
		if len(parts) < 2 {
			p.SendMessage(ColorWarning("Examine what?"))
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		for _, item := range p.location.items {
			if strings.ToLower(item.name) == itemName {
				description := fmt.Sprintf("%s: %s", ColorItem(item.name), item.description)
				if item.itemType == "weapon" && item.damage > 0 {
					description += fmt.Sprintf(" %s(Damage: +%d)%s", ColorDamage(""), item.damage, ColorReset)
				} else if item.itemType == "armor" && item.defense > 0 {
					description += fmt.Sprintf(" %s(Defense: +%d)%s", ColorEquipment(""), item.defense, ColorReset)
				}
				p.SendMessage(description)
				return
			}
		}
		
		for _, item := range p.inventory {
			if strings.ToLower(item.name) == itemName {
				p.SendMessage(fmt.Sprintf("%s: %s", ColorItem(item.name), item.description))
				return
			}
		}
		p.SendMessage(ColorError("You don't see that here."))
		
	case "who":
		p.SendMessage(fmt.Sprintf("%sPlayers online:%s", ColorBold, ColorReset))
		for _, player := range game.players {
			p.SendMessage(fmt.Sprintf("  %s (%s)", ColorPlayer(player.name), ColorRoomName(player.location.name)))
		}
		
	case "attack", "kill", "fight":
		if len(parts) < 2 {
			p.SendMessage(ColorWarning("Attack what?"))
			return
		}
		targetName := strings.ToLower(strings.Join(parts[1:], " "))
		game.PlayerAttackMonster(p, targetName)
		
	case "health", "hp":
		healthColor := ColorHealing("")
		if p.health < p.maxHealth/2 {
			healthColor = ColorDamage("")
		}
		p.SendMessage(fmt.Sprintf("%sHealth: %d/%d%s", healthColor, p.health, p.maxHealth, ColorReset))
		p.SendMessage(p.GetHealthStatus())
		
	case "equip", "wield", "wear":
		if len(parts) < 2 {
			p.SendMessage(ColorWarning("Equip what?"))
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		var item *Item
		var itemIndex int = -1
		for i, invItem := range p.inventory {
			if strings.ToLower(invItem.name) == itemName {
				item = invItem
				itemIndex = i
				break
			}
		}
		
		if item == nil {
			p.SendMessage(ColorError("You don't have that item."))
			return
		}
		
		if item.itemType == "weapon" {
			if p.weapon != nil {
				p.SendMessage(fmt.Sprintf("You unequip %s and equip %s.", ColorEquipment(p.weapon.name), ColorEquipment(item.name)))
				p.inventory = append(p.inventory, p.weapon)
			} else {
				p.SendMessage(fmt.Sprintf("You equip %s.", ColorEquipment(item.name)))
			}
			p.weapon = item
		} else if item.itemType == "armor" {
			if p.armor != nil {
				p.SendMessage(fmt.Sprintf("You remove %s and wear %s.", ColorEquipment(p.armor.name), ColorEquipment(item.name)))
				p.inventory = append(p.inventory, p.armor)
			} else {
				p.SendMessage(fmt.Sprintf("You wear %s.", ColorEquipment(item.name)))
			}
			p.armor = item
		} else {
			p.SendMessage(ColorError("You can't equip that item."))
			return
		}
		
		p.inventory = append(p.inventory[:itemIndex], p.inventory[itemIndex+1:]...)
		
	case "unequip", "remove":
		if len(parts) < 2 {
			p.SendMessage(ColorWarning("Remove what?"))
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		if p.weapon != nil && strings.ToLower(p.weapon.name) == itemName {
			p.SendMessage(fmt.Sprintf("You unequip %s.", ColorEquipment(p.weapon.name)))
			p.inventory = append(p.inventory, p.weapon)
			p.weapon = nil
		} else if p.armor != nil && strings.ToLower(p.armor.name) == itemName {
			p.SendMessage(fmt.Sprintf("You remove %s.", ColorEquipment(p.armor.name)))
			p.inventory = append(p.inventory, p.armor)
			p.armor = nil
		} else {
			p.SendMessage(ColorError("You don't have that equipped."))
		}
		
	case "equipment", "eq":
		p.SendMessage(fmt.Sprintf("%sEquipment:%s", ColorBold, ColorReset))
		if p.weapon != nil {
			p.SendMessage(fmt.Sprintf("  Weapon: %s %s(+%d damage)%s", ColorEquipment(p.weapon.name), ColorDamage(""), p.weapon.damage, ColorReset))
		} else {
			p.SendMessage("  Weapon: none")
		}
		if p.armor != nil {
			p.SendMessage(fmt.Sprintf("  Armor: %s %s(+%d defense)%s", ColorEquipment(p.armor.name), ColorEquipment(""), p.armor.defense, ColorReset))
		} else {
			p.SendMessage("  Armor: none")
		}
		
	case "use":
		if len(parts) < 2 {
			p.SendMessage(ColorWarning("Use what?"))
			return
		}
		itemName := strings.ToLower(strings.Join(parts[1:], " "))
		
		for i, item := range p.inventory {
			if strings.ToLower(item.name) == itemName {
				switch item.name {
				case "tome of knowledge":
					p.inventory = append(p.inventory[:i], p.inventory[i+1:]...)
					p.maxHealth += 10
					p.health += 10
					p.damage += 2
					p.SendMessage(ColorMagic("You study the ancient tome and feel your mind expand with knowledge!"))
					p.SendMessage(fmt.Sprintf("%sYour maximum health increases to %d and your base damage increases by 2!%s", ColorSuccess(""), p.maxHealth, ColorReset))
					p.location.Broadcast(fmt.Sprintf("%s glows with newfound wisdom!", ColorName(p.name)), p)
					return
				case "prayer book":
					p.health = p.maxHealth
					p.SendMessage(ColorMagic("You recite sacred prayers and feel divine healing wash over you!"))
					p.SendMessage(ColorHealing("You are fully healed!"))
					p.location.Broadcast(fmt.Sprintf("%s radiates with holy light!", ColorName(p.name)), p)
					return
				case "shiny coin":
					p.SendMessage("You flip the coin and make a wish, but nothing happens. It's just a coin.")
					return
				case "rusty key":
					p.SendMessage("The key doesn't seem to fit any locks here.")
					return
				default:
					p.SendMessage(ColorError("You can't use that item."))
					return
				}
			}
		}
		p.SendMessage(ColorError("You don't have that item."))
		
	case "rest":
		healAmount := p.maxHealth / 4
		if healAmount < 5 {
			healAmount = 5
		}
		if p.health >= p.maxHealth {
			p.SendMessage(ColorInfo("You are already at full health."))
			return
		}
		p.health += healAmount
		if p.health > p.maxHealth {
			p.health = p.maxHealth
		}
		p.SendMessage(fmt.Sprintf("%sYou rest and recover %d health points.%s", ColorHealing(""), healAmount, ColorReset))
		p.location.Broadcast(fmt.Sprintf("%s sits down to rest.", ColorName(p.name)), p)
		
	case "say":
		if len(parts) < 2 {
			p.SendMessage(ColorWarning("Say what?"))
			return
		}
		message := strings.Join(parts[1:], " ")
		p.SendMessage(fmt.Sprintf("You say: %s%s%s", ColorBrightWhite, message, ColorReset))
		p.location.Broadcast(fmt.Sprintf("%s says: %s%s%s", ColorName(p.name), ColorBrightWhite, message, ColorReset), p)
		
	case "quit", "q":
		p.SendMessage(ColorInfo("Goodbye!"))
		if p.conn != nil {
			p.conn.Close()
		}
		
	default:
		p.SendMessage(ColorError("Unknown command. Try: look, go <direction>, get <item>, drop <item>, inventory, examine <item>, equip <item>, equipment, attack <monster>, health, who, use <item>, rest, say, quit"))
	}
}