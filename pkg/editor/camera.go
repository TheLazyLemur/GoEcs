package editor

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// EditorCamera represents the camera used in editor mode
type EditorCamera struct {
	Position rl.Vector3
	Target   rl.Vector3
	Up       rl.Vector3
	FOV      float32

	// Control settings
	MoveSpeed   float32
	LookSpeed   float32
	MouseSens   float32

	// Internal state
	yaw   float32
	pitch float32
}

// NewEditorCamera creates a new editor camera with default values
func NewEditorCamera() *EditorCamera {
	return &EditorCamera{
		Position:  rl.NewVector3(10, 10, 10),
		Target:    rl.NewVector3(0, 0, 0),
		Up:        rl.NewVector3(0, 1, 0),
		FOV:       45.0,
		MoveSpeed: 10.0,
		LookSpeed: 100.0,
		MouseSens: 0.003,
		yaw:       -135.0, // Looking at origin from (10, 10, 10)
		pitch:     -35.0,
	}
}

// Update handles camera movement and rotation
func (ec *EditorCamera) Update(dt float32) {
	// Mouse look (right mouse button)
	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		mouseDelta := rl.GetMouseDelta()

		ec.yaw += mouseDelta.X * ec.MouseSens * 100.0
		ec.pitch -= mouseDelta.Y * ec.MouseSens * 100.0

		// Clamp pitch to avoid gimbal lock
		if ec.pitch > 89.0 {
			ec.pitch = 89.0
		}
		if ec.pitch < -89.0 {
			ec.pitch = -89.0
		}
	}

	// Calculate forward, right vectors from yaw/pitch
	forward := ec.GetForward()
	right := rl.Vector3Normalize(rl.Vector3CrossProduct(forward, ec.Up))

	speed := ec.MoveSpeed * dt

	// Speed boost with shift
	if rl.IsKeyDown(rl.KeyLeftShift) {
		speed *= 3.0
	}

	// WASD movement
	if rl.IsKeyDown(rl.KeyW) {
		ec.Position = rl.Vector3Add(ec.Position, rl.Vector3Scale(forward, speed))
	}
	if rl.IsKeyDown(rl.KeyS) {
		ec.Position = rl.Vector3Subtract(ec.Position, rl.Vector3Scale(forward, speed))
	}
	if rl.IsKeyDown(rl.KeyD) {
		ec.Position = rl.Vector3Add(ec.Position, rl.Vector3Scale(right, speed))
	}
	if rl.IsKeyDown(rl.KeyA) {
		ec.Position = rl.Vector3Subtract(ec.Position, rl.Vector3Scale(right, speed))
	}

	// Vertical movement
	if rl.IsKeyDown(rl.KeyE) {
		ec.Position.Y += speed
	}
	if rl.IsKeyDown(rl.KeyQ) {
		ec.Position.Y -= speed
	}

	// Update target to look at
	ec.Target = rl.Vector3Add(ec.Position, forward)
}

// GetForward returns the forward direction vector
func (ec *EditorCamera) GetForward() rl.Vector3 {
	yawRad := ec.yaw * rl.Deg2rad
	pitchRad := ec.pitch * rl.Deg2rad

	forward := rl.NewVector3(
		float32(math.Cos(float64(yawRad))*math.Cos(float64(pitchRad))),
		float32(math.Sin(float64(pitchRad))),
		float32(math.Sin(float64(yawRad))*math.Cos(float64(pitchRad))),
	)

	return rl.Vector3Normalize(forward)
}

// GetRaylibCamera converts the editor camera to raylib Camera3D
func (ec *EditorCamera) GetRaylibCamera() rl.Camera3D {
	return rl.NewCamera3D(
		ec.Position,
		ec.Target,
		ec.Up,
		ec.FOV,
		rl.CameraPerspective,
	)
}
