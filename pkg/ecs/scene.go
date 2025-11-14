package ecs

import (
	"github.com/yohamta/donburi"
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
	entity := s.world.CreateEntity()
	// Set the name
	NameComponent.Set(entity, NewName(name))
	return entity
}

// DestroyEntity removes an entity from the scene
func (s *Scene) DestroyEntity(entity EntityID) {
	s.world.DestroyEntity(entity)
}

// AddTransform adds a Transform component to an entity
func (s *Scene) AddTransform(entity EntityID) {
	TransformComponent.Set(entity, NewTransform())
}

// AddRenderMesh adds a RenderMesh component to an entity
func (s *Scene) AddRenderMesh(entity EntityID, meshID, materialID string) {
	RenderMeshComponent.Set(entity, NewRenderMesh(meshID, materialID))
}

// AddCamera adds a Camera component to an entity
func (s *Scene) AddCamera(entity EntityID) {
	CameraComponent.Set(entity, NewCamera())
}

// RemoveTransform removes the Transform component from an entity
func (s *Scene) RemoveTransform(entity EntityID) {
	if entity.HasComponent(TransformComponent) {
		donburi.Remove[Transform](entity, TransformComponent)
	}
}

// RemoveRenderMesh removes the RenderMesh component from an entity
func (s *Scene) RemoveRenderMesh(entity EntityID) {
	if entity.HasComponent(RenderMeshComponent) {
		donburi.Remove[RenderMesh](entity, RenderMeshComponent)
	}
}

// RemoveCamera removes the Camera component from an entity
func (s *Scene) RemoveCamera(entity EntityID) {
	if entity.HasComponent(CameraComponent) {
		donburi.Remove[Camera](entity, CameraComponent)
	}
}

// SetTransform updates the Transform component of an entity
func (s *Scene) SetTransform(entity EntityID, position, scale rl.Vector3, rotation rl.Quaternion) bool {
	if !entity.HasComponent(TransformComponent) {
		return false
	}

	transform := &Transform{
		Position: position,
		Rotation: rotation,
		Scale:    scale,
	}
	TransformComponent.Set(entity, transform)
	return true
}

// GetTransform retrieves the Transform component of an entity
func (s *Scene) GetTransform(entity EntityID) (*Transform, bool) {
	if !entity.HasComponent(TransformComponent) {
		return nil, false
	}
	transform := TransformComponent.Get(entity)
	return transform, true
}

// SetName updates the Name component of an entity
func (s *Scene) SetName(entity EntityID, name string) bool {
	if !entity.HasComponent(NameComponent) {
		return false
	}

	NameComponent.Set(entity, NewName(name))
	return true
}

// GetName retrieves the Name component of an entity
func (s *Scene) GetName(entity EntityID) (string, bool) {
	if !entity.HasComponent(NameComponent) {
		return "", false
	}
	name := NameComponent.Get(entity)
	return name.Value, true
}

// SetRenderMesh updates the RenderMesh component of an entity
func (s *Scene) SetRenderMesh(entity EntityID, meshID, materialID string) bool {
	if !entity.HasComponent(RenderMeshComponent) {
		return false
	}

	RenderMeshComponent.Set(entity, NewRenderMesh(meshID, materialID))
	return true
}

// GetRenderMesh retrieves the RenderMesh component of an entity
func (s *Scene) GetRenderMesh(entity EntityID) (*RenderMesh, bool) {
	if !entity.HasComponent(RenderMeshComponent) {
		return nil, false
	}
	mesh := RenderMeshComponent.Get(entity)
	return mesh, true
}

// GetCamera retrieves the Camera component of an entity
func (s *Scene) GetCamera(entity EntityID) (*Camera, bool) {
	if !entity.HasComponent(CameraComponent) {
		return nil, false
	}
	camera := CameraComponent.Get(entity)
	return camera, true
}

// GetAllEntities returns all entity IDs in the scene
func (s *Scene) GetAllEntities() []EntityID {
	return s.world.GetAllEntities()
}

// HasComponent checks if an entity has a specific component type
func (s *Scene) HasComponent(entity EntityID, componentType string) bool {
	// Map old string-based component types to new Donburi component types
	switch componentType {
	case "Transform":
		return entity.HasComponent(TransformComponent)
	case "RenderMesh":
		return entity.HasComponent(RenderMeshComponent)
	case "Name":
		return entity.HasComponent(NameComponent)
	case "Camera":
		return entity.HasComponent(CameraComponent)
	default:
		return false
	}
}
