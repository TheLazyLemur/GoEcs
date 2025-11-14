package ecs

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TransformCache stores computed world matrices for entities
type TransformCache struct {
	matrices map[EntityID]rl.Matrix
}

// NewTransformCache creates a new transform cache
func NewTransformCache() *TransformCache {
	return &TransformCache{
		matrices: make(map[EntityID]rl.Matrix),
	}
}

// TransformSystem processes all Transform components and builds world matrices
type TransformSystem struct {
	cache *TransformCache
}

// NewTransformSystem creates a new transform system
func NewTransformSystem() *TransformSystem {
	return &TransformSystem{
		cache: NewTransformCache(),
	}
}

// Update processes all entities with Transform components
func (ts *TransformSystem) Update(scene *Scene) {
	world := scene.GetWorld()
	entities := world.GetEntitiesWithComponent(ComponentTypeTransform)

	for _, entityID := range entities {
		comp, ok := world.GetComponent(entityID, ComponentTypeTransform)
		if !ok {
			continue
		}

		transform, ok := comp.(*Transform)
		if !ok {
			continue
		}

		// Build world matrix from position, rotation, and scale
		matrix := ts.buildWorldMatrix(transform)
		ts.cache.matrices[entityID] = matrix
	}
}

// GetWorldMatrix retrieves the cached world matrix for an entity
func (ts *TransformSystem) GetWorldMatrix(id EntityID) (rl.Matrix, bool) {
	matrix, ok := ts.cache.matrices[id]
	return matrix, ok
}

// buildWorldMatrix constructs a transformation matrix from Transform component
func (ts *TransformSystem) buildWorldMatrix(t *Transform) rl.Matrix {
	// Build individual transformation matrices
	translation := rl.MatrixTranslate(t.Position.X, t.Position.Y, t.Position.Z)
	rotation := rl.QuaternionToMatrix(t.Rotation)
	scale := rl.MatrixScale(t.Scale.X, t.Scale.Y, t.Scale.Z)

	// Combine: Scale -> Rotate -> Translate
	matrix := rl.MatrixMultiply(scale, rotation)
	matrix = rl.MatrixMultiply(matrix, translation)

	return matrix
}
