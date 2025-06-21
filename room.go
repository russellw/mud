package main

type Item struct {
	name        string
	description string
	itemType    string // "weapon", "armor", "misc"
	damage      int    // for weapons
	defense     int    // for armor
}

type Room struct {
	name        string
	description string
	players     []*Player
	items       []*Item
	monsters    []*Monster
	exits       map[string]*Room
}

func (r *Room) Broadcast(message string, except *Player) {
	for _, player := range r.players {
		if player != except {
			player.SendMessage(message)
		}
	}
}