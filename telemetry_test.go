package main

import (
	"testing"
	"time"
)

func TestTelemetryBasicFunctionality(t *testing.T) {
	telemetry := NewTelemetry()
	
	// Test initial state
	snapshot := telemetry.GetSnapshot()
	if snapshot.TotalConnections != 0 {
		t.Errorf("Expected 0 total connections, got %d", snapshot.TotalConnections)
	}
	if snapshot.CommandsExecuted != 0 {
		t.Errorf("Expected 0 commands executed, got %d", snapshot.CommandsExecuted)
	}
	
	// Test incrementing connections
	telemetry.IncrementConnections()
	telemetry.IncrementConnections()
	snapshot = telemetry.GetSnapshot()
	if snapshot.TotalConnections != 2 {
		t.Errorf("Expected 2 total connections, got %d", snapshot.TotalConnections)
	}
	if snapshot.ActiveConnections != 2 {
		t.Errorf("Expected 2 active connections, got %d", snapshot.ActiveConnections)
	}
	
	// Test decrementing active connections
	telemetry.DecrementActiveConnections()
	snapshot = telemetry.GetSnapshot()
	if snapshot.ActiveConnections != 1 {
		t.Errorf("Expected 1 active connection, got %d", snapshot.ActiveConnections)
	}
	
	// Test command tracking
	telemetry.IncrementCommandsExecuted("look")
	telemetry.IncrementCommandsExecuted("go")
	telemetry.IncrementCommandsExecuted("look")
	
	snapshot = telemetry.GetSnapshot()
	if snapshot.CommandsExecuted != 3 {
		t.Errorf("Expected 3 commands executed, got %d", snapshot.CommandsExecuted)
	}
	if snapshot.CommandCounts["look"] != 2 {
		t.Errorf("Expected 2 'look' commands, got %d", snapshot.CommandCounts["look"])
	}
	if snapshot.CommandCounts["go"] != 1 {
		t.Errorf("Expected 1 'go' command, got %d", snapshot.CommandCounts["go"])
	}
	
	// Test room visits
	telemetry.RecordRoomVisit("Town Square")
	telemetry.RecordRoomVisit("Forest")
	telemetry.RecordRoomVisit("Town Square")
	
	snapshot = telemetry.GetSnapshot()
	if snapshot.RoomVisits["Town Square"] != 2 {
		t.Errorf("Expected 2 visits to Town Square, got %d", snapshot.RoomVisits["Town Square"])
	}
	if snapshot.RoomVisits["Forest"] != 1 {
		t.Errorf("Expected 1 visit to Forest, got %d", snapshot.RoomVisits["Forest"])
	}
	
	// Test combat tracking
	telemetry.IncrementCombatActions()
	telemetry.IncrementMonsterKills()
	telemetry.IncrementPlayerDeaths()
	
	snapshot = telemetry.GetSnapshot()
	if snapshot.CombatActions != 1 {
		t.Errorf("Expected 1 combat action, got %d", snapshot.CombatActions)
	}
	if snapshot.MonsterKills != 1 {
		t.Errorf("Expected 1 monster kill, got %d", snapshot.MonsterKills)
	}
	if snapshot.PlayerDeaths != 1 {
		t.Errorf("Expected 1 player death, got %d", snapshot.PlayerDeaths)
	}
}

func TestTelemetryJSON(t *testing.T) {
	telemetry := NewTelemetry()
	
	telemetry.IncrementConnections()
	telemetry.IncrementCommandsExecuted("test")
	
	jsonData, err := telemetry.GetJSON()
	if err != nil {
		t.Errorf("Failed to get JSON: %v", err)
	}
	
	if len(jsonData) == 0 {
		t.Error("JSON data is empty")
	}
}

func TestTelemetrySummary(t *testing.T) {
	telemetry := NewTelemetry()
	
	telemetry.IncrementConnections()
	telemetry.IncrementCommandsExecuted("test")
	
	summary := telemetry.GetSummary()
	if len(summary) == 0 {
		t.Error("Summary is empty")
	}
	
	// Summary should contain key metrics
	if !contains(summary, "Connections") {
		t.Error("Summary missing connections info")
	}
	if !contains(summary, "Commands") {
		t.Error("Summary missing commands info")
	}
}

func TestTelemetryStartTime(t *testing.T) {
	before := time.Now()
	telemetry := NewTelemetry()
	after := time.Now()
	
	snapshot := telemetry.GetSnapshot()
	if snapshot.StartTime.Before(before) || snapshot.StartTime.After(after) {
		t.Error("Start time not properly set")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}