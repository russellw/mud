package main

type Game struct {
	rooms   map[string]*Room
	players []*Player
}

func NewGame() *Game {
	game := &Game{
		rooms:   make(map[string]*Room),
		players: make([]*Player, 0),
	}
	
	game.createWorld()
	
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
		items:       []*Item{{name: "wooden mug", description: "A sturdy wooden drinking mug"}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	forest := &Room{
		name:        "Dark Forest",
		description: "A dense forest with towering trees that block most of the sunlight. Strange sounds echo from the shadows.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "twisted branch", description: "A gnarled branch that could serve as a walking stick"}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	market := &Room{
		name:        "Marketplace",
		description: "A busy marketplace with merchants hawking their wares. Colorful stalls line the cobblestone square.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "shiny coin", description: "A gold coin that glints in the sunlight"}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	temple := &Room{
		name:        "Ancient Temple",
		description: "A sacred temple with marble columns and intricate carvings. A sense of peace fills the air.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "prayer book", description: "An old leather-bound book of prayers and rituals"}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	townSquare.exits["north"] = tavern
	townSquare.exits["south"] = forest
	townSquare.exits["east"] = market
	townSquare.exits["west"] = temple
	
	tavern.exits["south"] = townSquare
	forest.exits["north"] = townSquare
	market.exits["west"] = townSquare
	temple.exits["east"] = townSquare
	
	g.rooms["town_square"] = townSquare
	g.rooms["tavern"] = tavern
	g.rooms["forest"] = forest
	g.rooms["market"] = market
	g.rooms["temple"] = temple
	
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
	rat := NewMonster("giant rat", "A large, mangy rat with red eyes and yellowed teeth", 15, 3, true)
	rat.location = g.rooms["forest"]
	g.rooms["forest"].monsters = append(g.rooms["forest"].monsters, rat)
	
	wolf := NewMonster("dire wolf", "A massive wolf with silver fur and piercing blue eyes", 25, 6, true)
	wolf.location = g.rooms["forest"]
	g.rooms["forest"].monsters = append(g.rooms["forest"].monsters, wolf)
	
	bandit := NewMonster("bandit", "A shifty-looking human in leather armor, clutching a rusty dagger", 20, 5, false)
	bandit.location = g.rooms["market"]
	g.rooms["market"].monsters = append(g.rooms["market"].monsters, bandit)
	
	cultist := NewMonster("shadowy cultist", "A robed figure with glowing red eyes chanting in an unknown language", 18, 4, false)
	cultist.location = g.rooms["temple"]
	g.rooms["temple"].monsters = append(g.rooms["temple"].monsters, cultist)
}