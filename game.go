package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	rooms   map[string]*Room
	players []*Player
	running bool
}

func NewGame() *Game {
	game := &Game{
		rooms:   make(map[string]*Room),
		players: make([]*Player, 0),
		running: true,
	}
	
	game.createWorld()
	go game.gameLoop()
	
	return game
}

func (g *Game) createWorld() {
	townSquare := &Room{
		name:        "Town Square",
		description: "You are standing in a bustling town square. There are paths leading in all directions.",
		players:     make([]*Player, 0),
		items:       make([]*Item, 0),
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	tavern := &Room{
		name:        "The Prancing Pony Tavern",
		description: "A cozy tavern filled with the smell of ale and roasted meat. Wooden tables and chairs are scattered around.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "wooden mug", description: "A sturdy wooden drinking mug", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	forest := &Room{
		name:        "Dark Forest",
		description: "A dense forest with towering trees that block most of the sunlight. Strange sounds echo from the shadows.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "twisted branch", description: "A gnarled branch that could serve as a walking stick", itemType: "weapon", damage: 2, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	market := &Room{
		name:        "Marketplace",
		description: "A busy marketplace with merchants hawking their wares. Colorful stalls line the cobblestone square.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "shiny coin", description: "A gold coin that glints in the sunlight", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	temple := &Room{
		name:        "Ancient Temple",
		description: "A sacred temple with marble columns and intricate carvings. A sense of peace fills the air.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "prayer book", description: "An old leather-bound book of prayers and rituals", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	dungeon := &Room{
		name:        "Dungeon Entrance",
		description: "A crumbling stone entrance leads into darkness. Ancient torches flicker on the walls.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "rusty key", description: "An old iron key, corroded with age", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	deepForest := &Room{
		name:        "Deep Forest",
		description: "The forest grows darker here. Thick canopy blocks all sunlight. Something large moves in the shadows.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "iron sword", description: "A well-forged iron sword with a sharp edge", itemType: "weapon", damage: 8, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	catacombs := &Room{
		name:        "Ancient Catacombs",
		description: "Narrow stone passages wind through countless burial chambers. The air is thick with age and mystery.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "leather armor", description: "Sturdy leather armor that provides good protection", itemType: "armor", damage: 0, defense: 3}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	wizardTower := &Room{
		name:        "Wizard's Tower",
		description: "A tall stone tower filled with magical artifacts and glowing crystals. Books float in mid-air.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "magic staff", description: "A wooden staff topped with a glowing crystal orb", itemType: "weapon", damage: 10, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	dragonLair := &Room{
		name:        "Dragon's Lair",
		description: "A massive cavern with piles of gold and treasure. Scorch marks cover the walls. The air shimmers with heat.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "dragon scale", description: "A massive golden scale, still warm to the touch", itemType: "armor", damage: 0, defense: 8}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	cemetery := &Room{
		name:        "Moonlit Cemetery",
		description: "Ancient gravestones stretch as far as you can see. Mist swirls between the weathered monuments.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "silver cross", description: "A blessed silver cross that gleams in the moonlight", itemType: "weapon", damage: 6, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	armory := &Room{
		name:        "Castle Armory",
		description: "Weapons and armor line the walls of this military storehouse. Everything is kept in perfect condition.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "steel shield", description: "A heavy steel shield emblazoned with a royal crest", itemType: "armor", damage: 0, defense: 5}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Connect main areas
	townSquare.exits["north"] = tavern
	townSquare.exits["south"] = forest
	townSquare.exits["east"] = market
	townSquare.exits["west"] = temple
	
	tavern.exits["south"] = townSquare
	tavern.exits["up"] = armory
	
	forest.exits["north"] = townSquare
	forest.exits["south"] = deepForest
	forest.exits["down"] = dungeon
	
	market.exits["west"] = townSquare
	
	temple.exits["east"] = townSquare
	temple.exits["down"] = catacombs
	temple.exits["up"] = wizardTower
	
	// Connect extended areas
	dungeon.exits["up"] = forest
	dungeon.exits["north"] = catacombs
	
	deepForest.exits["north"] = forest
	deepForest.exits["west"] = cemetery
	deepForest.exits["east"] = dragonLair
	
	catacombs.exits["up"] = temple
	catacombs.exits["south"] = dungeon
	
	wizardTower.exits["down"] = temple
	
	dragonLair.exits["west"] = deepForest
	
	cemetery.exits["east"] = deepForest
	
	armory.exits["down"] = tavern
	
	g.rooms["town_square"] = townSquare
	g.rooms["tavern"] = tavern
	g.rooms["forest"] = forest
	g.rooms["market"] = market
	g.rooms["temple"] = temple
	g.rooms["dungeon"] = dungeon
	g.rooms["deep_forest"] = deepForest
	g.rooms["catacombs"] = catacombs
	g.rooms["wizard_tower"] = wizardTower
	g.rooms["dragon_lair"] = dragonLair
	g.rooms["cemetery"] = cemetery
	g.rooms["armory"] = armory
	
	g.spawnMonsters()
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

func (g *Game) spawnMonsters() {
	// Original forest monsters
	rat := NewMonster("giant rat", "A large, mangy rat with red eyes and yellowed teeth", 15, 3, true)
	rat.location = g.rooms["forest"]
	g.rooms["forest"].monsters = append(g.rooms["forest"].monsters, rat)
	
	wolf := NewMonster("dire wolf", "A massive wolf with silver fur and piercing blue eyes", 25, 6, true)
	wolf.location = g.rooms["forest"]
	g.rooms["forest"].monsters = append(g.rooms["forest"].monsters, wolf)
	
	// Market monsters
	bandit := NewMonster("bandit", "A shifty-looking human in leather armor, clutching a rusty dagger", 20, 5, true)
	bandit.location = g.rooms["market"]
	g.rooms["market"].monsters = append(g.rooms["market"].monsters, bandit)
	
	// Temple monsters
	cultist := NewMonster("shadowy cultist", "A robed figure with glowing red eyes chanting in an unknown language", 18, 4, true)
	cultist.location = g.rooms["temple"]
	g.rooms["temple"].monsters = append(g.rooms["temple"].monsters, cultist)
	
	// Deep forest monsters
	bear := NewMonster("cave bear", "A massive brown bear with razor-sharp claws and a thunderous roar", 40, 8, true)
	bear.location = g.rooms["deep_forest"]
	g.rooms["deep_forest"].monsters = append(g.rooms["deep_forest"].monsters, bear)
	
	// Dungeon monsters
	skeleton := NewMonster("skeleton warrior", "An ancient skeleton in rusted armor, wielding a bone sword", 20, 5, false)
	skeleton.location = g.rooms["dungeon"]
	g.rooms["dungeon"].monsters = append(g.rooms["dungeon"].monsters, skeleton)
	
	// Catacombs monsters
	zombie := NewMonster("shambling zombie", "A rotting corpse that moves with unnatural hunger", 25, 4, true)
	zombie.location = g.rooms["catacombs"]
	g.rooms["catacombs"].monsters = append(g.rooms["catacombs"].monsters, zombie)
	
	mummy := NewMonster("ancient mummy", "Wrapped in decaying bandages, this ancient guardian protects the tombs", 30, 6, false)
	mummy.location = g.rooms["catacombs"]
	g.rooms["catacombs"].monsters = append(g.rooms["catacombs"].monsters, mummy)
	
	// Wizard tower monsters
	imp := NewMonster("fire imp", "A small demonic creature wreathed in flames with a mischievous grin", 15, 7, true)
	imp.location = g.rooms["wizard_tower"]
	g.rooms["wizard_tower"].monsters = append(g.rooms["wizard_tower"].monsters, imp)
	
	// Dragon lair monsters
	dragon := NewMonster("ancient dragon", "A colossal red dragon with scales like molten gold and eyes like burning coals", 100, 15, true)
	dragon.location = g.rooms["dragon_lair"]
	g.rooms["dragon_lair"].monsters = append(g.rooms["dragon_lair"].monsters, dragon)
	
	// Cemetery monsters
	ghost := NewMonster("wandering spirit", "A translucent figure that wails mournfully as it drifts between the graves", 20, 3, false)
	ghost.location = g.rooms["cemetery"]
	g.rooms["cemetery"].monsters = append(g.rooms["cemetery"].monsters, ghost)
	
	wraith := NewMonster("vengeful wraith", "A dark specter filled with malice and hatred for the living", 35, 9, true)
	wraith.location = g.rooms["cemetery"]
	g.rooms["cemetery"].monsters = append(g.rooms["cemetery"].monsters, wraith)
	
	// Armory monsters
	guard := NewMonster("castle guard", "A heavily armored soldier sworn to protect the castle's treasures", 35, 7, false)
	guard.location = g.rooms["armory"]
	g.rooms["armory"].monsters = append(g.rooms["armory"].monsters, guard)
}

func (g *Game) gameLoop() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	
	for g.running {
		select {
		case <-ticker.C:
			g.processMonsterAI()
		}
	}
}

func (g *Game) processMonsterAI() {
	for _, room := range g.rooms {
		for _, monster := range room.monsters {
			if !monster.alive {
				continue
			}
			
			if monster.aggressive && len(room.players) > 0 {
				target := room.players[rand.Intn(len(room.players))]
				g.MonsterAttackPlayer(monster, target)
			} else if !monster.aggressive && len(room.players) > 0 {
				if rand.Float32() < 0.3 {
					target := room.players[rand.Intn(len(room.players))]
					g.MonsterAttackPlayer(monster, target)
				}
			}
		}
	}
}

func (g *Game) PlayerAttackMonster(player *Player, monsterName string) {
	var target *Monster
	for _, monster := range player.location.monsters {
		if monster.name == monsterName && monster.alive {
			target = monster
			break
		}
	}
	
	if target == nil {
		player.SendMessage("There is no such monster here.")
		return
	}
	
	baseDamage := player.damage
	if player.weapon != nil {
		baseDamage += player.weapon.damage
	}
	damage := baseDamage + rand.Intn(3) - 1
	if damage < 1 {
		damage = 1
	}
	
	isDead := target.TakeDamage(damage)
	
	if isDead {
		player.SendMessage(fmt.Sprintf("You kill the %s!", target.name))
		player.location.Broadcast(fmt.Sprintf("%s kills the %s!", player.name, target.name), player)
	} else {
		player.SendMessage(fmt.Sprintf("You attack the %s for %d damage!", target.name, damage))
		player.location.Broadcast(fmt.Sprintf("%s attacks the %s!", player.name, target.name), player)
	}
}

func (g *Game) MonsterAttackPlayer(monster *Monster, player *Player) {
	if !monster.alive {
		return
	}
	
	baseDamage := monster.damage + rand.Intn(5) - 2
	if baseDamage < 1 {
		baseDamage = 1
	}
	
	defense := 0
	if player.armor != nil {
		defense = player.armor.defense
	}
	
	damage := baseDamage - defense
	if damage < 1 {
		damage = 1
	}
	
	isDead := player.TakeDamage(damage)
	
	if isDead {
		player.SendMessage("You have been killed!")
		player.location.Broadcast(fmt.Sprintf("%s has been killed by %s!", player.name, monster.name), player)
		g.respawnPlayer(player)
	} else {
		player.SendMessage(fmt.Sprintf("The %s attacks you for %d damage!", monster.name, damage))
		player.location.Broadcast(fmt.Sprintf("%s attacks %s!", monster.name, player.name), player)
	}
}

func (g *Game) respawnPlayer(player *Player) {
	player.health = player.maxHealth
	
	if player.location != nil {
		for i, p := range player.location.players {
			if p == player {
				player.location.players = append(player.location.players[:i], player.location.players[i+1:]...)
				break
			}
		}
	}
	
	townSquare := g.rooms["town_square"]
	player.location = townSquare
	townSquare.players = append(townSquare.players, player)
	
	player.SendMessage("You respawn in the town square, fully healed.")
	player.location.Broadcast(fmt.Sprintf("%s respawns.", player.name), player)
}