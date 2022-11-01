package components

type HealthComponent struct {
	max     int
	current int
}

func (cmp *HealthComponent) GetHealthComponent() *HealthComponent {
	return cmp
}

type HealthFace interface {
	GetHealthComponent() *HealthComponent
}

func (cmp *HealthComponent) GetCurrentHealth() int {
	return cmp.current
}

func (cmp *HealthComponent) GetMaxHealth() int {
	return cmp.max
}

func (cmp *HealthComponent) ChangeHealth(val int){
    cmp.current += val
}

func (cmp *HealthComponent) IsDead() bool{
    return cmp.current <= 0
}

func NewHealthComponent(max int, current int) *HealthComponent {
    return &HealthComponent {
        max,
        current,
    }
}
