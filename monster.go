package main

type Monster struct {
	name        string
	description string
	health      int
	maxHealth   int
	damage      int
	location    *Room
	aggressive  bool
	alive       bool
}

func NewMonster(name, description string, health, damage int, aggressive bool) *Monster {
	return &Monster{
		name:        name,
		description: description,
		health:      health,
		maxHealth:   health,
		damage:      damage,
		aggressive:  aggressive,
		alive:       true,
	}
}

func (m *Monster) TakeDamage(damage int) bool {
	m.health -= damage
	if m.health <= 0 {
		m.health = 0
		m.alive = false
		return true
	}
	return false
}


func (m *Monster) Respawn() {
	m.health = m.maxHealth
	m.alive = true
}

func (m *Monster) GetStatus() string {
	if !m.alive {
		return m.name + ColorBrightBlack + " (dead)" + ColorReset
	}
	
	healthPercent := float64(m.health) / float64(m.maxHealth)
	
	if healthPercent > 0.8 {
		return m.name + ColorGreen + " (healthy)" + ColorReset
	} else if healthPercent > 0.5 {
		return m.name + ColorYellow + " (wounded)" + ColorReset
	} else if healthPercent > 0.2 {
		return m.name + ColorRed + " (badly wounded)" + ColorReset
	} else {
		return m.name + ColorBrightRed + " (near death)" + ColorReset
	}
}