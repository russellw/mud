package main

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame()
	
	if game == nil {
		t.Fatal("NewGame() returned nil")
	}
	
	if game.rooms == nil {
		t.Error("Game rooms map is nil")
	}
	
	if game.players == nil {
		t.Error("Game players slice is nil")
	}
	
	if len(game.players) != 0 {
		t.Errorf("Expected 0 players initially, got %d", len(game.players))
	}
	
	townSquare, exists := game.rooms["town_square"]
	if !exists {
		t.Error("Town square room not created")
	}
	
	if townSquare.name != "Town Square" {
		t.Errorf("Expected room name 'Town Square', got '%s'", townSquare.name)
	}
	
	if townSquare.description == "" {
		t.Error("Town square has no description")
	}
	
	if len(townSquare.players) != 0 {
		t.Errorf("Expected 0 players in town square initially, got %d", len(townSquare.players))
	}
}

func TestAddPlayer(t *testing.T) {
	game := NewGame()
	
	player := &Player{
		name: "TestPlayer",
		conn: nil,
	}
	
	game.AddPlayer(player)
	
	if len(game.players) != 1 {
		t.Errorf("Expected 1 player in game, got %d", len(game.players))
	}
	
	if game.players[0] != player {
		t.Error("Player not properly added to game")
	}
	
	if player.location == nil {
		t.Error("Player location not set")
	}
	
	if player.location.name != "Town Square" {
		t.Errorf("Expected player to be in Town Square, got '%s'", player.location.name)
	}
	
	townSquare := game.rooms["town_square"]
	if len(townSquare.players) != 1 {
		t.Errorf("Expected 1 player in town square, got %d", len(townSquare.players))
	}
	
	if townSquare.players[0] != player {
		t.Error("Player not properly added to room")
	}
}

func TestRemovePlayer(t *testing.T) {
	game := NewGame()
	
	player1 := &Player{name: "Player1", conn: nil}
	player2 := &Player{name: "Player2", conn: nil}
	
	game.AddPlayer(player1)
	game.AddPlayer(player2)
	
	if len(game.players) != 2 {
		t.Errorf("Expected 2 players before removal, got %d", len(game.players))
	}
	
	game.RemovePlayer(player1)
	
	if len(game.players) != 1 {
		t.Errorf("Expected 1 player after removal, got %d", len(game.players))
	}
	
	if game.players[0] != player2 {
		t.Error("Wrong player removed from game")
	}
	
	townSquare := game.rooms["town_square"]
	if len(townSquare.players) != 1 {
		t.Errorf("Expected 1 player in town square after removal, got %d", len(townSquare.players))
	}
	
	if townSquare.players[0] != player2 {
		t.Error("Wrong player removed from room")
	}
}

func TestRemovePlayerWithoutLocation(t *testing.T) {
	game := NewGame()
	
	player := &Player{
		name:     "TestPlayer",
		conn:     nil,
		location: nil,
	}
	
	game.players = append(game.players, player)
	
	game.RemovePlayer(player)
	
	if len(game.players) != 0 {
		t.Errorf("Expected 0 players after removal, got %d", len(game.players))
	}
}

func TestMultiplePlayersInRoom(t *testing.T) {
	game := NewGame()
	
	player1 := &Player{name: "Player1", conn: nil}
	player2 := &Player{name: "Player2", conn: nil}
	player3 := &Player{name: "Player3", conn: nil}
	
	game.AddPlayer(player1)
	game.AddPlayer(player2)
	game.AddPlayer(player3)
	
	townSquare := game.rooms["town_square"]
	
	if len(townSquare.players) != 3 {
		t.Errorf("Expected 3 players in town square, got %d", len(townSquare.players))
	}
	
	game.RemovePlayer(player2)
	
	if len(townSquare.players) != 2 {
		t.Errorf("Expected 2 players in town square after removal, got %d", len(townSquare.players))
	}
	
	foundPlayer1 := false
	foundPlayer3 := false
	for _, p := range townSquare.players {
		if p == player1 {
			foundPlayer1 = true
		}
		if p == player3 {
			foundPlayer3 = true
		}
		if p == player2 {
			t.Error("Removed player still found in room")
		}
	}
	
	if !foundPlayer1 || !foundPlayer3 {
		t.Error("Expected players not found in room after removal")
	}
}

func TestWorldCreation(t *testing.T) {
	game := NewGame()
	
	expectedRooms := []string{"town_square", "tavern", "forest", "market", "temple", "dungeon", "deep_forest", "catacombs", "wizard_tower", "dragon_lair", "cemetery", "armory", "pirate_cove", "volcanic_cavern", "ice_fortress", "sky_temple", "cursed_swamp", "crystal_mines", "haunted_library", "goblin_warren"}
	
	if len(game.rooms) != len(expectedRooms) {
		t.Errorf("Expected %d rooms, got %d", len(expectedRooms), len(game.rooms))
	}
	
	for _, roomKey := range expectedRooms {
		room, exists := game.rooms[roomKey]
		if !exists {
			t.Errorf("Room '%s' not found", roomKey)
			continue
		}
		
		if room.name == "" {
			t.Errorf("Room '%s' has empty name", roomKey)
		}
		
		if room.description == "" {
			t.Errorf("Room '%s' has empty description", roomKey)
		}
		
		if room.players == nil {
			t.Errorf("Room '%s' has nil players slice", roomKey)
		}
		
		if room.items == nil {
			t.Errorf("Room '%s' has nil items slice", roomKey)
		}
		
		if room.exits == nil {
			t.Errorf("Room '%s' has nil exits map", roomKey)
		}
	}
}

func TestRoomConnections(t *testing.T) {
	game := NewGame()
	
	townSquare := game.rooms["town_square"]
	tavern := game.rooms["tavern"]
	forest := game.rooms["forest"]
	market := game.rooms["market"]
	temple := game.rooms["temple"]
	
	if townSquare.exits["north"] != tavern {
		t.Error("Town square north exit should connect to tavern")
	}
	
	if townSquare.exits["south"] != forest {
		t.Error("Town square south exit should connect to forest")
	}
	
	if townSquare.exits["east"] != market {
		t.Error("Town square east exit should connect to market")
	}
	
	if townSquare.exits["west"] != temple {
		t.Error("Town square west exit should connect to temple")
	}
	
	if tavern.exits["south"] != townSquare {
		t.Error("Tavern south exit should connect to town square")
	}
	
	if forest.exits["north"] != townSquare {
		t.Error("Forest north exit should connect to town square")
	}
	
	if market.exits["west"] != townSquare {
		t.Error("Market west exit should connect to town square")
	}
	
	if temple.exits["east"] != townSquare {
		t.Error("Temple east exit should connect to town square")
	}
}

func TestRoomItems(t *testing.T) {
	game := NewGame()
	
	townSquare := game.rooms["town_square"]
	if len(townSquare.items) != 0 {
		t.Errorf("Town square should have 0 items, got %d", len(townSquare.items))
	}
	
	tavern := game.rooms["tavern"]
	if len(tavern.items) != 1 {
		t.Errorf("Tavern should have 1 item, got %d", len(tavern.items))
	} else if tavern.items[0].name != "wooden mug" {
		t.Errorf("Expected wooden mug in tavern, got %s", tavern.items[0].name)
	}
	
	forest := game.rooms["forest"]
	if len(forest.items) != 1 {
		t.Errorf("Forest should have 1 item, got %d", len(forest.items))
	} else if forest.items[0].name != "twisted branch" {
		t.Errorf("Expected twisted branch in forest, got %s", forest.items[0].name)
	}
	
	market := game.rooms["market"]
	if len(market.items) != 1 {
		t.Errorf("Market should have 1 item, got %d", len(market.items))
	} else if market.items[0].name != "shiny coin" {
		t.Errorf("Expected shiny coin in market, got %s", market.items[0].name)
	}
	
	temple := game.rooms["temple"]
	if len(temple.items) != 1 {
		t.Errorf("Temple should have 1 item, got %d", len(temple.items))
	} else if temple.items[0].name != "prayer book" {
		t.Errorf("Expected prayer book in temple, got %s", temple.items[0].name)
	}
}