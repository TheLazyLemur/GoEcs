package entities

import (
	cmp "GoEcs/components"

	"github.com/EngoEngine/ecs"
)

type PlayerEntity struct {
	ecs.BasicEntity
	*cmp.MovementComponent
    *cmp.HealthComponent
	*cmp.SpaceComponent
}
