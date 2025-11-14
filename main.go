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

	// Create viewport (full window for now)
	viewport := editor.NewViewport(0, 0, float32(screenWidth), float32(screenHeight))
	defer viewport.Cleanup()

	// Create asset manager and load default assets
	assetMgr := assets.NewAssetManager()
	assetMgr.CreateDefaultAssets()

	// Create scene
	scene := ecs.NewScene("TestScene")

	// Create systems
	transformSystem := ecs.NewTransformSystem()
	renderSystem := renderer.NewRenderSystem(assetMgr, transformSystem)

	// Create some test entities
	createTestEntities(scene)

	fmt.Println("Scene created with", len(scene.GetAllEntities()), "entities")

	// Main loop
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		// Update editor camera
		viewport.Camera.Update(dt)

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

		rl.EndMode3D()

		viewport.End()

		// === RENDER TO SCREEN ===
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw viewport texture to screen
		viewport.Draw()

		// Draw UI overlay
		rl.DrawText("EDITOR MODE - Phase 3 Complete", 10, 10, 20, rl.NewColor(80, 250, 123, 255))
		rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 10, 40, 20, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("Entities: %d", len(scene.GetAllEntities())), 10, 70, 20, rl.RayWhite)
		rl.DrawText("WASD+QE: Move camera | Right Mouse: Look | Shift: Speed boost", 10, 100, 20, rl.LightGray)
		rl.DrawText("Press ESC to exit", 10, 130, 20, rl.LightGray)

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
