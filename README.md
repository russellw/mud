# MUD - Multi-User Dungeon

A fully-featured text-based MMORPG (Multi-User Dungeon) server written in Go, featuring real-time combat, equipment systems, and an expansive fantasy world.

## Features

### 🌍 Expansive World
- **20 interconnected rooms** across 4 levels (-1 to 2)
- Diverse themed areas: pirate coves, volcanic caverns, ice fortresses, sky temples, haunted libraries, goblin warrens
- Complex multi-level navigation with vertical movement (up/down)

### ⚔️ Combat System
- **Real-time monster AI** with 3-second tick cycles
- **30 unique monsters** with varied behaviors (aggressive vs defensive)
- Dynamic damage calculation with equipment bonuses
- Player death and respawn mechanics

### 🛡️ Equipment & Items
- **Weapons**: Range from twisted branch (+2 damage) to celestial blade (+12 damage)
- **Armor**: Leather armor (+3 defense) to frost armor (+7 defense)
- **Consumables**: Tome of knowledge (permanent stat boost), prayer book (full heal)
- Equipment inspection showing detailed stats

### 🎮 Player Commands
- **Movement**: `go <direction>`, `up`, `down`
- **Combat**: `attack <monster>`, `fight <monster>`
- **Items**: `get <item>`, `drop <item>`, `examine <item>`, `inventory`
- **Equipment**: `equip <item>`, `unequip <item>`, `equipment`
- **Special**: `use <item>`, `rest`, `health`, `who`, `say <message>`

### 🏗️ Technical Features
- **Telnet protocol** for authentic MUD experience (port 4000)
- **Concurrent connection handling** with goroutines
- **Centralized game logic** with event-driven monster AI
- **Comprehensive test suite** with mock connections
- **PNG map generation** for world visualization

## Quick Start

```bash
# Build and run the server
go build
./mud

# Connect via telnet
telnet localhost 4000
```

## World Map

The game features a complex 4-level world structure:

- **Level 2**: Sky Temple, Ice Fortress (divine/celestial areas)
- **Level 1**: Wizard Tower, Castle Armory, Goblin Warren (elevated areas)  
- **Level 0**: Town Square, Forest, Market, Temple (main world)
- **Level -1**: Catacombs, Dungeon, Volcanic Cavern, Crystal Mines (underground)

Navigate between levels using `up` and `down` commands at connected rooms.

## Development

### Testing
```bash
# Run all tests
go test

# Generate gameplay transcript
go test -run TestGameplayTranscript

# Generate world map
go test -run TestGenerateMapPNG
```

### Architecture
- `main.go` - Network server and connection handling
- `game.go` - Core game logic and world creation
- `player.go` - Player commands and actions
- `room.go` - Room structures and broadcasting
- `monster.go` - Monster AI and behavior
- `*_test.go` - Comprehensive test suite

## Game Statistics

- **Rooms**: 20 interconnected areas
- **Monsters**: 30 unique creatures
- **Weapons**: 8 different weapon types (2-12 damage)
- **Armor**: 6 armor pieces (3-8 defense)
- **Levels**: 4-level vertical world structure
- **Commands**: 15+ player commands

Built with Go 1.24, tested extensively with automated gameplay scenarios.
