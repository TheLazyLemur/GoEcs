package editor

import (
	"github.com/TheLazyLemur/SpaceImpact/pkg/ecs"
	pkgmath "github.com/TheLazyLemur/SpaceImpact/pkg/math"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// GizmoMode represents the current gizmo operation mode
type GizmoMode int

const (
	GizmoModeNone GizmoMode = iota
	GizmoModeTranslate
	GizmoModeRotate
	GizmoModeScale
)

// GizmoAxis represents which axis is being manipulated
type GizmoAxis int

const (
	GizmoAxisNone GizmoAxis = iota
	GizmoAxisX
	GizmoAxisY
	GizmoAxisZ
)

// Gizmo handles 3D manipulation of entities
type Gizmo struct {
	Mode         GizmoMode
	ActiveAxis   GizmoAxis
	IsDragging   bool
	DragStartPos rl.Vector2
	EntityStart  rl.Vector3
}

// NewGizmo creates a new gizmo
func NewGizmo() *Gizmo {
	return &Gizmo{
		Mode:       GizmoModeTranslate,
		ActiveAxis: GizmoAxisNone,
		IsDragging: false,
	}
}

// DrawTranslationGizmo draws the translation gizmo at the entity position
func (g *Gizmo) DrawTranslationGizmo(position rl.Vector3, camera *EditorCamera) {
	gizmoSize := float32(1.0)
	arrowLength := gizmoSize
	arrowRadius := gizmoSize * 0.05
	coneHeight := gizmoSize * 0.2
	coneRadius := gizmoSize * 0.1

	// X axis (Red)
	xEnd := rl.Vector3Add(position, rl.NewVector3(arrowLength, 0, 0))
	xColor := rl.Red
	if g.ActiveAxis == GizmoAxisX {
		xColor = rl.Yellow
	}
	g.drawArrow(position, xEnd, arrowRadius, coneHeight, coneRadius, xColor)

	// Y axis (Green)
	yEnd := rl.Vector3Add(position, rl.NewVector3(0, arrowLength, 0))
	yColor := rl.Green
	if g.ActiveAxis == GizmoAxisY {
		yColor = rl.Yellow
	}
	g.drawArrow(position, yEnd, arrowRadius, coneHeight, coneRadius, yColor)

	// Z axis (Blue)
	zEnd := rl.Vector3Add(position, rl.NewVector3(0, 0, arrowLength))
	zColor := rl.Blue
	if g.ActiveAxis == GizmoAxisZ {
		zColor = rl.Yellow
	}
	g.drawArrow(position, zEnd, arrowRadius, coneHeight, coneRadius, zColor)
}

// drawArrow draws an arrow from start to end
func (g *Gizmo) drawArrow(start, end rl.Vector3, radius, coneHeight, coneRadius float32, color rl.Color) {
	// Draw cylinder for shaft
	direction := rl.Vector3Subtract(end, start)
	length := rl.Vector3Length(direction)
	shaftLength := length - coneHeight

	// Calculate midpoint for cylinder
	dir := rl.Vector3Normalize(direction)
	shaftEnd := rl.Vector3Add(start, rl.Vector3Scale(dir, shaftLength))

	// Draw shaft
	rl.DrawLine3D(start, shaftEnd, color)

	// Draw cone for arrowhead
	rl.DrawSphere(end, coneRadius, color)
}

// Update handles gizmo interaction
func (g *Gizmo) Update(
	viewport *Viewport,
	state *EditorState,
	transformSystem *ecs.TransformSystem,
) {
	if !state.HasSelection() {
		g.IsDragging = false
		g.ActiveAxis = GizmoAxisNone
		return
	}

	// Get first selected entity position
	entityID := state.GetFirstSelected()
	transform, ok := state.Scene.GetTransform(entityID)
	if !ok {
		return
	}

	mousePos := rl.GetMousePosition()

	// Start dragging
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && viewport.IsMouseOver() {
		// Check which axis was clicked
		axis := g.pickGizmoAxis(viewport, transform.Position, mousePos)
		if axis != GizmoAxisNone {
			g.IsDragging = true
			g.ActiveAxis = axis
			g.DragStartPos = mousePos
			g.EntityStart = transform.Position
		}
	}

	// Update dragging
	if g.IsDragging && rl.IsMouseButtonDown(rl.MouseLeftButton) {
		mouseDelta := rl.Vector2Subtract(mousePos, g.DragStartPos)
		g.updateTranslation(state, entityID, mouseDelta, viewport)
	}

	// Stop dragging
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		g.IsDragging = false
		g.ActiveAxis = GizmoAxisNone
	}
}

// pickGizmoAxis determines which gizmo axis was clicked
func (g *Gizmo) pickGizmoAxis(viewport *Viewport, position rl.Vector3, mousePos rl.Vector2) GizmoAxis {
	// Convert mouse position to viewport-relative coordinates
	viewportMouseX := mousePos.X - viewport.Bounds.X
	viewportMouseY := mousePos.Y - viewport.Bounds.Y

	// Create viewport-relative mouse position
	viewportMouse := rl.NewVector2(viewportMouseX, viewportMouseY)

	// Generate ray
	camera := viewport.Camera.GetRaylibCamera()
	viewportRect := rl.NewRectangle(0, 0, viewport.Bounds.Width, viewport.Bounds.Height)
	ray := pkgmath.GetMouseRay(viewportMouse, camera, viewportRect)

	gizmoSize := float32(1.0)
	threshold := float32(0.2)

	// Test X axis
	xEnd := rl.Vector3Add(position, rl.NewVector3(gizmoSize, 0, 0))
	xBounds := pkgmath.NewBoundingBox(
		rl.NewVector3(position.X-threshold, position.Y-threshold, position.Z-threshold),
		rl.NewVector3(xEnd.X+threshold, position.Y+threshold, position.Z+threshold),
	)
	if hit, _ := pkgmath.RayIntersectsBoundingBox(ray, xBounds); hit {
		return GizmoAxisX
	}

	// Test Y axis
	yEnd := rl.Vector3Add(position, rl.NewVector3(0, gizmoSize, 0))
	yBounds := pkgmath.NewBoundingBox(
		rl.NewVector3(position.X-threshold, position.Y-threshold, position.Z-threshold),
		rl.NewVector3(position.X+threshold, yEnd.Y+threshold, position.Z+threshold),
	)
	if hit, _ := pkgmath.RayIntersectsBoundingBox(ray, yBounds); hit {
		return GizmoAxisY
	}

	// Test Z axis
	zEnd := rl.Vector3Add(position, rl.NewVector3(0, 0, gizmoSize))
	zBounds := pkgmath.NewBoundingBox(
		rl.NewVector3(position.X-threshold, position.Y-threshold, position.Z-threshold),
		rl.NewVector3(position.X+threshold, position.Y+threshold, zEnd.Z+threshold),
	)
	if hit, _ := pkgmath.RayIntersectsBoundingBox(ray, zBounds); hit {
		return GizmoAxisZ
	}

	return GizmoAxisNone
}

// updateTranslation updates entity position based on mouse movement
func (g *Gizmo) updateTranslation(
	state *EditorState,
	entityID ecs.EntityID,
	mouseDelta rl.Vector2,
	viewport *Viewport,
) {
	// Simple translation based on screen-space mouse movement
	sensitivity := float32(0.01)
	offset := rl.Vector3Zero()

	switch g.ActiveAxis {
	case GizmoAxisX:
		offset.X = mouseDelta.X * sensitivity
	case GizmoAxisY:
		offset.Y = -mouseDelta.Y * sensitivity // Invert Y for intuitive up/down
	case GizmoAxisZ:
		offset.Z = mouseDelta.Y * sensitivity
	}

	// Apply offset to entity
	newPos := rl.Vector3Add(g.EntityStart, offset)

	transform, ok := state.Scene.GetTransform(entityID)
	if ok {
		state.Scene.SetTransform(
			entityID,
			newPos,
			transform.Scale,
			transform.Rotation,
		)
	}
}
