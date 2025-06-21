package main

type Room struct {
	name        string
	description string
	players     []*Player
}

func (r *Room) Broadcast(message string, except *Player) {
	for _, player := range r.players {
		if player != except {
			player.SendMessage(message)
		}
	}
}