package components

type SpaceComponent struct {
	width  int
	height int
	x      float32
	y      float32
}

func (cmp *SpaceComponent) GetSpaceComponent() *SpaceComponent {
	return cmp
}

type SpaceFace interface {
	GetSpaceComponent() *SpaceComponent
}

func (c *SpaceComponent) ChangeX(newX float32) {
	c.x += newX
}

func (c *SpaceComponent) ChangeY(newY float32) {
	c.y += newY
}

func (c *SpaceComponent) GetX() float32 {
	return c.x
}

func (c *SpaceComponent) GetY() float32 {
	return c.y
}

func (c *SpaceComponent) GetHeight() float32 {
	return float32(c.height)
}

func (c *SpaceComponent) GetWidth() float32 {
	return float32(c.width)
}

func NewSpaceComponent() *SpaceComponent {
	return &SpaceComponent{
		height: 10,
		width:  10,
	}
}
