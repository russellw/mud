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