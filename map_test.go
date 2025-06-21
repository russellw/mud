package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type RoomPosition struct {
	x, y, z int
}

var roomPositions = map[string]RoomPosition{
	"town_square":  {2, 2, 0},
	"tavern":       {2, 1, 0},
	"forest":       {2, 3, 0},
	"market":       {3, 2, 0},
	"temple":       {1, 2, 0},
	"armory":       {2, 1, 1},
	"deep_forest":  {2, 4, 0},
	"dungeon":      {2, 3, -1},
	"catacombs":    {1, 2, -1},
	"wizard_tower": {1, 2, 1},
	"dragon_lair":  {3, 4, 0},
	"cemetery":     {1, 4, 0},
}

func TestRoomConnectivity(t *testing.T) {
	game := NewGame()
	
	// Test all rooms exist
	expectedRooms := []string{
		"town_square", "tavern", "forest", "market", "temple",
		"dungeon", "deep_forest", "catacombs", "wizard_tower", 
		"dragon_lair", "cemetery", "armory",
	}
	
	for _, roomKey := range expectedRooms {
		if _, exists := game.rooms[roomKey]; !exists {
			t.Errorf("Expected room '%s' not found", roomKey)
		}
	}
	
	// Test bidirectional connectivity
	connectivityTests := []struct {
		from, to, direction string
	}{
		{"town_square", "tavern", "north"},
		{"town_square", "forest", "south"},
		{"town_square", "market", "east"},
		{"town_square", "temple", "west"},
		{"tavern", "armory", "up"},
		{"forest", "deep_forest", "south"},
		{"forest", "dungeon", "down"},
		{"temple", "catacombs", "down"},
		{"temple", "wizard_tower", "up"},
		{"deep_forest", "cemetery", "west"},
		{"deep_forest", "dragon_lair", "east"},
		{"dungeon", "catacombs", "north"},
	}
	
	for _, test := range connectivityTests {
		fromRoom := game.rooms[test.from]
		toRoom := game.rooms[test.to]
		
		if fromRoom == nil {
			t.Errorf("From room '%s' not found", test.from)
			continue
		}
		if toRoom == nil {
			t.Errorf("To room '%s' not found", test.to)
			continue
		}
		
		// Test forward connection
		if exit, exists := fromRoom.exits[test.direction]; !exists {
			t.Errorf("Room '%s' missing '%s' exit", test.from, test.direction)
		} else if exit != toRoom {
			t.Errorf("Room '%s' '%s' exit points to wrong room", test.from, test.direction)
		}
		
		// Test reverse connection
		reverseDir := getReverseDirection(test.direction)
		if reverseExit, exists := toRoom.exits[reverseDir]; !exists {
			t.Errorf("Room '%s' missing reverse '%s' exit", test.to, reverseDir)
		} else if reverseExit != fromRoom {
			t.Errorf("Room '%s' '%s' exit doesn't point back to '%s'", test.to, reverseDir, test.from)
		}
	}
	
	t.Logf("All room connectivity tests passed")
}

func getReverseDirection(dir string) string {
	reverseMap := map[string]string{
		"north": "south",
		"south": "north",
		"east":  "west",
		"west":  "east",
		"up":    "down",
		"down":  "up",
	}
	return reverseMap[dir]
}

func TestRoomContents(t *testing.T) {
	game := NewGame()
	
	// Test specific room contents
	contentTests := []struct {
		room         string
		expectedItem string
		itemType     string
		monsters     []string
	}{
		{"tavern", "wooden mug", "misc", []string{}},
		{"forest", "twisted branch", "weapon", []string{"giant rat", "dire wolf"}},
		{"market", "shiny coin", "misc", []string{"bandit"}},
		{"temple", "prayer book", "misc", []string{"shadowy cultist"}},
		{"dungeon", "rusty key", "misc", []string{"skeleton warrior"}},
		{"deep_forest", "iron sword", "weapon", []string{"cave bear"}},
		{"catacombs", "leather armor", "armor", []string{"shambling zombie", "ancient mummy"}},
		{"wizard_tower", "magic staff", "weapon", []string{"fire imp"}},
		{"dragon_lair", "dragon scale", "armor", []string{"ancient dragon"}},
		{"cemetery", "silver cross", "weapon", []string{"wandering spirit", "vengeful wraith"}},
		{"armory", "steel shield", "armor", []string{"castle guard"}},
	}
	
	for _, test := range contentTests {
		room := game.rooms[test.room]
		if room == nil {
			t.Errorf("Room '%s' not found", test.room)
			continue
		}
		
		// Test items
		if len(room.items) == 0 {
			if test.expectedItem != "" {
				t.Errorf("Room '%s' should have item '%s' but has no items", test.room, test.expectedItem)
			}
		} else {
			found := false
			for _, item := range room.items {
				if item.name == test.expectedItem {
					found = true
					if item.itemType != test.itemType {
						t.Errorf("Item '%s' in room '%s' should be type '%s' but is '%s'", 
							test.expectedItem, test.room, test.itemType, item.itemType)
					}
					break
				}
			}
			if !found && test.expectedItem != "" {
				t.Errorf("Room '%s' missing expected item '%s'", test.room, test.expectedItem)
			}
		}
		
		// Test monsters
		if len(room.monsters) != len(test.monsters) {
			t.Errorf("Room '%s' should have %d monsters but has %d", 
				test.room, len(test.monsters), len(room.monsters))
		} else {
			for _, expectedMonster := range test.monsters {
				found := false
				for _, monster := range room.monsters {
					if monster.name == expectedMonster {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Room '%s' missing expected monster '%s'", test.room, expectedMonster)
				}
			}
		}
	}
	
	t.Logf("All room content tests passed")
}

func TestGenerateMapPNG(t *testing.T) {
	game := NewGame()
	
	// Create image
	const (
		cellSize = 120
		margin   = 50
		mapWidth = 5
		mapHeight = 5
		levels   = 3 // -1, 0, 1
	)
	
	imgWidth := mapWidth*cellSize + 2*margin + 200 // extra space for legend
	imgHeight := levels*mapHeight*cellSize + 2*margin + 100
	
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	
	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{240, 240, 240, 255}}, image.Point{}, draw.Src)
	
	// Colors for different room types
	colors := map[string]color.RGBA{
		"town_square":  {100, 200, 100, 255}, // Green
		"tavern":       {200, 150, 100, 255}, // Brown
		"forest":       {50, 150, 50, 255},   // Dark green
		"market":       {200, 200, 100, 255}, // Yellow
		"temple":       {150, 150, 200, 255}, // Light blue
		"armory":       {150, 150, 150, 255}, // Gray
		"deep_forest":  {30, 100, 30, 255},   // Very dark green
		"dungeon":      {100, 50, 50, 255},   // Dark red
		"catacombs":    {80, 80, 80, 255},    // Dark gray
		"wizard_tower": {150, 100, 200, 255}, // Purple
		"dragon_lair":  {200, 50, 50, 255},   // Red
		"cemetery":     {120, 120, 120, 255}, // Medium gray
	}
	
	// Draw rooms for each level
	for level := -1; level <= 1; level++ {
		levelOffset := (level + 1) * (mapHeight * cellSize + 30)
		
		// Draw level label
		drawText(img, fmt.Sprintf("Level %d", level), margin-30, margin+levelOffset-10, color.RGBA{0, 0, 0, 255})
		
		for roomKey, pos := range roomPositions {
			if pos.z != level {
				continue
			}
			
			room := game.rooms[roomKey]
			if room == nil {
				continue
			}
			
			x := margin + pos.x*cellSize
			y := margin + levelOffset + pos.y*cellSize
			
			// Draw room rectangle
			roomColor := colors[roomKey]
			if roomColor.A == 0 {
				roomColor = color.RGBA{200, 200, 200, 255} // Default gray
			}
			
			drawRect(img, x, y, cellSize-10, cellSize-10, roomColor)
			
			// Draw room name
			lines := splitRoomName(room.name)
			for i, line := range lines {
				drawText(img, line, x+5, y+15+i*12, color.RGBA{0, 0, 0, 255})
			}
			
			// Draw monster count
			if len(room.monsters) > 0 {
				drawText(img, fmt.Sprintf("M:%d", len(room.monsters)), x+5, y+cellSize-25, color.RGBA{200, 0, 0, 255})
			}
			
			// Draw item indicator
			if len(room.items) > 0 {
				drawText(img, fmt.Sprintf("I:%d", len(room.items)), x+5, y+cellSize-12, color.RGBA{0, 0, 200, 255})
			}
		}
		
		// Draw connections for this level
		for roomKey, pos := range roomPositions {
			if pos.z != level {
				continue
			}
			
			room := game.rooms[roomKey]
			if room == nil {
				continue
			}
			
			x1 := margin + pos.x*cellSize + cellSize/2
			y1 := margin + levelOffset + pos.y*cellSize + cellSize/2
			
			for dir, connectedRoom := range room.exits {
				connectedKey := ""
				for key, r := range game.rooms {
					if r == connectedRoom {
						connectedKey = key
						break
					}
				}
				
				if connectedPos, exists := roomPositions[connectedKey]; exists && connectedPos.z == level {
					x2 := margin + connectedPos.x*cellSize + cellSize/2
					y2 := margin + levelOffset + connectedPos.y*cellSize + cellSize/2
					
					// Draw connection line
					drawLine(img, x1, y1, x2, y2, color.RGBA{100, 100, 100, 255})
					
					// Draw direction label
					midX := (x1 + x2) / 2
					midY := (y1 + y2) / 2
					drawText(img, strings.ToUpper(dir[:1]), midX-3, midY+3, color.RGBA{50, 50, 50, 255})
				}
			}
		}
	}
	
	// Draw legend
	legendX := mapWidth*cellSize + margin + 20
	legendY := margin
	
	drawText(img, "MUD World Map", legendX, legendY, color.RGBA{0, 0, 0, 255})
	legendY += 25
	
	drawText(img, "Legend:", legendX, legendY, color.RGBA{0, 0, 0, 255})
	legendY += 20
	
	drawText(img, "M:X = X Monsters", legendX, legendY, color.RGBA{200, 0, 0, 255})
	legendY += 15
	drawText(img, "I:X = X Items", legendX, legendY, color.RGBA{0, 0, 200, 255})
	legendY += 25
	
	// Room type legend
	legendEntries := []struct {
		name  string
		color color.RGBA
	}{
		{"Town Square", colors["town_square"]},
		{"Tavern", colors["tavern"]},
		{"Forest Areas", colors["forest"]},
		{"Market", colors["market"]},
		{"Temple", colors["temple"]},
		{"Underground", colors["dungeon"]},
		{"Dragon Lair", colors["dragon_lair"]},
	}
	
	for _, entry := range legendEntries {
		drawRect(img, legendX, legendY-8, 15, 10, entry.color)
		drawText(img, entry.name, legendX+20, legendY, color.RGBA{0, 0, 0, 255})
		legendY += 18
	}
	
	// Save PNG
	file, err := os.Create("mud_map.png")
	if err != nil {
		t.Fatalf("Failed to create map file: %v", err)
	}
	defer file.Close()
	
	err = png.Encode(file, img)
	if err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}
	
	t.Logf("Map saved to mud_map.png")
	t.Logf("Map dimensions: %dx%d pixels", imgWidth, imgHeight)
	t.Logf("Rooms mapped: %d", len(roomPositions))
}

func drawRect(img *image.RGBA, x, y, width, height int, c color.RGBA) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			if x+dx < img.Bounds().Max.X && y+dy < img.Bounds().Max.Y {
				img.Set(x+dx, y+dy, c)
			}
		}
	}
	
	// Draw border
	borderColor := color.RGBA{0, 0, 0, 255}
	for dx := 0; dx < width; dx++ {
		if x+dx < img.Bounds().Max.X {
			if y < img.Bounds().Max.Y {
				img.Set(x+dx, y, borderColor)
			}
			if y+height-1 < img.Bounds().Max.Y {
				img.Set(x+dx, y+height-1, borderColor)
			}
		}
	}
	for dy := 0; dy < height; dy++ {
		if y+dy < img.Bounds().Max.Y {
			if x < img.Bounds().Max.X {
				img.Set(x, y+dy, borderColor)
			}
			if x+width-1 < img.Bounds().Max.X {
				img.Set(x+width-1, y+dy, borderColor)
			}
		}
	}
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	
	var sx, sy int
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}
	
	err := dx - dy
	x, y := x1, y1
	
	for {
		if x >= 0 && x < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
			img.Set(x, y, c)
		}
		
		if x == x2 && y == y2 {
			break
		}
		
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func drawText(img *image.RGBA, text string, x, y int, c color.RGBA) {
	face := basicfont.Face7x13
	
	drawer := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{c},
		Face: face,
		Dot:  fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)},
	}
	
	drawer.DrawString(text)
}

func splitRoomName(name string) []string {
	words := strings.Fields(name)
	if len(words) <= 2 {
		return []string{name}
	}
	
	// Split long names into multiple lines
	if len(name) > 15 {
		mid := len(words) / 2
		line1 := strings.Join(words[:mid], " ")
		line2 := strings.Join(words[mid:], " ")
		return []string{line1, line2}
	}
	
	return []string{name}
}