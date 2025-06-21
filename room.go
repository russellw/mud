package main

type Item struct {
	name        string
	description string
}

type Room struct {
	name        string
	description string
	players     []*Player
	items       []*Item
	exits       map[string]*Room
}

func (r *Room) Broadcast(message string, except *Player) {
	for _, player := range r.players {
		if player != except {
			player.SendMessage(message)
		}
	}
}