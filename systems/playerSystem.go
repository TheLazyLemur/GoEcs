package systems

import (
	ent "github.com/TheLazyLemur/SpaceImpact/entities"
	cmp "github.com/TheLazyLemur/SpaceImpact/components"

	"github.com/EngoEngine/ecs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerSystem struct {
	Entities map[uint64]ent.PlayerEntity
}

type Playerable interface {
	ecs.BasicFace
	cmp.MovementFace
	cmp.HealthFace
	cmp.SpaceFace
}

func (m *PlayerSystem) Add(basic ecs.BasicEntity, a *cmp.MovementComponent, b *cmp.HealthComponent, c *cmp.SpaceComponent) {
	m.Entities[basic.ID()] = ent.PlayerEntity{
		BasicEntity:       basic,
		MovementComponent: a,
		HealthComponent:   b,
		SpaceComponent:    c,
	}
}

func (m *PlayerSystem) Remove(basic ecs.BasicEntity) {
	delete(m.Entities, basic.ID())
}

func (m *PlayerSystem) Update(dt float32) {
	for _, p := range m.Entities {

		if rl.IsKeyDown(rl.KeyRight) {
			p.SpaceComponent.ChangeX(1 * p.MovementComponent.Speed * dt)
		}

		if rl.IsKeyDown(rl.KeyLeft) {
			p.SpaceComponent.ChangeX(-1 * p.MovementComponent.Speed * dt)
		}

		if rl.IsKeyDown(rl.KeyUp) {
			p.SpaceComponent.ChangeY(-1 * p.MovementComponent.Speed * dt)
		}

		if rl.IsKeyDown(rl.KeyDown) {
			p.SpaceComponent.ChangeY(1 * p.MovementComponent.Speed * dt)
		}

		rec := rl.Rectangle{
			X:      p.SpaceComponent.GetX(),
			Y:      p.SpaceComponent.GetY(),
			Width:  p.SpaceComponent.GetWidth(),
			Height: p.SpaceComponent.GetHeight(),
		}

		rl.DrawRectangleRec(rec, rl.Green)
	}
}

func (m *PlayerSystem) AddByInterface(o ecs.Identifier) {
	obj := o.(Playerable)
	m.Add(*obj.GetBasicEntity(), obj.GetMovementComponent(), obj.GetHealthComponent(), obj.GetSpaceComponent())
}
