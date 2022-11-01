package components

type MovementComponent struct {
	Speed float32
}

func (cmp *MovementComponent) GetMovementComponent() *MovementComponent {
	return cmp
}

type MovementFace interface {
	GetMovementComponent() *MovementComponent
}
