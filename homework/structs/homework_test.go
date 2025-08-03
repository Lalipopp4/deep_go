package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i, s := range name {
			person.name[i] = byte(s)
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.otherMeta += uint32(mana << 5)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.otherMeta += uint32(health << 15)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rsel += uint16(respect << 12)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rsel += uint16(strength << 8)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rsel += uint16(experience << 4)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rsel += uint16(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.otherMeta = person.otherMeta | 0b10000
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.otherMeta = person.otherMeta | 0b1000
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.otherMeta = person.otherMeta | 0b100
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.otherMeta += uint32(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType

	Len2  = 3
	Len4  = 15
	Len10 = 1023
)

type GamePerson struct {
	x, y, z int32
	name    [42]byte
	// bin mask for respect, strength, experience, level
	rsel uint16
	// misc 		0000000
	// mana     	0000000000
	// health	 	0000000000
	// house    	0
	// gun   		0
	// family   	0
	// type     	00
	otherMeta uint32
	// gold
	gold uint32
}

func NewGamePerson(options ...Option) GamePerson {
	gp := GamePerson{}
	for _, op := range options {
		op(&gp)
	}

	return gp
}

func (p *GamePerson) Name() string {
	return unsafe.String(&p.name[0], 42)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.otherMeta&(Len10<<5)) >> 5
}

func (p *GamePerson) Health() int {
	return int(p.otherMeta&(Len10<<5)) >> 5
}

func (p *GamePerson) Respect() int {
	return int(p.rsel&(Len4<<12)) >> 12
}

func (p *GamePerson) Strength() int {
	return int(p.rsel&(Len4<<8)) >> 8
}

func (p *GamePerson) Experience() int {
	return int(p.rsel&(Len4<<4)) >> 4
}

func (p *GamePerson) Level() int {
	return int(p.rsel & Len4)
}

func (p *GamePerson) HasHouse() bool {
	return p.otherMeta&0b10000 != 0
}

func (p *GamePerson) HasGun() bool {
	return p.otherMeta&0b1000 != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.otherMeta&0b100 != 0
}

func (p *GamePerson) Type() int {
	return int(p.otherMeta & Len2)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
