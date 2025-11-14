package math

import (
	gomath "math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ray represents a 3D ray with origin and direction
type Ray struct {
	Origin    rl.Vector3
	Direction rl.Vector3
}

// BoundingBox represents an axis-aligned bounding box
type BoundingBox struct {
	Min rl.Vector3
	Max rl.Vector3
}

// NewBoundingBox creates a new bounding box
func NewBoundingBox(min, max rl.Vector3) BoundingBox {
	return BoundingBox{Min: min, Max: max}
}

// GetDefaultMeshBounds returns default bounds for common mesh types
func GetDefaultMeshBounds(meshID string) BoundingBox {
	switch meshID {
	case "cube":
		return NewBoundingBox(
			rl.NewVector3(-0.5, -0.5, -0.5),
			rl.NewVector3(0.5, 0.5, 0.5),
		)
	case "sphere":
		return NewBoundingBox(
			rl.NewVector3(-0.5, -0.5, -0.5),
			rl.NewVector3(0.5, 0.5, 0.5),
		)
	case "plane":
		return NewBoundingBox(
			rl.NewVector3(-5.0, 0, -5.0),
			rl.NewVector3(5.0, 0.01, 5.0),
		)
	default:
		return NewBoundingBox(
			rl.NewVector3(-0.5, -0.5, -0.5),
			rl.NewVector3(0.5, 0.5, 0.5),
		)
	}
}

// TransformBoundingBox applies a transformation matrix to a bounding box
func TransformBoundingBox(box BoundingBox, matrix rl.Matrix) BoundingBox {
	// Transform all 8 corners of the box
	corners := []rl.Vector3{
		{box.Min.X, box.Min.Y, box.Min.Z},
		{box.Max.X, box.Min.Y, box.Min.Z},
		{box.Min.X, box.Max.Y, box.Min.Z},
		{box.Max.X, box.Max.Y, box.Min.Z},
		{box.Min.X, box.Min.Y, box.Max.Z},
		{box.Max.X, box.Min.Y, box.Max.Z},
		{box.Min.X, box.Max.Y, box.Max.Z},
		{box.Max.X, box.Max.Y, box.Max.Z},
	}

	transformedMin := rl.Vector3Transform(corners[0], matrix)
	transformedMax := transformedMin

	for i := 1; i < 8; i++ {
		transformed := rl.Vector3Transform(corners[i], matrix)
		transformedMin = Vector3Min(transformedMin, transformed)
		transformedMax = Vector3Max(transformedMax, transformed)
	}

	return NewBoundingBox(transformedMin, transformedMax)
}

// RayIntersectsBoundingBox tests if a ray intersects with an AABB
// Returns (hit, distance)
func RayIntersectsBoundingBox(ray Ray, box BoundingBox) (bool, float32) {
	tMin := float32(0.0)
	tMax := float32(1000000.0) // Large number

	// Test intersection with each pair of planes
	for i := 0; i < 3; i++ {
		var origin, direction, boxMin, boxMax float32

		switch i {
		case 0:
			origin = ray.Origin.X
			direction = ray.Direction.X
			boxMin = box.Min.X
			boxMax = box.Max.X
		case 1:
			origin = ray.Origin.Y
			direction = ray.Direction.Y
			boxMin = box.Min.Y
			boxMax = box.Max.Y
		case 2:
			origin = ray.Origin.Z
			direction = ray.Direction.Z
			boxMin = box.Min.Z
			boxMax = box.Max.Z
		}

		if abs32(direction) < 0.0001 {
			// Ray is parallel to slab
			if origin < boxMin || origin > boxMax {
				return false, 0
			}
		} else {
			// Compute intersection t values
			t1 := (boxMin - origin) / direction
			t2 := (boxMax - origin) / direction

			if t1 > t2 {
				t1, t2 = t2, t1
			}

			if t1 > tMin {
				tMin = t1
			}
			if t2 < tMax {
				tMax = t2
			}

			if tMin > tMax {
				return false, 0
			}
		}
	}

	return true, tMin
}

// Vector3Min returns component-wise minimum
func Vector3Min(a, b rl.Vector3) rl.Vector3 {
	return rl.NewVector3(
		min32(a.X, b.X),
		min32(a.Y, b.Y),
		min32(a.Z, b.Z),
	)
}

// Vector3Max returns component-wise maximum
func Vector3Max(a, b rl.Vector3) rl.Vector3 {
	return rl.NewVector3(
		max32(a.X, b.X),
		max32(a.Y, b.Y),
		max32(a.Z, b.Z),
	)
}

func min32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func abs32(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

// GetMouseRay generates a ray from the camera through the mouse position
func GetMouseRay(mousePos rl.Vector2, camera rl.Camera3D, viewport rl.Rectangle) Ray {
	// Normalize mouse coordinates to [-1, 1] range
	x := (2.0*mousePos.X)/viewport.Width - 1.0
	y := 1.0 - (2.0*mousePos.Y)/viewport.Height

	// Get camera matrices
	view := rl.GetCameraMatrix(camera)
	projection := rl.GetCameraMatrix2D(rl.Camera2D{
		Offset: rl.NewVector2(0, 0),
		Target: rl.NewVector2(0, 0),
		Rotation: 0,
		Zoom: 1,
	}) // Placeholder - we'll manually create projection

	// Create proper projection matrix
	aspect := viewport.Width / viewport.Height
	projection = createPerspectiveProjection(camera.Fovy, aspect, 0.01, 1000.0)

	// Calculate inverse matrices
	viewProjection := rl.MatrixMultiply(view, projection)
	invViewProjection := rl.MatrixInvert(viewProjection)

	// Near and far points in NDC
	nearPoint := rl.NewVector3(x, y, -1.0)
	farPoint := rl.NewVector3(x, y, 1.0)

	// Transform to world space
	nearWorld := rl.Vector3Transform(nearPoint, invViewProjection)
	farWorld := rl.Vector3Transform(farPoint, invViewProjection)

	// Calculate ray direction
	direction := rl.Vector3Subtract(farWorld, nearWorld)
	direction = rl.Vector3Normalize(direction)

	return Ray{
		Origin:    camera.Position,
		Direction: direction,
	}
}

// createPerspectiveProjection creates a perspective projection matrix
func createPerspectiveProjection(fovy, aspect, near, far float32) rl.Matrix {
	top := near * float32(gomath.Tan(float64(fovy*rl.Deg2rad/2.0)))
	bottom := -top
	right := top * aspect
	left := -right

	return rl.MatrixFrustum(left, right, bottom, top, near, far)
}
