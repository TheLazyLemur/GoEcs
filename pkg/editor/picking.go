package editor

import (
	"github.com/TheLazyLemur/SpaceImpact/pkg/ecs"
	pkgmath "github.com/TheLazyLemur/SpaceImpact/pkg/math"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// PickResult represents the result of a picking operation
type PickResult struct {
	Hit      bool
	EntityID ecs.EntityID
	Distance float32
}

// PickEntityInViewport performs raycasting to find the entity under the mouse
func PickEntityInViewport(
	viewport *Viewport,
	mousePos rl.Vector2,
	scene *ecs.Scene,
	transformSystem *ecs.TransformSystem,
) PickResult {
	// Convert mouse position to viewport-relative coordinates
	viewportMouseX := mousePos.X - viewport.Bounds.X
	viewportMouseY := mousePos.Y - viewport.Bounds.Y

	// Check if mouse is within viewport bounds
	if viewportMouseX < 0 || viewportMouseX > viewport.Bounds.Width ||
		viewportMouseY < 0 || viewportMouseY > viewport.Bounds.Height {
		return PickResult{Hit: false}
	}

	// Create viewport-relative mouse position
	viewportMouse := rl.NewVector2(viewportMouseX, viewportMouseY)

	// Generate ray from camera through mouse position
	camera := viewport.Camera.GetRaylibCamera()
	viewportRect := rl.NewRectangle(0, 0, viewport.Bounds.Width, viewport.Bounds.Height)
	ray := pkgmath.GetMouseRay(viewportMouse, camera, viewportRect)

	// Find closest hit
	closestHit := PickResult{Hit: false, Distance: 999999.0}

	// Get all entities with Transform and RenderMesh using Donburi query
	world := scene.GetWorld()
	q := query.NewQuery(filter.Contains(ecs.TransformComponent, ecs.RenderMeshComponent))

	q.Each(world.GetDonburiWorld(), func(entry *donburi.Entry) {
		// Get render mesh to determine bounding box
		renderMesh := ecs.RenderMeshComponent.Get(entry)

		// Get local bounding box for mesh type
		localBounds := pkgmath.GetDefaultMeshBounds(renderMesh.MeshID)

		// Get world transform matrix
		worldMatrix, ok := transformSystem.GetWorldMatrix(entry)
		if !ok {
			return
		}

		// Transform bounding box to world space
		worldBounds := pkgmath.TransformBoundingBox(localBounds, worldMatrix)

		// Test ray intersection
		hit, distance := pkgmath.RayIntersectsBoundingBox(ray, worldBounds)
		if hit {
			// Only update if this is closer AND distance is positive (in front of camera)
			if distance > 0 && distance < closestHit.Distance {
				closestHit = PickResult{
					Hit:      true,
					EntityID: entry,
					Distance: distance,
				}
			}
		}
	})

	return closestHit
}
