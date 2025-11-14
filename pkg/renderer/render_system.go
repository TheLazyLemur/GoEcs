package renderer

import (
	"github.com/TheLazyLemur/SpaceImpact/pkg/assets"
	"github.com/TheLazyLemur/SpaceImpact/pkg/ecs"
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
	entities := world.GetEntitiesWithComponents(ecs.ComponentTypeTransform, ecs.ComponentTypeRenderMesh)

	for _, entityID := range entities {
		rs.renderEntity(world, entityID)
	}
}

// renderEntity renders a single entity
func (rs *RenderSystem) renderEntity(world *ecs.World, entityID ecs.EntityID) {
	// Get transform
	transformComp, ok := world.GetComponent(entityID, ecs.ComponentTypeTransform)
	if !ok {
		return
	}
	transform, ok := transformComp.(*ecs.Transform)
	if !ok {
		return
	}

	// Get render mesh
	meshComp, ok := world.GetComponent(entityID, ecs.ComponentTypeRenderMesh)
	if !ok {
		return
	}
	renderMesh, ok := meshComp.(*ecs.RenderMesh)
	if !ok {
		return
	}

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
	worldMatrix, ok := rs.transformSystem.GetWorldMatrix(entityID)
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
