package ecs

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Scene wraps the ECS world and provides a controlled API for mutations
type Scene struct {
	world *World
	name  string
}

// NewScene creates a new scene
func NewScene(name string) *Scene {
	return &Scene{
		world: NewWorld(),
		name:  name,
	}
}

// GetWorld returns the underlying ECS world (read-only access)
func (s *Scene) GetWorld() *World {
	return s.world
}

// GetSceneName returns the scene name
func (s *Scene) GetSceneName() string {
	return s.name
}

// CreateEntity creates a new entity and returns its ID
func (s *Scene) CreateEntity(name string) EntityID {
	id := s.world.CreateEntity()
	s.world.AddComponent(id, ComponentTypeName, NewName(name))
	return id
}

// DestroyEntity removes an entity from the scene
func (s *Scene) DestroyEntity(id EntityID) {
	s.world.DestroyEntity(id)
}

// AddTransform adds a Transform component to an entity
func (s *Scene) AddTransform(id EntityID) {
	s.world.AddComponent(id, ComponentTypeTransform, NewTransform())
}

// AddRenderMesh adds a RenderMesh component to an entity
func (s *Scene) AddRenderMesh(id EntityID, meshID, materialID string) {
	s.world.AddComponent(id, ComponentTypeRenderMesh, NewRenderMesh(meshID, materialID))
}

// AddCamera adds a Camera component to an entity
func (s *Scene) AddCamera(id EntityID) {
	s.world.AddComponent(id, ComponentTypeCamera, NewCamera())
}

// RemoveComponent removes a component from an entity
func (s *Scene) RemoveComponent(id EntityID, componentType string) {
	s.world.RemoveComponent(id, componentType)
}

// SetTransform updates the Transform component of an entity
func (s *Scene) SetTransform(id EntityID, position, scale rl.Vector3, rotation rl.Quaternion) bool {
	if !s.world.HasComponent(id, ComponentTypeTransform) {
		return false
	}

	transform := &Transform{
		Position: position,
		Rotation: rotation,
		Scale:    scale,
	}
	s.world.AddComponent(id, ComponentTypeTransform, transform)
	return true
}

// GetTransform retrieves the Transform component of an entity
func (s *Scene) GetTransform(id EntityID) (*Transform, bool) {
	comp, ok := s.world.GetComponent(id, ComponentTypeTransform)
	if !ok {
		return nil, false
	}
	transform, ok := comp.(*Transform)
	return transform, ok
}

// SetName updates the Name component of an entity
func (s *Scene) SetName(id EntityID, name string) bool {
	if !s.world.HasComponent(id, ComponentTypeName) {
		return false
	}

	s.world.AddComponent(id, ComponentTypeName, NewName(name))
	return true
}

// GetName retrieves the Name component of an entity
func (s *Scene) GetName(id EntityID) (string, bool) {
	comp, ok := s.world.GetComponent(id, ComponentTypeName)
	if !ok {
		return "", false
	}
	name, ok := comp.(*Name)
	if !ok {
		return "", false
	}
	return name.Value, true
}

// SetRenderMesh updates the RenderMesh component of an entity
func (s *Scene) SetRenderMesh(id EntityID, meshID, materialID string) bool {
	if !s.world.HasComponent(id, ComponentTypeRenderMesh) {
		return false
	}

	s.world.AddComponent(id, ComponentTypeRenderMesh, NewRenderMesh(meshID, materialID))
	return true
}

// GetRenderMesh retrieves the RenderMesh component of an entity
func (s *Scene) GetRenderMesh(id EntityID) (*RenderMesh, bool) {
	comp, ok := s.world.GetComponent(id, ComponentTypeRenderMesh)
	if !ok {
		return nil, false
	}
	mesh, ok := comp.(*RenderMesh)
	return mesh, ok
}

// GetCamera retrieves the Camera component of an entity
func (s *Scene) GetCamera(id EntityID) (*Camera, bool) {
	comp, ok := s.world.GetComponent(id, ComponentTypeCamera)
	if !ok {
		return nil, false
	}
	camera, ok := comp.(*Camera)
	return camera, ok
}

// GetAllEntities returns all entity IDs in the scene
func (s *Scene) GetAllEntities() []EntityID {
	return s.world.GetAllEntities()
}

// HasComponent checks if an entity has a specific component
func (s *Scene) HasComponent(id EntityID, componentType string) bool {
	return s.world.HasComponent(id, componentType)
}
