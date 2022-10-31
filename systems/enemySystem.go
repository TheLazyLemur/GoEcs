package systems

import (
	"fmt"

	"GoEcs/components"
	"GoEcs/entities"
	"github.com/EngoEngine/ecs"
)

type EnemySystem struct {
	Entities map[uint64]entities.EnemyEntity
}

type Enemyable interface {
	ecs.BasicFace
	components.MovementFace
	components.HealthFace
}

func (m *EnemySystem) Add(basic ecs.BasicEntity, a *components.MovementComponent, b *components.HealthComponent) {
	m.Entities[basic.ID()] = entities.EnemyEntity{BasicEntity: basic, MovementComponent: a, HealthComponent: b}
}
func (m *EnemySystem) Remove(basic ecs.BasicEntity) {
	delete(m.Entities, basic.ID())
}

func (m *EnemySystem) Update(dt float32) {
	for _, e := range m.Entities {

		e.HealthComponent.Current--
		if e.HealthComponent.Current <= 0 {
			fmt.Println("Removing entity", e.BasicEntity.ID())
			m.Remove(e.BasicEntity)
			continue
		}

		fmt.Println(e.HealthComponent.Current, "/", e.HealthComponent.Max)
	}
}

func (m *EnemySystem) AddByInterface(o ecs.Identifier) {
	obj := o.(Enemyable)
	m.Add(*obj.GetBasicEntity(), obj.GetMovementComponent(), obj.GetHealthComponent())
}


