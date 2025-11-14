package editor

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// HierarchyPanel represents the entity hierarchy panel
type HierarchyPanel struct {
	Bounds     rl.Rectangle
	ScrollView rl.Vector2
}

// NewHierarchyPanel creates a new hierarchy panel
func NewHierarchyPanel(x, y, width, height float32) *HierarchyPanel {
	return &HierarchyPanel{
		Bounds:     rl.NewRectangle(x, y, width, height),
		ScrollView: rl.NewVector2(0, 0),
	}
}

// Draw renders the hierarchy panel
func (hp *HierarchyPanel) Draw(state *EditorState) {
	if !state.ShowHierarchy {
		return
	}

	// Draw panel background
	rl.DrawRectangleRec(hp.Bounds, rl.NewColor(40, 40, 40, 255))
	rl.DrawRectangleLinesEx(hp.Bounds, 1.0, rl.NewColor(80, 80, 80, 255))

	// Draw title bar
	titleBounds := rl.NewRectangle(hp.Bounds.X, hp.Bounds.Y, hp.Bounds.Width, 25)
	rl.DrawRectangleRec(titleBounds, rl.NewColor(60, 60, 60, 255))
	rl.DrawText("Hierarchy", int32(hp.Bounds.X+10), int32(hp.Bounds.Y+5), 16, rl.RayWhite)

	// Content area (inside panel)
	contentBounds := rl.NewRectangle(
		hp.Bounds.X+10,
		hp.Bounds.Y+30,
		hp.Bounds.Width-20,
		hp.Bounds.Height-40,
	)

	// Get all entities
	entities := state.Scene.GetAllEntities()

	// Draw each entity
	yOffset := float32(0)
	itemHeight := float32(25)

	rl.BeginScissorMode(
		int32(contentBounds.X),
		int32(contentBounds.Y),
		int32(contentBounds.Width),
		int32(contentBounds.Height),
	)

	for i, entityID := range entities {
		// Get entity name
		name, ok := state.Scene.GetName(entityID)
		if !ok {
			name = fmt.Sprintf("Entity %d", entityID)
		}

		itemY := contentBounds.Y + yOffset - hp.ScrollView.Y
		itemBounds := rl.NewRectangle(
			contentBounds.X,
			itemY,
			contentBounds.Width,
			itemHeight,
		)

		// Check if visible
		if itemY+itemHeight < contentBounds.Y || itemY > contentBounds.Y+contentBounds.Height {
			yOffset += itemHeight
			continue
		}

		// Draw item background
		isSelected := state.IsSelected(entityID)
		bgColor := rl.NewColor(30, 30, 30, 255)
		if isSelected {
			bgColor = rl.NewColor(100, 150, 255, 200)
		}

		// Hover effect
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), itemBounds) {
			if !isSelected {
				bgColor = rl.NewColor(50, 50, 50, 255)
			}

			// Handle click
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
					state.AddToSelection(entityID)
				} else {
					state.SelectEntity(entityID)
				}
			}
		}

		rl.DrawRectangleRec(itemBounds, bgColor)

		// Draw entity name
		textColor := rl.RayWhite
		if isSelected {
			textColor = rl.White
		}
		rl.DrawText(
			name,
			int32(itemBounds.X+5),
			int32(itemBounds.Y+5),
			16,
			textColor,
		)

		// Draw entity ID (small, gray)
		idText := fmt.Sprintf("#%d", entityID)
		rl.DrawText(
			idText,
			int32(itemBounds.X+itemBounds.Width-60),
			int32(itemBounds.Y+7),
			12,
			rl.Gray,
		)

		yOffset += itemHeight

		// Draw separator
		if i < len(entities)-1 {
			rl.DrawLine(
				int32(itemBounds.X),
				int32(itemBounds.Y+itemBounds.Height),
				int32(itemBounds.X+itemBounds.Width),
				int32(itemBounds.Y+itemBounds.Height),
				rl.NewColor(60, 60, 60, 255),
			)
		}
	}

	rl.EndScissorMode()

	// Handle mouse wheel scrolling over panel
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), contentBounds) {
		wheel := rl.GetMouseWheelMove()
		hp.ScrollView.Y -= wheel * 20

		// Clamp scroll
		maxScroll := yOffset - contentBounds.Height
		if maxScroll < 0 {
			maxScroll = 0
		}
		if hp.ScrollView.Y < 0 {
			hp.ScrollView.Y = 0
		}
		if hp.ScrollView.Y > maxScroll {
			hp.ScrollView.Y = maxScroll
		}
	}
}
