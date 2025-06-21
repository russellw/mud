package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

type MockConnection struct {
	messages []string
}

func (m *MockConnection) Write(p []byte) (n int, err error) {
	message := strings.TrimSpace(string(p))
	if message != "" {
		m.messages = append(m.messages, message)
	}
	return len(p), nil
}

func (m *MockConnection) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("mock connection read not implemented")
}

func (m *MockConnection) Close() error {
	return nil
}

func (m *MockConnection) LocalAddr() net.Addr {
	return nil
}

func (m *MockConnection) RemoteAddr() net.Addr {
	return nil
}

func (m *MockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func createMockPlayer(name string) *Player {
	mockConn := &MockConnection{messages: make([]string, 0)}
	return &Player{
		conn:      mockConn,
		name:      name,
		inventory: make([]*Item, 0),
		health:    30,
		maxHealth: 30,
		damage:    5,
	}
}

func getPlayerMessages(player *Player) []string {
	if mockConn, ok := player.conn.(*MockConnection); ok {
		return mockConn.messages
	}
	return []string{}
}

func clearPlayerMessages(player *Player) {
	if mockConn, ok := player.conn.(*MockConnection); ok {
		mockConn.messages = make([]string, 0)
	}
}

func TestGameplayTranscript(t *testing.T) {
	transcript := []string{}
	transcript = append(transcript, "=== MUD Gameplay Transcript ===")
	transcript = append(transcript, fmt.Sprintf("Generated at: %s", time.Now().Format("2006-01-02 15:04:05")))
	transcript = append(transcript, "")
	
	game := NewGame()
	defer func() {
		game.running = false
	}()
	
	alice := createMockPlayer("Alice")
	bob := createMockPlayer("Bob")
	
	transcript = append(transcript, ">>> Alice and Bob enter the game")
	game.AddPlayer(alice)
	game.AddPlayer(bob)
	
	transcript = append(transcript, "\n>>> Alice looks around the town square")
	alice.HandleCommand(game, "look")
	messages := getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice sees: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice checks her health")
	alice.HandleCommand(game, "health")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice goes to the dangerous forest")
	alice.HandleCommand(game, "south")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Bob goes to the marketplace")
	bob.HandleCommand(game, "east")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Alice attacks the giant rat")
	alice.HandleCommand(game, "attack giant rat")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Bob attacks the bandit")
	bob.HandleCommand(game, "fight bandit")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Waiting for monster AI to activate (3 seconds)...")
	time.Sleep(4 * time.Second)
	
	transcript = append(transcript, "\n>>> Checking Alice's condition after monster attacks")
	alice.HandleCommand(game, "health")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Checking Bob's condition after monster attacks")
	bob.HandleCommand(game, "health")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Alice tries to get the twisted branch")
	alice.HandleCommand(game, "get twisted branch")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice checks her inventory")
	alice.HandleCommand(game, "inventory")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Bob tries to examine the shiny coin")
	bob.HandleCommand(game, "examine shiny coin")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Alice equips the twisted branch weapon")
	alice.HandleCommand(game, "equip twisted branch")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice checks her equipment")
	alice.HandleCommand(game, "equipment")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice attacks again with weapon equipped")
	alice.HandleCommand(game, "attack giant rat")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice tries to flee back to town square")
	alice.HandleCommand(game, "north")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Bob gets the shiny coin and goes to tavern")
	bob.HandleCommand(game, "get shiny coin")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	bob.HandleCommand(game, "west")
	bob.HandleCommand(game, "north")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Bob explores the armory above the tavern")
	bob.HandleCommand(game, "up")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Bob gets the steel shield")
	bob.HandleCommand(game, "get steel shield")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Bob equips the armor")
	bob.HandleCommand(game, "equip steel shield")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Alice explores the temple and its underground")
	alice.HandleCommand(game, "west")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	alice.HandleCommand(game, "down")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice encounters the dangerous catacombs")
	alice.HandleCommand(game, "look")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice gets the leather armor")
	alice.HandleCommand(game, "get leather armor")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice equips armor and fights zombie")
	alice.HandleCommand(game, "equip leather armor")
	alice.HandleCommand(game, "attack shambling zombie")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice escapes to town square")
	alice.HandleCommand(game, "up")
	alice.HandleCommand(game, "east")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice explores the new pirate cove")
	alice.HandleCommand(game, "east")
	alice.HandleCommand(game, "south")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice examines the cutlass with enhanced examine")
	alice.HandleCommand(game, "examine cutlass")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice gets the cutlass and encounters sea monsters")
	alice.HandleCommand(game, "get cutlass")
	alice.HandleCommand(game, "attack bloodthirsty pirate")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice explores the volcanic cavern underground")
	alice.HandleCommand(game, "down")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice gets obsidian dagger and fights lava salamander")
	alice.HandleCommand(game, "get obsidian dagger")
	alice.HandleCommand(game, "equip obsidian dagger")
	alice.HandleCommand(game, "attack lava salamander")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice uses rest command to heal")
	alice.HandleCommand(game, "rest")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Bob goes to ice fortress via wizard tower")
	bob.HandleCommand(game, "down")
	bob.HandleCommand(game, "west")
	bob.HandleCommand(game, "down")
	bob.HandleCommand(game, "up")
	bob.HandleCommand(game, "up")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Bob gets frost armor and fights ice monsters")
	bob.HandleCommand(game, "get frost armor")
	bob.HandleCommand(game, "unequip steel shield")
	bob.HandleCommand(game, "equip frost armor")
	bob.HandleCommand(game, "attack frost yeti")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Bob explores sky temple and uses rest")
	bob.HandleCommand(game, "west")
	bob.HandleCommand(game, "rest")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Bob gets celestial blade and fights seraph")
	bob.HandleCommand(game, "get celestial blade")
	bob.HandleCommand(game, "equip celestial blade")
	bob.HandleCommand(game, "attack golden seraph")
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob: %s", msg))
	}
	clearPlayerMessages(bob)
	
	transcript = append(transcript, "\n>>> Alice goes to temple and uses prayer book for healing")
	alice.HandleCommand(game, "up")
	alice.HandleCommand(game, "north")
	alice.HandleCommand(game, "west")
	alice.HandleCommand(game, "get prayer book")
	alice.HandleCommand(game, "use prayer book")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Checking who is online from different locations")
	alice.HandleCommand(game, "who")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice says something about the expanded world")
	alice.HandleCommand(game, "say This world is amazing! I found sea monsters, lava creatures, and divine weapons!")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Final status check")
	transcript = append(transcript, fmt.Sprintf("Alice health: %d/%d", alice.health, alice.maxHealth))
	transcript = append(transcript, fmt.Sprintf("Bob health: %d/%d", bob.health, bob.maxHealth))
	transcript = append(transcript, fmt.Sprintf("Alice location: %s", alice.location.name))
	transcript = append(transcript, fmt.Sprintf("Bob location: %s", bob.location.name))
	
	transcriptContent := strings.Join(transcript, "\n")
	
	err := os.WriteFile("gameplay_transcript.txt", []byte(transcriptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write transcript file: %v", err)
	}
	
	t.Logf("Gameplay transcript written to gameplay_transcript.txt")
	t.Logf("Transcript contains %d lines", len(transcript))
	
	if alice.health <= 0 {
		t.Logf("Alice was killed and respawned during the test")
	}
	if bob.health <= 0 {
		t.Logf("Bob was killed and respawned during the test")
	}
}