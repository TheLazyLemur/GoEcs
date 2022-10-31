package components

type MovementComponent struct {
	Speed int
}

func (a *MovementComponent) GetMovementComponent() *MovementComponent {
	return a
}

type MovementFace interface {
	GetMovementComponent() *MovementComponent
}
