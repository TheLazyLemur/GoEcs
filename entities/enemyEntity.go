package entities

import (
	cmp "github.com/TheLazyLemur/SpaceImpact/components"

	"github.com/EngoEngine/ecs"
)

type EnemyEntity struct {
	ecs.BasicEntity
	*cmp.AiComponent
	*cmp.HealthComponent
}
