package entities

import (
	"GoEcs/components"

	"github.com/EngoEngine/ecs"
)

type EnemyEntity struct {
	ecs.BasicEntity
	*components.MovementComponent
	*components.HealthComponent
}
