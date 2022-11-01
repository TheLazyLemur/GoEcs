package systems

import (
	"fmt"

	cmp "GoEcs/components"
	ent "GoEcs/entities"

	"github.com/EngoEngine/ecs"
	"github.com/gen2brain/raylib-go/raylib"
)

type EnemySystem struct {
	Entities map[uint64]ent.EnemyEntity
}

type Enemyable interface {
	ecs.BasicFace
	cmp.AiFace
	cmp.HealthFace
}

func (m *EnemySystem) Add(basic ecs.BasicEntity, a *cmp.AiComponent, b *cmp.HealthComponent) {
	m.Entities[basic.ID()] = ent.EnemyEntity{
		BasicEntity:       basic,
		AiComponent: a,
		HealthComponent:   b,
	}
}

func (m *EnemySystem) Remove(basic ecs.BasicEntity) {
	delete(m.Entities, basic.ID())
}

func (m *EnemySystem) Update(dt float32) {
	for i, e := range m.Entities {
		e.HealthComponent.Current--
		if e.HealthComponent.Current <= 0 {
			m.Remove(e.BasicEntity)
			continue
		}
		rl.DrawText(fmt.Sprintf("%d/%d", e.HealthComponent.Current, e.HealthComponent.Max), 0, int32(0+i*20), 20, rl.LightGray)
	}
}

func (m *EnemySystem) AddByInterface(o ecs.Identifier) {
	obj := o.(Enemyable)
	m.Add(*obj.GetBasicEntity(), obj.GetAiComponent(), obj.GetHealthComponent())
}
