package main

import (
	"GoEcs/components"
	"GoEcs/entities"
	"GoEcs/systems"

	"github.com/EngoEngine/ecs"
)

func main() {
	w := ecs.World{}

	var myable *systems.Enemyable

	w.AddSystemInterface(&systems.EnemySystem{Entities: make(map[uint64]entities.EnemyEntity)}, myable, nil)

	entity := entities.EnemyEntity{
		BasicEntity: ecs.NewBasic(),
		MovementComponent: &components.MovementComponent{
			Speed: 100,
		},
		HealthComponent: &components.HealthComponent{
			Max:     100,
			Current: 100,
		},
	}

	entity2 := entities.EnemyEntity{
		BasicEntity: ecs.NewBasic(),
		MovementComponent: &components.MovementComponent{
			Speed: 100,
		},
		HealthComponent: &components.HealthComponent{
			Max:     10,
			Current: 10,
		},
	}

	w.SortSystems()

	w.AddEntity(&entity)
	w.AddEntity(&entity2)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
	w.Update(0.1)
}
