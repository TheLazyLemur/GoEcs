package editor

import (
	"fmt"

	"github.com/TheLazyLemur/SpaceImpact/pkg/ecs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// InspectorPanel represents the component inspector panel
type InspectorPanel struct {
	Bounds     rl.Rectangle
	ScrollView rl.Vector2
}

// NewInspectorPanel creates a new inspector panel
func NewInspectorPanel(x, y, width, height float32) *InspectorPanel {
	return &InspectorPanel{
		Bounds:     rl.NewRectangle(x, y, width, height),
		ScrollView: rl.NewVector2(0, 0),
	}
}

// Draw renders the inspector panel
func (ip *InspectorPanel) Draw(state *EditorState) {
	if !state.ShowInspector {
		return
	}

	// Draw panel background
	title := "Inspector"
	if state.HasSelection() {
		entityID := state.GetFirstSelected()
		name, ok := state.Scene.GetName(entityID)
		if ok {
			title = fmt.Sprintf("Inspector - %s", name)
		} else {
			title = fmt.Sprintf("Inspector - Entity #%d", entityID)
		}
	}

	// Draw panel background
	rl.DrawRectangleRec(ip.Bounds, rl.NewColor(40, 40, 40, 255))
	rl.DrawRectangleLinesEx(ip.Bounds, 1.0, rl.NewColor(80, 80, 80, 255))

	// Draw title bar
	titleBounds := rl.NewRectangle(ip.Bounds.X, ip.Bounds.Y, ip.Bounds.Width, 25)
	rl.DrawRectangleRec(titleBounds, rl.NewColor(60, 60, 60, 255))
	rl.DrawText(title, int32(ip.Bounds.X+10), int32(ip.Bounds.Y+5), 14, rl.RayWhite)

	// Content area
	contentBounds := rl.NewRectangle(
		ip.Bounds.X+10,
		ip.Bounds.Y+30,
		ip.Bounds.Width-20,
		ip.Bounds.Height-40,
	)

	// If nothing selected, show message
	if !state.HasSelection() {
		rl.DrawText(
			"No entity selected",
			int32(contentBounds.X+10),
			int32(contentBounds.Y+10),
			16,
			rl.Gray,
		)
		return
	}

	// Get first selected entity
	entityID := state.GetFirstSelected()

	rl.BeginScissorMode(
		int32(contentBounds.X),
		int32(contentBounds.Y),
		int32(contentBounds.Width),
		int32(contentBounds.Height),
	)

	yOffset := float32(0)

	// Draw Name component
	if state.Scene.HasComponent(entityID, ecs.ComponentTypeName) {
		yOffset += ip.drawNameComponent(state, entityID, contentBounds, yOffset)
	}

	// Draw Transform component
	if state.Scene.HasComponent(entityID, ecs.ComponentTypeTransform) {
		yOffset += ip.drawTransformComponent(state, entityID, contentBounds, yOffset)
	}

	// Draw RenderMesh component
	if state.Scene.HasComponent(entityID, ecs.ComponentTypeRenderMesh) {
		yOffset += ip.drawRenderMeshComponent(state, entityID, contentBounds, yOffset)
	}

	rl.EndScissorMode()
}

func (ip *InspectorPanel) drawNameComponent(state *EditorState, entityID ecs.EntityID, bounds rl.Rectangle, yOffset float32) float32 {
	y := bounds.Y + yOffset - ip.ScrollView.Y

	// Component header
	rl.DrawRectangle(
		int32(bounds.X),
		int32(y),
		int32(bounds.Width),
		25,
		rl.NewColor(70, 70, 70, 255),
	)
	rl.DrawText("Name", int32(bounds.X+5), int32(y+5), 16, rl.RayWhite)

	y += 30

	// Name field (read-only for now - text input is complex)
	name, _ := state.Scene.GetName(entityID)
	rl.DrawText("Name:", int32(bounds.X+5), int32(y), 14, rl.LightGray)
	rl.DrawText(name, int32(bounds.X+60), int32(y), 14, rl.RayWhite)

	return 70 // Total height used
}

func (ip *InspectorPanel) drawTransformComponent(state *EditorState, entityID ecs.EntityID, bounds rl.Rectangle, yOffset float32) float32 {
	y := bounds.Y + yOffset - ip.ScrollView.Y

	// Component header
	rl.DrawRectangle(
		int32(bounds.X),
		int32(y),
		int32(bounds.Width),
		25,
		rl.NewColor(70, 70, 70, 255),
	)
	rl.DrawText("Transform", int32(bounds.X+5), int32(y+5), 16, rl.RayWhite)

	y += 30

	transform, ok := state.Scene.GetTransform(entityID)
	if !ok {
		return 60
	}

	// Position
	rl.DrawText("Position:", int32(bounds.X+5), int32(y), 14, rl.LightGray)
	y += 20
	posText := fmt.Sprintf("X: %.2f  Y: %.2f  Z: %.2f", transform.Position.X, transform.Position.Y, transform.Position.Z)
	rl.DrawText(posText, int32(bounds.X+10), int32(y), 12, rl.RayWhite)

	y += 25

	// Scale
	rl.DrawText("Scale:", int32(bounds.X+5), int32(y), 14, rl.LightGray)
	y += 20
	scaleText := fmt.Sprintf("X: %.2f  Y: %.2f  Z: %.2f", transform.Scale.X, transform.Scale.Y, transform.Scale.Z)
	rl.DrawText(scaleText, int32(bounds.X+10), int32(y), 12, rl.RayWhite)

	y += 25

	// Rotation (quaternion - simplified display)
	rl.DrawText("Rotation:", int32(bounds.X+5), int32(y), 14, rl.LightGray)
	y += 20
	rotText := fmt.Sprintf("X: %.2f  Y: %.2f  Z: %.2f  W: %.2f",
		transform.Rotation.X, transform.Rotation.Y, transform.Rotation.Z, transform.Rotation.W)
	rl.DrawText(rotText, int32(bounds.X+10), int32(y), 12, rl.RayWhite)

	return 165 // Total height used
}

func (ip *InspectorPanel) drawRenderMeshComponent(state *EditorState, entityID ecs.EntityID, bounds rl.Rectangle, yOffset float32) float32 {
	y := bounds.Y + yOffset - ip.ScrollView.Y

	// Component header
	rl.DrawRectangle(
		int32(bounds.X),
		int32(y),
		int32(bounds.Width),
		25,
		rl.NewColor(70, 70, 70, 255),
	)
	rl.DrawText("Render Mesh", int32(bounds.X+5), int32(y+5), 16, rl.RayWhite)

	y += 30

	renderMesh, ok := state.Scene.GetRenderMesh(entityID)
	if !ok {
		return 60
	}

	// Mesh ID
	rl.DrawText("Mesh:", int32(bounds.X+5), int32(y), 14, rl.LightGray)
	rl.DrawText(renderMesh.MeshID, int32(bounds.X+70), int32(y), 14, rl.RayWhite)
	y += 25

	// Material ID
	rl.DrawText("Material:", int32(bounds.X+5), int32(y), 14, rl.LightGray)
	rl.DrawText(renderMesh.MaterialID, int32(bounds.X+70), int32(y), 14, rl.RayWhite)

	return 90 // Total height used
}
