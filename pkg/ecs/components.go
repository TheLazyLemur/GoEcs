package ecs

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Component type names as constants
const (
	ComponentTypeTransform  = "Transform"
	ComponentTypeRenderMesh = "RenderMesh"
	ComponentTypeName       = "Name"
	ComponentTypeCamera     = "Camera"
)

// Transform component holds position, rotation, and scale
type Transform struct {
	Position rl.Vector3
	Rotation rl.Quaternion
	Scale    rl.Vector3
}

// NewTransform creates a new Transform with default values
func NewTransform() *Transform {
	return &Transform{
		Position: rl.NewVector3(0, 0, 0),
		Rotation: rl.NewQuaternion(0, 0, 0, 1), // Identity quaternion
		Scale:    rl.NewVector3(1, 1, 1),
	}
}

// RenderMesh component holds references to mesh and material
type RenderMesh struct {
	MeshID     string
	MaterialID string
}

// NewRenderMesh creates a new RenderMesh component
func NewRenderMesh(meshID, materialID string) *RenderMesh {
	return &RenderMesh{
		MeshID:     meshID,
		MaterialID: materialID,
	}
}

// Name component gives an entity a human-readable name
type Name struct {
	Value string
}

// NewName creates a new Name component
func NewName(name string) *Name {
	return &Name{Value: name}
}

// Camera component marks an entity as a camera
type Camera struct {
	FOV        float32
	Near       float32
	Far        float32
	IsActive   bool
	Projection rl.CameraProjection
}

// NewCamera creates a new Camera component with default values
func NewCamera() *Camera {
	return &Camera{
		FOV:        45.0,
		Near:       0.1,
		Far:        1000.0,
		IsActive:   true,
		Projection: rl.CameraPerspective,
	}
}
