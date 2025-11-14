package main

import (
	"flag"
	"fmt"

	"github.com/TheLazyLemur/SpaceImpact/pkg/assets"
	"github.com/TheLazyLemur/SpaceImpact/pkg/ecs"
	"github.com/TheLazyLemur/SpaceImpact/pkg/editor"
	"github.com/TheLazyLemur/SpaceImpact/pkg/renderer"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	editorMode = flag.Bool("editor", false, "Launch in editor mode")
)

func main() {
	flag.Parse()

	if *editorMode {
		fmt.Println("Launching in EDITOR mode...")
		runEditor()
	} else {
		fmt.Println("Launching in GAME mode...")
		runGame()
	}
}

func runEditor() {
	// Initialize window
	screenWidth := int32(1920)
	screenHeight := int32(1080)
	rl.InitWindow(screenWidth, screenHeight, "GoECS Editor")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	// Create asset manager and load default assets
	assetMgr := assets.NewAssetManager()
	assetMgr.CreateDefaultAssets()

	// Create scene
	scene := ecs.NewScene("TestScene")

	// Create editor state
	editorState := editor.NewEditorState(scene)

	// Create UI panels
	hierarchyPanel := editor.NewHierarchyPanel(0, 0, editorState.HierarchyWidth, float32(screenHeight))
	inspectorPanel := editor.NewInspectorPanel(
		float32(screenWidth)-editorState.InspectorWidth,
		0,
		editorState.InspectorWidth,
		float32(screenHeight),
	)

	// Create viewport (center area between panels)
	viewportX := editorState.HierarchyWidth
	viewportY := float32(0)
	viewportWidth := float32(screenWidth) - editorState.HierarchyWidth - editorState.InspectorWidth
	viewportHeight := float32(screenHeight)

	viewport := editor.NewViewport(viewportX, viewportY, viewportWidth, viewportHeight)
	defer viewport.Cleanup()

	// Add viewport to state
	editorState.Viewports = append(editorState.Viewports, viewport)

	// Create systems
	transformSystem := ecs.NewTransformSystem()
	renderSystem := renderer.NewRenderSystem(assetMgr, transformSystem)

	// Create gizmo
	gizmo := editor.NewGizmo()

	// Create some test entities
	createTestEntities(scene)

	fmt.Println("Scene created with", len(scene.GetAllEntities()), "entities")

	// Main loop
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		// Handle mode switching
		if rl.IsKeyPressed(rl.KeyF5) {
			if editorState.Mode == editor.ModeEdit {
				editorState.SetMode(editor.ModePlay)
				fmt.Println("Entering PLAY mode")
			} else {
				editorState.SetMode(editor.ModeEdit)
				fmt.Println("Entering EDIT mode")
			}
		}

		// Update gizmo (only in edit mode)
		if editorState.Mode == editor.ModeEdit {
			gizmo.Update(viewport, editorState, transformSystem)
		}

		// Handle viewport picking (only in edit mode and not dragging gizmo)
		if editorState.Mode == editor.ModeEdit && !gizmo.IsDragging {
			// Check for left mouse click in viewport
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				mousePos := rl.GetMousePosition()

				// Make sure click is in viewport (not on panels)
				if viewport.IsMouseOver() {
					// Perform picking
					pickResult := editor.PickEntityInViewport(viewport, mousePos, scene, transformSystem)

					if pickResult.Hit {
						// Handle multi-select with Ctrl
						if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
							editorState.AddToSelection(pickResult.EntityID)
						} else {
							editorState.SelectEntity(pickResult.EntityID)
						}

						name, _ := scene.GetName(pickResult.EntityID)
						fmt.Printf("Selected: %s (ID: %d, Distance: %.2f)\n", name, pickResult.EntityID, pickResult.Distance)
					} else {
						// Clicked in viewport but hit nothing - clear selection
						if !rl.IsKeyDown(rl.KeyLeftControl) && !rl.IsKeyDown(rl.KeyRightControl) {
							editorState.ClearSelection()
						}
					}
				}
			}
		}

		// Update editor camera (only in edit mode)
		if editorState.Mode == editor.ModeEdit {
			viewport.Camera.Update(dt)
		}

		// Update systems
		transformSystem.Update(scene)

		// === RENDER TO VIEWPORT TEXTURE ===
		viewport.Begin()

		camera := viewport.Camera.GetRaylibCamera()
		rl.BeginMode3D(camera)

		// Draw grid
		rl.DrawGrid(20, 1.0)

		// Render all entities
		renderSystem.Render(scene)

		// Draw selection highlights
		drawSelectionHighlights(editorState, transformSystem)

		// Draw gizmo for selected entity (only in edit mode)
		if editorState.Mode == editor.ModeEdit && editorState.HasSelection() {
			entityID := editorState.GetFirstSelected()
			transform, ok := editorState.Scene.GetTransform(entityID)
			if ok {
				gizmo.DrawTranslationGizmo(transform.Position, viewport.Camera)
			}
		}

		rl.EndMode3D()

		viewport.End()

		// === RENDER TO SCREEN ===
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(30, 30, 30, 255))

		// Draw viewport texture to screen
		viewport.Draw()

		// Draw UI panels
		hierarchyPanel.Draw(editorState)
		inspectorPanel.Draw(editorState)

		// Draw top bar with mode indicator
		modeColor := rl.NewColor(80, 250, 123, 255)
		if editorState.Mode == editor.ModePlay {
			modeColor = rl.NewColor(255, 184, 108, 255)
		}
		rl.DrawText(
			fmt.Sprintf("Mode: %s | FPS: %d | Entities: %d",
				editorState.Mode.String(),
				rl.GetFPS(),
				len(scene.GetAllEntities()),
			),
			10, 10, 16, modeColor,
		)

		// Draw help text
		rl.DrawText("Click to select | Drag gizmo arrows to move | F5: Play/Edit | ESC: Exit", 10, int32(screenHeight)-25, 14, rl.LightGray)

		rl.EndDrawing()
	}
}

func createTestEntities(scene *ecs.Scene) {
	identityQuat := rl.NewQuaternion(0, 0, 0, 1) // Identity quaternion

	// Create a red cube
	cube1 := scene.CreateEntity("Red Cube")
	scene.AddTransform(cube1)
	scene.AddRenderMesh(cube1, "cube", "red")
	scene.SetTransform(cube1, rl.NewVector3(-3, 0, 0), rl.NewVector3(2, 2, 2), identityQuat)

	// Create a green sphere
	sphere1 := scene.CreateEntity("Green Sphere")
	scene.AddTransform(sphere1)
	scene.AddRenderMesh(sphere1, "sphere", "green")
	scene.SetTransform(sphere1, rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2), identityQuat)

	// Create a blue cube
	cube2 := scene.CreateEntity("Blue Cube")
	scene.AddTransform(cube2)
	scene.AddRenderMesh(cube2, "cube", "blue")
	scene.SetTransform(cube2, rl.NewVector3(3, 0, 0), rl.NewVector3(2, 2, 2), identityQuat)

	// Create a floor plane
	floor := scene.CreateEntity("Floor")
	scene.AddTransform(floor)
	scene.AddRenderMesh(floor, "plane", "default")
	scene.SetTransform(floor, rl.NewVector3(0, -1, 0), rl.NewVector3(1, 1, 1), identityQuat)
}

func drawSelectionHighlights(editorState *editor.EditorState, transformSystem *ecs.TransformSystem) {
	if !editorState.HasSelection() {
		return
	}

	// Draw wireframe highlights for selected entities
	for _, entityID := range editorState.Selected {
		transform, ok := editorState.Scene.GetTransform(entityID)
		if !ok {
			continue
		}

		// Get render mesh to determine shape
		renderMesh, ok := editorState.Scene.GetRenderMesh(entityID)
		if !ok {
			continue
		}

		// Draw highlight based on mesh type
		highlightColor := rl.NewColor(255, 184, 108, 255) // Orange highlight

		switch renderMesh.MeshID {
		case "cube":
			// Draw wireframe cube at entity position with scale
			size := rl.NewVector3(2.0, 2.0, 2.0)
			size.X *= transform.Scale.X
			size.Y *= transform.Scale.Y
			size.Z *= transform.Scale.Z
			rl.DrawCubeWiresV(transform.Position, size, highlightColor)

		case "sphere":
			// Draw wireframe sphere
			radius := transform.Scale.X
			rl.DrawSphereWires(transform.Position, radius, 16, 16, highlightColor)

		case "plane":
			// Draw wireframe for plane (as a thin box)
			size := rl.NewVector3(10, 0.1, 10)
			size.X *= transform.Scale.X
			size.Y *= transform.Scale.Y
			size.Z *= transform.Scale.Z
			rl.DrawCubeWiresV(transform.Position, size, highlightColor)
		}
	}
}

func runGame() {
	// Initialize window
	screenWidth := int32(800)
	screenHeight := int32(600)
	rl.InitWindow(screenWidth, screenHeight, "Space Impact")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	// Main loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawText("GAME MODE", 10, 10, 20, rl.Green)
		rl.DrawText("(Not implemented yet)", 10, 40, 20, rl.Gray)
		rl.DrawText("Press ESC to exit", 10, 70, 20, rl.LightGray)

		rl.EndDrawing()
	}
}
