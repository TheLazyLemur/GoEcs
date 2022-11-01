package entities

import (
	cmp "GoEcs/components"

	"github.com/EngoEngine/ecs"
)

type EnemyEntity struct {
	ecs.BasicEntity
	*cmp.AiComponent
	*cmp.HealthComponent
}
