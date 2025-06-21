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