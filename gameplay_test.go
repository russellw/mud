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
	
	transcript = append(transcript, "\n>>> Alice tries to flee back to town square")
	alice.HandleCommand(game, "north")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Checking who is online")
	alice.HandleCommand(game, "who")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	clearPlayerMessages(alice)
	
	transcript = append(transcript, "\n>>> Alice says something to Bob")
	alice.HandleCommand(game, "say That was scary! The monsters are really aggressive.")
	messages = getPlayerMessages(alice)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Alice: %s", msg))
	}
	
	messages = getPlayerMessages(bob)
	for _, msg := range messages {
		transcript = append(transcript, fmt.Sprintf("Bob hears: %s", msg))
	}
	clearPlayerMessages(alice)
	clearPlayerMessages(bob)
	
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