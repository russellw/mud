package main

import (
	"fmt"
	"math/rand"
)

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

func (m *Monster) Attack(target *Player) {
	if !m.alive {
		return
	}
	
	damage := m.damage + rand.Intn(5) - 2
	if damage < 1 {
		damage = 1
	}
	
	isDead := target.TakeDamage(damage)
	
	if m.location != nil {
		if isDead {
			m.location.Broadcast(fmt.Sprintf("%s kills %s!", m.name, target.name), nil)
		} else {
			m.location.Broadcast(fmt.Sprintf("%s attacks %s for %d damage!", m.name, target.name, damage), nil)
		}
	}
}

func (m *Monster) Respawn() {
	m.health = m.maxHealth
	m.alive = true
}

func (m *Monster) GetStatus() string {
	if !m.alive {
		return m.name + " (dead)"
	}
	
	healthPercent := float64(m.health) / float64(m.maxHealth)
	
	if healthPercent > 0.8 {
		return m.name + " (healthy)"
	} else if healthPercent > 0.5 {
		return m.name + " (wounded)"
	} else if healthPercent > 0.2 {
		return m.name + " (badly wounded)"
	} else {
		return m.name + " (near death)"
	}
}