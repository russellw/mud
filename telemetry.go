package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type TelemetryData struct {
	mu                    sync.RWMutex
	StartTime             time.Time            `json:"start_time"`
	TotalConnections      int64                `json:"total_connections"`
	ActiveConnections     int64                `json:"active_connections"`
	PlayersCreated        int64                `json:"players_created"`
	CommandsExecuted      int64                `json:"commands_executed"`
	CombatActions         int64                `json:"combat_actions"`
	MonsterKills          int64                `json:"monster_kills"`
	PlayerDeaths          int64                `json:"player_deaths"`
	RoomVisits            map[string]int64     `json:"room_visits"`
	CommandCounts         map[string]int64     `json:"command_counts"`
	LastUpdate            time.Time            `json:"last_update"`
}

type Telemetry struct {
	data   *TelemetryData
	logger *log.Logger
}

func NewTelemetry() *Telemetry {
	return &Telemetry{
		data: &TelemetryData{
			StartTime:     time.Now(),
			RoomVisits:    make(map[string]int64),
			CommandCounts: make(map[string]int64),
			LastUpdate:    time.Now(),
		},
		logger: log.New(log.Writer(), "[TELEMETRY] ", log.LstdFlags),
	}
}

func (t *Telemetry) IncrementConnections() {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()
	t.data.TotalConnections++
	t.data.ActiveConnections++
	t.data.LastUpdate = time.Now()
}

func (t *Telemetry) DecrementActiveConnections() {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()
	if t.data.ActiveConnections > 0 {
		t.data.ActiveConnections--
	}
	t.data.LastUpdate = time.Now()
}

func (t *Telemetry) IncrementPlayersCreated() {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()
	t.data.PlayersCreated++
	t.data.LastUpdate = time.Now()
}

func (t *Telemetry) IncrementCommandsExecuted(command string) {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()
	t.data.CommandsExecuted++
	t.data.CommandCounts[command]++
	t.data.LastUpdate = time.Now()
}

func (t *Telemetry) IncrementCombatActions() {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()
	t.data.CombatActions++
	t.data.LastUpdate = time.Now()
}

func (t *Telemetry) IncrementMonsterKills() {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()
	t.data.MonsterKills++
	t.data.LastUpdate = time.Now()
}

func (t *Telemetry) IncrementPlayerDeaths() {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()
	t.data.PlayerDeaths++
	t.data.LastUpdate = time.Now()
}

func (t *Telemetry) RecordRoomVisit(roomName string) {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()
	t.data.RoomVisits[roomName]++
	t.data.LastUpdate = time.Now()
}

func (t *Telemetry) GetSnapshot() TelemetryData {
	t.data.mu.RLock()
	defer t.data.mu.RUnlock()
	
	snapshot := TelemetryData{
		StartTime:         t.data.StartTime,
		TotalConnections:  t.data.TotalConnections,
		ActiveConnections: t.data.ActiveConnections,
		PlayersCreated:    t.data.PlayersCreated,
		CommandsExecuted:  t.data.CommandsExecuted,
		CombatActions:     t.data.CombatActions,
		MonsterKills:      t.data.MonsterKills,
		PlayerDeaths:      t.data.PlayerDeaths,
		RoomVisits:        make(map[string]int64),
		CommandCounts:     make(map[string]int64),
		LastUpdate:        t.data.LastUpdate,
	}
	
	for k, v := range t.data.RoomVisits {
		snapshot.RoomVisits[k] = v
	}
	
	for k, v := range t.data.CommandCounts {
		snapshot.CommandCounts[k] = v
	}
	
	return snapshot
}

func (t *Telemetry) GetJSON() ([]byte, error) {
	snapshot := t.GetSnapshot()
	return json.MarshalIndent(snapshot, "", "  ")
}

func (t *Telemetry) LogStats() {
	uptime := time.Since(t.data.StartTime)
	snapshot := t.GetSnapshot()
	
	t.logger.Printf("=== MUD Server Statistics ===")
	t.logger.Printf("Uptime: %v", uptime.Round(time.Second))
	t.logger.Printf("Total Connections: %d", snapshot.TotalConnections)
	t.logger.Printf("Active Connections: %d", snapshot.ActiveConnections)
	t.logger.Printf("Players Created: %d", snapshot.PlayersCreated)
	t.logger.Printf("Commands Executed: %d", snapshot.CommandsExecuted)
	t.logger.Printf("Combat Actions: %d", snapshot.CombatActions)
	t.logger.Printf("Monster Kills: %d", snapshot.MonsterKills)
	t.logger.Printf("Player Deaths: %d", snapshot.PlayerDeaths)
	
	if len(snapshot.RoomVisits) > 0 {
		t.logger.Printf("Most Popular Rooms:")
		for room, visits := range snapshot.RoomVisits {
			if visits > 0 {
				t.logger.Printf("  %s: %d visits", room, visits)
			}
		}
	}
	
	if len(snapshot.CommandCounts) > 0 {
		t.logger.Printf("Top Commands:")
		for cmd, count := range snapshot.CommandCounts {
			if count > 0 {
				t.logger.Printf("  %s: %d times", cmd, count)
			}
		}
	}
	t.logger.Printf("============================")
}

func (t *Telemetry) StartPeriodicLogging(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			t.LogStats()
		}
	}()
}

var GlobalTelemetry = NewTelemetry()

func (t *Telemetry) GetSummary() string {
	snapshot := t.GetSnapshot()
	uptime := time.Since(snapshot.StartTime)
	
	return fmt.Sprintf("Uptime: %v | Connections: %d/%d | Commands: %d | Combat: %d | Kills: %d | Deaths: %d",
		uptime.Round(time.Second),
		snapshot.ActiveConnections,
		snapshot.TotalConnections,
		snapshot.CommandsExecuted,
		snapshot.CombatActions,
		snapshot.MonsterKills,
		snapshot.PlayerDeaths)
}