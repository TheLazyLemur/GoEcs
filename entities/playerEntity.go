package entities

import (
	cmp "github.com/TheLazyLemur/SpaceImpact/components"

	"github.com/EngoEngine/ecs"
)

type PlayerEntity struct {
	ecs.BasicEntity
	*cmp.MovementComponent
	*cmp.HealthComponent
	*cmp.SpaceComponent
}
