package main

import (
	"strings"
	"testing"
)

func createTestPlayer(name string) *Player {
	return &Player{
		name:      name,
		conn:      nil,
		inventory: make([]*Item, 0),
	}
}

func TestPlayerMovement(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	if player.location.name != "Town Square" {
		t.Errorf("Player should start in Town Square, got %s", player.location.name)
	}
	
	player.HandleCommand(game, "north")
	
	if player.location.name != "The Prancing Pony Tavern" {
		t.Errorf("Player should be in tavern after going north, got %s", player.location.name)
	}
	
	tavernPlayers := game.rooms["tavern"].players
	if len(tavernPlayers) != 1 || tavernPlayers[0] != player {
		t.Error("Player not found in tavern room")
	}
	
	townSquarePlayers := game.rooms["town_square"].players
	if len(townSquarePlayers) != 0 {
		t.Error("Player still found in town square after moving")
	}
}

func TestPlayerMovementDirections(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	testCases := []struct {
		command      string
		expectedRoom string
	}{
		{"n", "The Prancing Pony Tavern"},
		{"s", "Town Square"},
		{"south", "Dark Forest"},
		{"north", "Town Square"},
		{"e", "Marketplace"},
		{"west", "Town Square"},
		{"w", "Ancient Temple"},
		{"east", "Town Square"},
	}
	
	for _, tc := range testCases {
		player.HandleCommand(game, tc.command)
		if player.location.name != tc.expectedRoom {
			t.Errorf("After command '%s', expected room '%s', got '%s'", 
				tc.command, tc.expectedRoom, player.location.name)
		}
	}
}

func TestPlayerMovementInvalidDirection(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	originalRoom := player.location.name
	
	player.HandleCommand(game, "northeast")
	
	if player.location.name != originalRoom {
		t.Error("Player moved to invalid direction")
	}
}

func TestItemPickup(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	player.HandleCommand(game, "north")
	
	tavern := game.rooms["tavern"]
	if len(tavern.items) != 1 {
		t.Errorf("Tavern should have 1 item initially, got %d", len(tavern.items))
	}
	
	player.HandleCommand(game, "get wooden mug")
	
	if len(player.inventory) != 1 {
		t.Errorf("Player should have 1 item in inventory, got %d", len(player.inventory))
	}
	
	if player.inventory[0].name != "wooden mug" {
		t.Errorf("Expected wooden mug in inventory, got %s", player.inventory[0].name)
	}
	
	if len(tavern.items) != 0 {
		t.Errorf("Tavern should have 0 items after pickup, got %d", len(tavern.items))
	}
}

func TestItemDrop(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	mug := &Item{name: "wooden mug", description: "A sturdy wooden drinking mug"}
	player.inventory = append(player.inventory, mug)
	
	townSquare := game.rooms["town_square"]
	initialItems := len(townSquare.items)
	
	player.HandleCommand(game, "drop wooden mug")
	
	if len(player.inventory) != 0 {
		t.Errorf("Player should have 0 items in inventory after drop, got %d", len(player.inventory))
	}
	
	if len(townSquare.items) != initialItems+1 {
		t.Errorf("Room should have %d items after drop, got %d", initialItems+1, len(townSquare.items))
	}
	
	found := false
	for _, item := range townSquare.items {
		if item.name == "wooden mug" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Dropped item not found in room")
	}
}

func TestItemPickupNonexistent(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	originalInventorySize := len(player.inventory)
	
	player.HandleCommand(game, "get nonexistent item")
	
	if len(player.inventory) != originalInventorySize {
		t.Error("Inventory changed when picking up nonexistent item")
	}
}

func TestItemDropNonexistent(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	townSquare := game.rooms["town_square"]
	originalRoomItems := len(townSquare.items)
	
	player.HandleCommand(game, "drop nonexistent item")
	
	if len(townSquare.items) != originalRoomItems {
		t.Error("Room items changed when dropping nonexistent item")
	}
}

func TestInventoryCommand(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	player.HandleCommand(game, "inventory")
	
	mug := &Item{name: "wooden mug", description: "A sturdy wooden drinking mug"}
	coin := &Item{name: "shiny coin", description: "A gold coin"}
	player.inventory = append(player.inventory, mug, coin)
	
	player.HandleCommand(game, "inv")
	player.HandleCommand(game, "i")
}

func TestExamineCommand(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	player.HandleCommand(game, "north")
	
	player.HandleCommand(game, "examine wooden mug")
	
	mug := &Item{name: "wooden mug", description: "A sturdy wooden drinking mug"}
	player.inventory = append(player.inventory, mug)
	
	player.HandleCommand(game, "ex wooden mug")
}

func TestWhoCommand(t *testing.T) {
	game := NewGame()
	player1 := createTestPlayer("Player1")
	player2 := createTestPlayer("Player2")
	
	game.AddPlayer(player1)
	game.AddPlayer(player2)
	
	player1.HandleCommand(game, "who")
}

func TestGoCommand(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	player.HandleCommand(game, "go north")
	
	if player.location.name != "The Prancing Pony Tavern" {
		t.Errorf("Player should be in tavern after 'go north', got %s", player.location.name)
	}
	
	player.HandleCommand(game, "go south")
	
	if player.location.name != "Town Square" {
		t.Errorf("Player should be in town square after 'go south', got %s", player.location.name)
	}
}

func TestGoCommandNoDirection(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	originalRoom := player.location.name
	player.HandleCommand(game, "go")
	
	if player.location.name != originalRoom {
		t.Error("Player moved without specifying direction")
	}
}

func TestSayCommand(t *testing.T) {
	game := NewGame()
	player1 := createTestPlayer("Player1")
	player2 := createTestPlayer("Player2")
	
	game.AddPlayer(player1)
	game.AddPlayer(player2)
	
	player1.HandleCommand(game, "say Hello everyone!")
	
	player1.HandleCommand(game, "say")
}

func TestCaseInsensitiveCommands(t *testing.T) {
	game := NewGame()
	player := createTestPlayer("TestPlayer")
	game.AddPlayer(player)
	
	commands := []string{"LOOK", "Look", "l", "L", "NORTH", "North", "n", "N"}
	
	for _, cmd := range commands {
		originalRoom := player.location.name
		player.HandleCommand(game, cmd)
		
		if strings.Contains(strings.ToLower(cmd), "north") || strings.Contains(strings.ToLower(cmd), "n") {
			if originalRoom == "Town Square" && player.location.name == "Town Square" {
				t.Errorf("Command '%s' should have moved player north", cmd)
			}
			player.HandleCommand(game, "south")
		}
	}
}