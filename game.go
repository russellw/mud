package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	rooms   map[string]*Room
	players []*Player
	running bool
}

func NewGame() *Game {
	game := &Game{
		rooms:   make(map[string]*Room),
		players: make([]*Player, 0),
		running: true,
	}
	
	game.createWorld()
	go game.gameLoop()
	
	return game
}

func (g *Game) createWorld() {
	townSquare := &Room{
		name:        "Town Square",
		description: "You are standing in a bustling town square. There are paths leading in all directions.",
		players:     make([]*Player, 0),
		items:       make([]*Item, 0),
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	tavern := &Room{
		name:        "The Prancing Pony Tavern",
		description: "A cozy tavern filled with the smell of ale and roasted meat. Wooden tables and chairs are scattered around.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "wooden mug", description: "A sturdy wooden drinking mug", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	forest := &Room{
		name:        "Dark Forest",
		description: "A dense forest with towering trees that block most of the sunlight. Strange sounds echo from the shadows.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "twisted branch", description: "A gnarled branch that could serve as a walking stick", itemType: "weapon", damage: 2, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	market := &Room{
		name:        "Marketplace",
		description: "A busy marketplace with merchants hawking their wares. Colorful stalls line the cobblestone square.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "shiny coin", description: "A gold coin that glints in the sunlight", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	temple := &Room{
		name:        "Ancient Temple",
		description: "A sacred temple with marble columns and intricate carvings. A sense of peace fills the air.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "prayer book", description: "An old leather-bound book of prayers and rituals", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	dungeon := &Room{
		name:        "Dungeon Entrance",
		description: "A crumbling stone entrance leads into darkness. Ancient torches flicker on the walls.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "rusty key", description: "An old iron key, corroded with age", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	deepForest := &Room{
		name:        "Deep Forest",
		description: "The forest grows darker here. Thick canopy blocks all sunlight. Something large moves in the shadows.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "iron sword", description: "A well-forged iron sword with a sharp edge", itemType: "weapon", damage: 8, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	catacombs := &Room{
		name:        "Ancient Catacombs",
		description: "Narrow stone passages wind through countless burial chambers. The air is thick with age and mystery.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "leather armor", description: "Sturdy leather armor that provides good protection", itemType: "armor", damage: 0, defense: 3}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	wizardTower := &Room{
		name:        "Wizard's Tower",
		description: "A tall stone tower filled with magical artifacts and glowing crystals. Books float in mid-air.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "magic staff", description: "A wooden staff topped with a glowing crystal orb", itemType: "weapon", damage: 10, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	dragonLair := &Room{
		name:        "Dragon's Lair",
		description: "A massive cavern with piles of gold and treasure. Scorch marks cover the walls. The air shimmers with heat.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "dragon scale", description: "A massive golden scale, still warm to the touch", itemType: "armor", damage: 0, defense: 8}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	cemetery := &Room{
		name:        "Moonlit Cemetery",
		description: "Ancient gravestones stretch as far as you can see. Mist swirls between the weathered monuments.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "silver cross", description: "A blessed silver cross that gleams in the moonlight", itemType: "weapon", damage: 6, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	armory := &Room{
		name:        "Castle Armory",
		description: "Weapons and armor line the walls of this military storehouse. Everything is kept in perfect condition.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "steel shield", description: "A heavy steel shield emblazoned with a royal crest", itemType: "armor", damage: 0, defense: 5}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// New areas - Pirate Cove
	pirateCove := &Room{
		name:        "Hidden Pirate Cove",
		description: "A secluded beach cove with a rotting wooden pier. Seagulls cry overhead and waves crash against the rocky shore.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "cutlass", description: "A curved pirate sword with a brass handguard", itemType: "weapon", damage: 7, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Volcanic Caverns
	volcanicCavern := &Room{
		name:        "Volcanic Cavern",
		description: "A steaming cavern deep underground. Lava pools cast an orange glow on the obsidian walls. The air shimmers with heat.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "obsidian dagger", description: "A razor-sharp dagger carved from volcanic glass", itemType: "weapon", damage: 9, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Ice Fortress
	iceFortress := &Room{
		name:        "Frozen Fortress",
		description: "An ancient fortress made entirely of ice and snow. Icicles hang like spears from the ceiling. Your breath forms clouds in the frigid air.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "frost armor", description: "Crystalline armor that radiates cold, providing excellent protection", itemType: "armor", damage: 0, defense: 7}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Sky Temple
	skyTemple := &Room{
		name:        "Sky Temple",
		description: "A magnificent temple floating high in the clouds. Golden columns support a crystal dome that captures the sunlight.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "celestial blade", description: "A legendary sword that glows with divine light", itemType: "weapon", damage: 12, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Cursed Swamp
	cursedSwamp := &Room{
		name:        "Cursed Swamp",
		description: "A fetid swamp where twisted trees emerge from stagnant water. Strange lights flicker in the mist and the air reeks of decay.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "swamp boots", description: "Waterproof boots that protect against poison and disease", itemType: "armor", damage: 0, defense: 4}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Crystal Mines
	crystalMines := &Room{
		name:        "Crystal Mines",
		description: "Deep underground tunnels where precious crystals grow from the walls. The gems cast rainbow patterns of light throughout the cavern.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "crystal wand", description: "A wand topped with a multifaceted crystal that pulses with magical energy", itemType: "weapon", damage: 11, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Haunted Library
	hauntedLibrary := &Room{
		name:        "Haunted Library",
		description: "A vast library with towering shelves of ancient books. Spectral figures drift between the stacks and whispers echo in the darkness.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "tome of knowledge", description: "An ancient book that increases the reader's wisdom and magical understanding", itemType: "misc", damage: 0, defense: 0}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Goblin Warren
	goblinWarren := &Room{
		name:        "Goblin Warren",
		description: "A maze of tunnels and chambers carved into the hillside. The walls are covered in crude goblin drawings and the floor is littered with bones.",
		players:     make([]*Player, 0),
		items:       []*Item{{name: "goblin mail", description: "Crude but effective armor made from scavenged metal pieces", itemType: "armor", damage: 0, defense: 6}},
		monsters:    make([]*Monster, 0),
		exits:       make(map[string]*Room),
	}
	
	// Connect main areas
	townSquare.exits["north"] = tavern
	townSquare.exits["south"] = forest
	townSquare.exits["east"] = market
	townSquare.exits["west"] = temple
	
	tavern.exits["south"] = townSquare
	tavern.exits["up"] = armory
	
	forest.exits["north"] = townSquare
	forest.exits["south"] = deepForest
	forest.exits["down"] = dungeon
	
	market.exits["west"] = townSquare
	
	temple.exits["east"] = townSquare
	temple.exits["down"] = catacombs
	temple.exits["up"] = wizardTower
	
	// Connect extended areas
	dungeon.exits["up"] = forest
	dungeon.exits["north"] = catacombs
	
	deepForest.exits["north"] = forest
	deepForest.exits["west"] = cemetery
	deepForest.exits["east"] = dragonLair
	
	catacombs.exits["up"] = temple
	catacombs.exits["south"] = dungeon
	
	wizardTower.exits["down"] = temple
	
	dragonLair.exits["west"] = deepForest
	
	cemetery.exits["east"] = deepForest
	
	armory.exits["down"] = tavern
	
	// Connect new areas
	market.exits["south"] = pirateCove
	pirateCove.exits["north"] = market
	pirateCove.exits["down"] = volcanicCavern
	
	volcanicCavern.exits["up"] = pirateCove
	volcanicCavern.exits["north"] = crystalMines
	
	crystalMines.exits["south"] = volcanicCavern
	crystalMines.exits["up"] = goblinWarren
	
	goblinWarren.exits["down"] = crystalMines
	goblinWarren.exits["west"] = cursedSwamp
	
	cursedSwamp.exits["east"] = goblinWarren
	cursedSwamp.exits["north"] = hauntedLibrary
	
	hauntedLibrary.exits["south"] = cursedSwamp
	hauntedLibrary.exits["up"] = skyTemple
	
	skyTemple.exits["down"] = hauntedLibrary
	skyTemple.exits["east"] = iceFortress
	
	iceFortress.exits["west"] = skyTemple
	iceFortress.exits["down"] = wizardTower
	
	wizardTower.exits["up"] = iceFortress
	
	g.rooms["town_square"] = townSquare
	g.rooms["tavern"] = tavern
	g.rooms["forest"] = forest
	g.rooms["market"] = market
	g.rooms["temple"] = temple
	g.rooms["dungeon"] = dungeon
	g.rooms["deep_forest"] = deepForest
	g.rooms["catacombs"] = catacombs
	g.rooms["wizard_tower"] = wizardTower
	g.rooms["dragon_lair"] = dragonLair
	g.rooms["cemetery"] = cemetery
	g.rooms["armory"] = armory
	g.rooms["pirate_cove"] = pirateCove
	g.rooms["volcanic_cavern"] = volcanicCavern
	g.rooms["ice_fortress"] = iceFortress
	g.rooms["sky_temple"] = skyTemple
	g.rooms["cursed_swamp"] = cursedSwamp
	g.rooms["crystal_mines"] = crystalMines
	g.rooms["haunted_library"] = hauntedLibrary
	g.rooms["goblin_warren"] = goblinWarren
	
	g.spawnMonsters()
}

func (g *Game) AddPlayer(player *Player) {
	g.players = append(g.players, player)
	startRoom := g.rooms["town_square"]
	player.location = startRoom
	startRoom.players = append(startRoom.players, player)
}

func (g *Game) RemovePlayer(player *Player) {
	for i, p := range g.players {
		if p == player {
			g.players = append(g.players[:i], g.players[i+1:]...)
			break
		}
	}
	
	if player.location != nil {
		for i, p := range player.location.players {
			if p == player {
				player.location.players = append(player.location.players[:i], player.location.players[i+1:]...)
				break
			}
		}
	}
}

func (g *Game) spawnMonsters() {
	// Original forest monsters
	rat := NewMonster("giant rat", "A large, mangy rat with red eyes and yellowed teeth", 15, 3, true)
	rat.location = g.rooms["forest"]
	g.rooms["forest"].monsters = append(g.rooms["forest"].monsters, rat)
	
	wolf := NewMonster("dire wolf", "A massive wolf with silver fur and piercing blue eyes", 25, 6, true)
	wolf.location = g.rooms["forest"]
	g.rooms["forest"].monsters = append(g.rooms["forest"].monsters, wolf)
	
	// Market monsters
	bandit := NewMonster("bandit", "A shifty-looking human in leather armor, clutching a rusty dagger", 20, 5, true)
	bandit.location = g.rooms["market"]
	g.rooms["market"].monsters = append(g.rooms["market"].monsters, bandit)
	
	// Temple monsters
	cultist := NewMonster("shadowy cultist", "A robed figure with glowing red eyes chanting in an unknown language", 18, 4, true)
	cultist.location = g.rooms["temple"]
	g.rooms["temple"].monsters = append(g.rooms["temple"].monsters, cultist)
	
	// Deep forest monsters
	bear := NewMonster("cave bear", "A massive brown bear with razor-sharp claws and a thunderous roar", 40, 8, true)
	bear.location = g.rooms["deep_forest"]
	g.rooms["deep_forest"].monsters = append(g.rooms["deep_forest"].monsters, bear)
	
	// Dungeon monsters
	skeleton := NewMonster("skeleton warrior", "An ancient skeleton in rusted armor, wielding a bone sword", 20, 5, false)
	skeleton.location = g.rooms["dungeon"]
	g.rooms["dungeon"].monsters = append(g.rooms["dungeon"].monsters, skeleton)
	
	// Catacombs monsters
	zombie := NewMonster("shambling zombie", "A rotting corpse that moves with unnatural hunger", 25, 4, true)
	zombie.location = g.rooms["catacombs"]
	g.rooms["catacombs"].monsters = append(g.rooms["catacombs"].monsters, zombie)
	
	mummy := NewMonster("ancient mummy", "Wrapped in decaying bandages, this ancient guardian protects the tombs", 30, 6, false)
	mummy.location = g.rooms["catacombs"]
	g.rooms["catacombs"].monsters = append(g.rooms["catacombs"].monsters, mummy)
	
	// Wizard tower monsters
	imp := NewMonster("fire imp", "A small demonic creature wreathed in flames with a mischievous grin", 15, 7, true)
	imp.location = g.rooms["wizard_tower"]
	g.rooms["wizard_tower"].monsters = append(g.rooms["wizard_tower"].monsters, imp)
	
	// Dragon lair monsters
	dragon := NewMonster("ancient dragon", "A colossal red dragon with scales like molten gold and eyes like burning coals", 100, 15, true)
	dragon.location = g.rooms["dragon_lair"]
	g.rooms["dragon_lair"].monsters = append(g.rooms["dragon_lair"].monsters, dragon)
	
	// Cemetery monsters
	ghost := NewMonster("wandering spirit", "A translucent figure that wails mournfully as it drifts between the graves", 20, 3, false)
	ghost.location = g.rooms["cemetery"]
	g.rooms["cemetery"].monsters = append(g.rooms["cemetery"].monsters, ghost)
	
	wraith := NewMonster("vengeful wraith", "A dark specter filled with malice and hatred for the living", 35, 9, true)
	wraith.location = g.rooms["cemetery"]
	g.rooms["cemetery"].monsters = append(g.rooms["cemetery"].monsters, wraith)
	
	// Armory monsters
	guard := NewMonster("castle guard", "A heavily armored soldier sworn to protect the castle's treasures", 35, 7, false)
	guard.location = g.rooms["armory"]
	g.rooms["armory"].monsters = append(g.rooms["armory"].monsters, guard)
	
	// Pirate Cove monsters
	pirate := NewMonster("bloodthirsty pirate", "A scarred sailor with a wooden leg and a gleaming cutlass", 30, 8, true)
	pirate.location = g.rooms["pirate_cove"]
	g.rooms["pirate_cove"].monsters = append(g.rooms["pirate_cove"].monsters, pirate)
	
	kraken := NewMonster("sea kraken", "A massive tentacled beast that emerges from the depths to terrorize sailors", 80, 12, true)
	kraken.location = g.rooms["pirate_cove"]
	g.rooms["pirate_cove"].monsters = append(g.rooms["pirate_cove"].monsters, kraken)
	
	// Volcanic Cavern monsters
	salamander := NewMonster("lava salamander", "A lizard-like creature with scales that glow like embers", 40, 10, true)
	salamander.location = g.rooms["volcanic_cavern"]
	g.rooms["volcanic_cavern"].monsters = append(g.rooms["volcanic_cavern"].monsters, salamander)
	
	phoenix := NewMonster("flame phoenix", "A magnificent bird wreathed in eternal fire that rises from the ashes", 60, 14, false)
	phoenix.location = g.rooms["volcanic_cavern"]
	g.rooms["volcanic_cavern"].monsters = append(g.rooms["volcanic_cavern"].monsters, phoenix)
	
	// Ice Fortress monsters
	yeti := NewMonster("frost yeti", "A massive white-furred beast with icicles for claws", 50, 11, true)
	yeti.location = g.rooms["ice_fortress"]
	g.rooms["ice_fortress"].monsters = append(g.rooms["ice_fortress"].monsters, yeti)
	
	iceGolem := NewMonster("ice golem", "A towering construct made of solid ice and ancient magic", 70, 9, false)
	iceGolem.location = g.rooms["ice_fortress"]
	g.rooms["ice_fortress"].monsters = append(g.rooms["ice_fortress"].monsters, iceGolem)
	
	// Sky Temple monsters
	seraph := NewMonster("golden seraph", "A six-winged celestial being radiating divine light and power", 90, 16, false)
	seraph.location = g.rooms["sky_temple"]
	g.rooms["sky_temple"].monsters = append(g.rooms["sky_temple"].monsters, seraph)
	
	stormElemental := NewMonster("storm elemental", "A swirling vortex of wind and lightning with glowing eyes", 45, 13, true)
	stormElemental.location = g.rooms["sky_temple"]
	g.rooms["sky_temple"].monsters = append(g.rooms["sky_temple"].monsters, stormElemental)
	
	// Cursed Swamp monsters
	swampTroll := NewMonster("bog troll", "A massive troll covered in moss and slime, reeking of decay", 55, 9, true)
	swampTroll.location = g.rooms["cursed_swamp"]
	g.rooms["cursed_swamp"].monsters = append(g.rooms["cursed_swamp"].monsters, swampTroll)
	
	willOWisp := NewMonster("will-o'-wisp", "A dancing ball of eerie light that leads travelers astray", 20, 6, false)
	willOWisp.location = g.rooms["cursed_swamp"]
	g.rooms["cursed_swamp"].monsters = append(g.rooms["cursed_swamp"].monsters, willOWisp)
	
	// Crystal Mines monsters
	crystalSpider := NewMonster("crystal spider", "A spider with a crystalline carapace that refracts light into deadly beams", 35, 8, true)
	crystalSpider.location = g.rooms["crystal_mines"]
	g.rooms["crystal_mines"].monsters = append(g.rooms["crystal_mines"].monsters, crystalSpider)
	
	earthElemental := NewMonster("earth elemental", "A hulking creature of living stone and gems", 65, 10, false)
	earthElemental.location = g.rooms["crystal_mines"]
	g.rooms["crystal_mines"].monsters = append(g.rooms["crystal_mines"].monsters, earthElemental)
	
	// Haunted Library monsters
	librarian := NewMonster("spectral librarian", "The ghostly keeper of forbidden knowledge, eternally bound to the library", 40, 7, false)
	librarian.location = g.rooms["haunted_library"]
	g.rooms["haunted_library"].monsters = append(g.rooms["haunted_library"].monsters, librarian)
	
	bookWyrm := NewMonster("ancient book wyrm", "A serpentine dragon that devours knowledge and breathes ink", 75, 11, true)
	bookWyrm.location = g.rooms["haunted_library"]
	g.rooms["haunted_library"].monsters = append(g.rooms["haunted_library"].monsters, bookWyrm)
	
	// Goblin Warren monsters
	goblinKing := NewMonster("goblin king", "The cruel ruler of the goblin warren, adorned with stolen treasures", 45, 8, true)
	goblinKing.location = g.rooms["goblin_warren"]
	g.rooms["goblin_warren"].monsters = append(g.rooms["goblin_warren"].monsters, goblinKing)
	
	goblinShaman := NewMonster("goblin shaman", "A wicked spellcaster who communes with dark spirits", 30, 9, false)
	goblinShaman.location = g.rooms["goblin_warren"]
	g.rooms["goblin_warren"].monsters = append(g.rooms["goblin_warren"].monsters, goblinShaman)
	
	goblinWarrior := NewMonster("goblin warrior", "A fierce goblin fighter with crude weapons and a vicious temperament", 25, 6, true)
	goblinWarrior.location = g.rooms["goblin_warren"]
	g.rooms["goblin_warren"].monsters = append(g.rooms["goblin_warren"].monsters, goblinWarrior)
}

func (g *Game) gameLoop() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	
	for g.running {
		select {
		case <-ticker.C:
			g.processMonsterAI()
		}
	}
}

func (g *Game) processMonsterAI() {
	for _, room := range g.rooms {
		for _, monster := range room.monsters {
			if !monster.alive {
				continue
			}
			
			if monster.aggressive && len(room.players) > 0 {
				target := room.players[rand.Intn(len(room.players))]
				g.MonsterAttackPlayer(monster, target)
			} else if !monster.aggressive && len(room.players) > 0 {
				if rand.Float32() < 0.3 {
					target := room.players[rand.Intn(len(room.players))]
					g.MonsterAttackPlayer(monster, target)
				}
			}
		}
	}
}

func (g *Game) PlayerAttackMonster(player *Player, monsterName string) {
	var target *Monster
	for _, monster := range player.location.monsters {
		if monster.name == monsterName && monster.alive {
			target = monster
			break
		}
	}
	
	if target == nil {
		player.SendMessage(ColorError("There is no such monster here."))
		return
	}
	
	baseDamage := player.damage
	if player.weapon != nil {
		baseDamage += player.weapon.damage
	}
	damage := baseDamage + rand.Intn(3) - 1
	if damage < 1 {
		damage = 1
	}
	
	isDead := target.TakeDamage(damage)
	
	if isDead {
		player.SendMessage(fmt.Sprintf("%sYou kill the %s!%s", ColorSuccess(""), ColorMonster(target.name), ColorReset))
		player.location.Broadcast(fmt.Sprintf("%s kills the %s!", ColorName(player.name), ColorMonster(target.name)), player)
	} else {
		player.SendMessage(fmt.Sprintf("You attack the %s for %s%d damage%s!", ColorMonster(target.name), ColorDamage(""), damage, ColorReset))
		player.location.Broadcast(fmt.Sprintf("%s attacks the %s!", ColorName(player.name), ColorMonster(target.name)), player)
	}
}

func (g *Game) MonsterAttackPlayer(monster *Monster, player *Player) {
	if !monster.alive {
		return
	}
	
	baseDamage := monster.damage + rand.Intn(5) - 2
	if baseDamage < 1 {
		baseDamage = 1
	}
	
	defense := 0
	if player.armor != nil {
		defense = player.armor.defense
	}
	
	damage := baseDamage - defense
	if damage < 1 {
		damage = 1
	}
	
	isDead := player.TakeDamage(damage)
	
	if isDead {
		player.SendMessage(ColorDamage("You have been killed!"))
		player.location.Broadcast(fmt.Sprintf("%s has been killed by %s!", ColorName(player.name), ColorMonster(monster.name)), player)
		g.respawnPlayer(player)
	} else {
		player.SendMessage(fmt.Sprintf("The %s attacks you for %s%d damage%s!", ColorMonster(monster.name), ColorDamage(""), damage, ColorReset))
		player.location.Broadcast(fmt.Sprintf("%s attacks %s!", ColorMonster(monster.name), ColorName(player.name)), player)
	}
}

func (g *Game) respawnPlayer(player *Player) {
	player.health = player.maxHealth
	
	if player.location != nil {
		for i, p := range player.location.players {
			if p == player {
				player.location.players = append(player.location.players[:i], player.location.players[i+1:]...)
				break
			}
		}
	}
	
	townSquare := g.rooms["town_square"]
	player.location = townSquare
	townSquare.players = append(townSquare.players, player)
	
	player.SendMessage(ColorHealing("You respawn in the town square, fully healed."))
	player.location.Broadcast(fmt.Sprintf("%s respawns.", ColorName(player.name)), player)
}