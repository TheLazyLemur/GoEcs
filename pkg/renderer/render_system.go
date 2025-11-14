package renderer

import (
	"github.com/TheLazyLemur/SpaceImpact/pkg/assets"
	"github.com/TheLazyLemur/SpaceImpact/pkg/ecs"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// RenderSystem handles rendering of entities
type RenderSystem struct {
	assetManager    *assets.AssetManager
	transformSystem *ecs.TransformSystem
}

// NewRenderSystem creates a new render system
func NewRenderSystem(assetMgr *assets.AssetManager, transformSys *ecs.TransformSystem) *RenderSystem {
	return &RenderSystem{
		assetManager:    assetMgr,
		transformSystem: transformSys,
	}
}

// Render draws all entities with Transform and RenderMesh components
func (rs *RenderSystem) Render(scene *ecs.Scene) {
	world := scene.GetWorld()

	// Query all entities with both Transform and RenderMesh components
	q := query.NewQuery(filter.Contains(ecs.TransformComponent, ecs.RenderMeshComponent))

	q.Each(world.GetDonburiWorld(), func(entry *donburi.Entry) {
		rs.renderEntity(entry)
	})
}

// renderEntity renders a single entity
func (rs *RenderSystem) renderEntity(entity ecs.EntityID) {
	// Get transform
	transform := ecs.TransformComponent.Get(entity)

	// Get render mesh
	renderMesh := ecs.RenderMeshComponent.Get(entity)

	// Get mesh from registry
	mesh, ok := rs.assetManager.GetMeshes().Get(renderMesh.MeshID)
	if !ok {
		// Draw placeholder for missing mesh
		rs.drawMissingAsset(transform.Position)
		return
	}

	// Get material from registry
	material, ok := rs.assetManager.GetMaterials().Get(renderMesh.MaterialID)
	if !ok {
		// Use default material if not found
		material, _ = rs.assetManager.GetMaterials().Get("default")
	}

	// Get world matrix
	worldMatrix, ok := rs.transformSystem.GetWorldMatrix(entity)
	if !ok {
		worldMatrix = rl.MatrixIdentity()
	}

	// Draw mesh
	rl.DrawMesh(mesh, material, worldMatrix)
}

// drawMissingAsset draws a magenta wireframe cube as placeholder
func (rs *RenderSystem) drawMissingAsset(position rl.Vector3) {
	rl.DrawCubeWires(position, 1.0, 1.0, 1.0, rl.Magenta)
	rl.DrawCube(position, 0.1, 0.1, 0.1, rl.Magenta)
}
