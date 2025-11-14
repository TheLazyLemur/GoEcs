package ecs

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// EntityID is a type alias for Donburi's Entity
type EntityID = *donburi.Entry

// World wraps Donburi's World
type World struct {
	world donburi.World
}

// NewWorld creates a new ECS world
func NewWorld() *World {
	return &World{
		world: donburi.NewWorld(),
	}
}

// GetDonburiWorld returns the underlying Donburi world
func (w *World) GetDonburiWorld() donburi.World {
	return w.world
}

// CreateEntity creates a new entity and returns its ID
func (w *World) CreateEntity() EntityID {
	return w.world.Entry(w.world.Create(NameComponent))
}

// DestroyEntity removes an entity and all its components
func (w *World) DestroyEntity(entity EntityID) {
	w.world.Remove(entity.Entity())
}

// EntityExists checks if an entity exists
func (w *World) EntityExists(entity EntityID) bool {
	return entity.Valid()
}

// GetEntitiesWithComponent returns all entities that have a specific component
// This method is deprecated, kept for compatibility - use direct queries instead
func (w *World) GetEntitiesWithComponent(componentType interface{}) []EntityID {
	// This method is not used in the current implementation
	// Direct queries should be used instead
	return []EntityID{}
}

// GetEntitiesWithComponents returns all entities that have all specified components
// This method is deprecated, kept for compatibility - use direct queries instead
func (w *World) GetEntitiesWithComponents(componentTypes ...interface{}) []EntityID {
	// This method is not used in the current implementation
	// Direct queries should be used instead
	return []EntityID{}
}

// GetAllEntities returns all entity IDs
func (w *World) GetAllEntities() []EntityID {
	entities := []EntityID{}
	// Query for entities with Name component (all entities should have one)
	q := query.NewQuery(filter.Contains(NameComponent))
	q.Each(w.world, func(entry *donburi.Entry) {
		entities = append(entities, entry)
	})
	return entities
}
