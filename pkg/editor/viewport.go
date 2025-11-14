package editor

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Viewport represents a 3D rendering viewport
type Viewport struct {
	Bounds        rl.Rectangle
	RenderTexture rl.RenderTexture2D
	Camera        *EditorCamera
}

// NewViewport creates a new viewport with the given bounds
func NewViewport(x, y, width, height float32) *Viewport {
	renderTexture := rl.LoadRenderTexture(int32(width), int32(height))

	return &Viewport{
		Bounds:        rl.NewRectangle(x, y, width, height),
		RenderTexture: renderTexture,
		Camera:        NewEditorCamera(),
	}
}

// Resize updates the viewport bounds and recreates the render texture
func (v *Viewport) Resize(x, y, width, height float32) {
	if v.Bounds.Width == width && v.Bounds.Height == height {
		return // No need to resize
	}

	// Unload old render texture
	rl.UnloadRenderTexture(v.RenderTexture)

	// Create new render texture
	v.Bounds = rl.NewRectangle(x, y, width, height)
	v.RenderTexture = rl.LoadRenderTexture(int32(width), int32(height))
}

// Begin starts rendering to this viewport's texture
func (v *Viewport) Begin() {
	rl.BeginTextureMode(v.RenderTexture)
	rl.ClearBackground(rl.NewColor(40, 42, 54, 255))
}

// End stops rendering to this viewport's texture
func (v *Viewport) End() {
	rl.EndTextureMode()
}

// Draw renders the viewport's texture to the screen
func (v *Viewport) Draw() {
	// Source rectangle (flip Y because render texture is upside down)
	source := rl.NewRectangle(
		0, 0,
		float32(v.RenderTexture.Texture.Width),
		float32(-v.RenderTexture.Texture.Height), // Negative to flip
	)

	// Destination rectangle
	dest := v.Bounds

	// Draw the render texture to screen
	rl.DrawTexturePro(
		v.RenderTexture.Texture,
		source,
		dest,
		rl.NewVector2(0, 0),
		0.0,
		rl.White,
	)
}

// Cleanup unloads the render texture
func (v *Viewport) Cleanup() {
	rl.UnloadRenderTexture(v.RenderTexture)
}

// IsMouseOver checks if the mouse is over this viewport
func (v *Viewport) IsMouseOver() bool {
	mousePos := rl.GetMousePosition()
	return rl.CheckCollisionPointRec(mousePos, v.Bounds)
}
