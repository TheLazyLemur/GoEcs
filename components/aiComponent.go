package components

type AiComponent struct {
	Speed int
}

func (cmp *AiComponent) GetAiComponent() *AiComponent {
	return cmp
}

type AiFace interface {
	GetAiComponent() *AiComponent
}
