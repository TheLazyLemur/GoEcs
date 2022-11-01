package components

type HealthComponent struct {
	Max     int
	Current int
}

func (cmp *HealthComponent) GetHealthComponent() *HealthComponent {
	return cmp
}

type HealthFace interface {
	GetHealthComponent() *HealthComponent
}
