package ecs

import (
	"sync"
)

// EntityID is a unique identifier for an entity
type EntityID uint64

// Component is a marker interface for all components
type Component interface{}

// World holds all entities and their components
type World struct {
	nextEntityID EntityID
	entities     map[EntityID]bool
	components   map[EntityID]map[string]Component
	mu           sync.RWMutex
}

// NewWorld creates a new ECS world
func NewWorld() *World {
	return &World{
		nextEntityID: 1,
		entities:     make(map[EntityID]bool),
		components:   make(map[EntityID]map[string]Component),
	}
}

// CreateEntity creates a new entity and returns its ID
func (w *World) CreateEntity() EntityID {
	w.mu.Lock()
	defer w.mu.Unlock()

	id := w.nextEntityID
	w.nextEntityID++
	w.entities[id] = true
	w.components[id] = make(map[string]Component)
	return id
}

// DestroyEntity removes an entity and all its components
func (w *World) DestroyEntity(id EntityID) {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.entities, id)
	delete(w.components, id)
}

// AddComponent adds a component to an entity
func (w *World) AddComponent(id EntityID, componentType string, component Component) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.entities[id]; !exists {
		return
	}

	w.components[id][componentType] = component
}

// RemoveComponent removes a component from an entity
func (w *World) RemoveComponent(id EntityID, componentType string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if comps, exists := w.components[id]; exists {
		delete(comps, componentType)
	}
}

// GetComponent retrieves a component from an entity
func (w *World) GetComponent(id EntityID, componentType string) (Component, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if comps, exists := w.components[id]; exists {
		comp, ok := comps[componentType]
		return comp, ok
	}
	return nil, false
}

// HasComponent checks if an entity has a specific component
func (w *World) HasComponent(id EntityID, componentType string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if comps, exists := w.components[id]; exists {
		_, ok := comps[componentType]
		return ok
	}
	return false
}

// GetAllEntities returns all entity IDs
func (w *World) GetAllEntities() []EntityID {
	w.mu.RLock()
	defer w.mu.RUnlock()

	entities := make([]EntityID, 0, len(w.entities))
	for id := range w.entities {
		entities = append(entities, id)
	}
	return entities
}

// GetEntitiesWithComponent returns all entities that have a specific component
func (w *World) GetEntitiesWithComponent(componentType string) []EntityID {
	w.mu.RLock()
	defer w.mu.RUnlock()

	entities := make([]EntityID, 0)
	for id, comps := range w.components {
		if _, ok := comps[componentType]; ok {
			entities = append(entities, id)
		}
	}
	return entities
}

// GetEntitiesWithComponents returns all entities that have all specified components
func (w *World) GetEntitiesWithComponents(componentTypes ...string) []EntityID {
	w.mu.RLock()
	defer w.mu.RUnlock()

	entities := make([]EntityID, 0)
	for id, comps := range w.components {
		hasAll := true
		for _, ct := range componentTypes {
			if _, ok := comps[ct]; !ok {
				hasAll = false
				break
			}
		}
		if hasAll {
			entities = append(entities, id)
		}
	}
	return entities
}

// EntityExists checks if an entity exists
func (w *World) EntityExists(id EntityID) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.entities[id]
}
