package components

type HealthComponent struct {
	Max     int
	Current int
}

func (a *HealthComponent) GetHealthComponent() *HealthComponent {
	return a
}

type HealthFace interface {
	GetHealthComponent() *HealthComponent
}
