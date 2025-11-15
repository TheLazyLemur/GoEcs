# GoEcs Engine + Editor Technical Specification

## Version 3.0

## Architecture Overview

This project implements a Go-based game engine with an integrated editor, following a layered architecture that separates the game world (ECS-based) from the editor tooling (non-ECS based).

### Layer Architecture

```
┌─────────────────────────────────────────────────────────┐
│ Layer 3: Editor Layer (Non-ECS)                        │
│  - EditorCamera, Viewport, Gizmos, Panels              │
│  - Manipulates Layer 2 ONLY through Scene API          │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼ (Scene API - Single Mutation Entry Point)
┌─────────────────────────────────────────────────────────┐
│ Layer 2: Game World (Donburi ECS)                      │
│  - Entities, Components, Systems                        │
│  - All game state lives here                           │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│ Layer 1: Rendering & Assets (raylib-go)                │
│  - Meshes, Materials, Textures                         │
│  - Low-level rendering primitives                      │
└─────────────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Separation of Concerns**: Editor is NOT ECS-based. It's a separate layer that observes and manipulates the ECS world.
2. **Single Mutation Entry Point**: All ECS modifications go through the Scene API wrapper.
3. **Type Safety**: Donburi ECS provides compile-time type safety for components.
4. **Immutability**: Editor state is separate from game state.

## Layer 2: Game World (ECS)

### ECS Implementation: Donburi

We use [Donburi](https://github.com/yohamta/donburi) v1.15.7 as our ECS foundation.

**Why Donburi?**
- Clean, lightweight, idiomatic Go
- Type-safe component access via generics
- Query-based entity iteration
- No reflection overhead
- Battle-tested in production games

### Core Types

#### EntityID
```go
type EntityID = *donburi.Entry
```
Entities are represented as pointers to Donburi entries. This allows:
- Direct component access via entry methods
- Efficient entity validity checks
- Type-safe component operations

#### Components

All components are defined using Donburi's `ComponentType[T]`:

```go
var (
    TransformComponent  = donburi.NewComponentType[Transform]()
    RenderMeshComponent = donburi.NewComponentType[RenderMesh]()
    NameComponent       = donburi.NewComponentType[Name]()
    CameraComponent     = donburi.NewComponentType[Camera]()
)
```

**Transform Component**
```go
type Transform struct {
    Position rl.Vector3
    Rotation rl.Quaternion  // Identity: (0,0,0,1)
    Scale    rl.Vector3     // Default: (1,1,1)
}
```
- Pure data struct
- Rotation uses quaternions for robust 3D rotation
- Position/Scale use raylib's Vector3

**RenderMesh Component**
```go
type RenderMesh struct {
    MeshID     string  // Reference to mesh registry
    MaterialID string  // Reference to material registry
}
```
- References GPU resources by string ID
- Actual mesh/material data lives in asset registries
- Allows multiple entities to share meshes efficiently

**Name Component**
```go
type Name struct {
    Value string
}
```
- Human-readable entity identifier
- All entities have a Name component by default
- Used for hierarchy panel display

**Camera Component**
```go
type Camera struct {
    FOV        float32
    Near       float32
    Far        float32
    IsActive   bool
    Projection rl.CameraProjection
}
```
- Defines camera properties for rendering
- Multiple cameras can exist (only one active)
- Currently unused but planned for scene cameras

### Systems

#### TransformSystem

**Purpose**: Builds world transformation matrices for all entities with Transform components.

**Architecture**:
```go
type TransformSystem struct {
    cache *TransformCache  // map[donburi.Entity]rl.Matrix
}
```

**Algorithm**:
1. Query all entities with `TransformComponent`
2. For each entity:
   - Extract position, rotation, scale
   - Build TRS matrix: `Matrix = Translate * Rotate * Scale`
   - Cache result by entity ID

**Usage**:
```go
transformSystem.Update(scene)
matrix, ok := transformSystem.GetWorldMatrix(entity)
```

**Matrix Composition Order**:
```
Local Space → Scale → Rotate → Translate → World Space
```

#### RenderSystem

**Purpose**: Renders all entities with both Transform and RenderMesh components.

**Algorithm**:
1. Query entities with `TransformComponent` AND `RenderMeshComponent`
2. For each entity:
   - Get mesh from asset registry (by MeshID)
   - Get material from asset registry (by MaterialID)
   - Get world matrix from TransformSystem cache
   - Call `rl.DrawMesh(mesh, material, worldMatrix)`
3. If mesh/material missing, draw magenta placeholder

**Missing Asset Handling**:
- Mesh not found → Draw magenta wireframe cube
- Material not found → Use default material

### Scene API (Mutation Wrapper)

The Scene API is the **ONLY** way to modify the ECS world from editor code.

**Entity Lifecycle**:
```go
entity := scene.CreateEntity("MyEntity")  // Creates entity with Name component
scene.AddTransform(entity)                // Add Transform component
scene.AddRenderMesh(entity, "cube", "default")  // Add RenderMesh
scene.DestroyEntity(entity)               // Remove entity and all components
```

**Component Access**:
```go
// Getters
transform, ok := scene.GetTransform(entity)
name, ok := scene.GetName(entity)
renderMesh, ok := scene.GetRenderMesh(entity)

// Setters
scene.SetTransform(entity, position, scale, rotation)
scene.SetName(entity, "NewName")
scene.SetRenderMesh(entity, "sphere", "red")

// Removal
scene.RemoveTransform(entity)
scene.RemoveRenderMesh(entity)
```

**Query API**:
```go
entities := scene.GetAllEntities()
hasTransform := scene.HasComponent(entity, "Transform")
```

**Design Rationale**:
- Prevents direct World access from editor
- Allows validation and logging of mutations
- Future: Can add undo/redo tracking here
- Future: Can add change notifications for editor updates

## Layer 3: Editor Layer (Non-ECS)

### Editor Architecture

The editor is **NOT** ECS-based. It's a traditional OOP system that manipulates the ECS world through the Scene API.

```
Editor Components:
├── EditorState (selection, mode, viewports)
├── EditorCamera (FPS-style camera, WASD controls)
├── Viewport (RenderTexture pipeline)
├── Gizmos (translation/rotation/scale handles)
└── Panels
    ├── HierarchyPanel (entity list)
    └── InspectorPanel (component editor)
```

### EditorState

**Purpose**: Central editor state container (non-ECS).

```go
type EditorState struct {
    Mode           EditorMode  // Edit, Play, Pause
    Scene          *ecs.Scene
    Selected       []ecs.EntityID
    Viewports      []*Viewport
    ShowHierarchy  bool
    ShowInspector  bool
}
```

**Selection Management**:
```go
state.SelectEntity(entity)      // Replace selection
state.AddToSelection(entity)    // Multi-select (Ctrl+Click)
state.ClearSelection()
state.IsSelected(entity) bool
state.GetFirstSelected() EntityID
```

### EditorCamera

**Purpose**: Free-flying FPS camera for viewport navigation.

**Controls**:
- WASD: Move horizontally
- Q/E: Move up/down
- Right Mouse + Drag: Look around
- Mouse Wheel: Adjust move speed

**Implementation**:
```go
type EditorCamera struct {
    Position   rl.Vector3
    yaw, pitch float32
    MoveSpeed  float32
    MouseSens  float32
}
```

**Note**: This is separate from the ECS Camera component. The editor camera is for viewing/editing, while Camera components are for in-game cameras.

### Viewport System

**Purpose**: Render-to-texture pipeline for viewport rendering.

**Architecture**:
```go
type Viewport struct {
    Bounds        rl.Rectangle     // Screen-space bounds
    RenderTexture rl.RenderTexture2D
    Camera        *EditorCamera
}
```

**Rendering Pipeline**:
```
1. Begin RenderTexture
2. BeginMode3D (editor camera)
3. Draw grid
4. RenderSystem.Render(scene)
5. Draw selection highlights
6. Draw gizmos
7. EndMode3D
8. End RenderTexture
9. Draw RenderTexture to main window
```

**Benefits**:
- Can apply post-processing effects
- Multiple viewports possible
- Easy viewport resolution changes
- Viewport-relative mouse picking

### Picking System

**Purpose**: Raycasting for viewport entity selection.

**Algorithm**:
1. Convert mouse position to viewport-relative coordinates
2. Generate ray from camera through mouse position
3. For each entity with Transform + RenderMesh:
   - Get mesh bounding box (local space)
   - Transform AABB to world space using world matrix
   - Test ray-AABB intersection
   - Track closest hit with positive distance
4. Return closest entity (if any)

**Ray Generation**:
```go
ray := pkgmath.GetMouseRay(viewportMouse, camera, viewportRect)
```

**Behind-Camera Filtering**:
```go
if distance > 0 && distance < closestHit.Distance {
    // Only consider hits in front of camera
}
```

### Gizmo System

**Purpose**: 3D manipulation handles for entity transforms.

**Current Implementation**: Translation Gizmo (Phase 8)

**Architecture**:
```go
type Gizmo struct {
    mode       GizmoMode  // Translation, Rotation, Scale
    activeAxis GizmoAxis  // X, Y, Z, None
}
```

**Translation Gizmo**:
- Three colored arrows: Red (X), Green (Y), Blue (Z)
- Click and drag to move entity along axis
- Sensitivity: 0.05 units per pixel
- Visual feedback: Highlight active axis

**Interaction Flow**:
1. Hover detection (raycast to axis handles)
2. Click to grab axis
3. Drag to move (constrained to axis)
4. Update entity transform via Scene API
5. Release to finish

**Planned**: Rotation and Scale gizmos (Phase 13)

### UI Panels

#### HierarchyPanel

**Purpose**: Display all entities in the scene as a scrollable list.

**Features**:
- Entity name display
- Selection highlighting (blue background)
- Click to select (Ctrl+Click for multi-select)
- Hover effects
- Scroll support for long lists

**Current Rendering**: Custom drawing with raylib primitives
**Planned**: Migrate to raygui ListView (Phase 9+)

#### InspectorPanel

**Purpose**: Display and edit components for selected entity.

**Layout**:
```
┌─────────────────────────┐
│ Inspector - EntityName  │ ← Title bar
├─────────────────────────┤
│ Name Component          │ ← Collapsible header
│   Name: "Cube"          │ ← Read-only (for now)
├─────────────────────────┤
│ Transform Component     │
│   Position: X Y Z       │
│   Scale: X Y Z          │
│   Rotation: X Y Z W     │
├─────────────────────────┤
│ RenderMesh Component    │
│   Mesh: cube            │
│   Material: default     │
└─────────────────────────┘
```

**Current Rendering**: Custom drawing with raylib primitives
**Planned**: Migrate to raygui controls (TextBox, Spinner, etc.) for editing

## Layer 1: Rendering & Assets

### Asset Management

**Purpose**: Centralized registry for GPU resources.

```go
type AssetManager struct {
    meshes    *MeshRegistry
    materials *MaterialRegistry
    textures  *TextureRegistry
}
```

**Registry Pattern**:
```go
type MeshRegistry struct {
    meshes map[string]rl.Mesh
    mu     sync.RWMutex
}

registry.Register("cube", rl.GenMeshCube(1, 1, 1))
mesh, ok := registry.Get("cube")
```

**Default Assets**:
- Meshes: cube, sphere, plane, cylinder
- Materials: default (white), red, green, blue
- Textures: (none yet)

**Thread Safety**: All registries use RWMutex for concurrent access.

### Rendering Pipeline

**Main Loop**:
```
1. Input handling (camera, picking, gizmo)
2. Update systems (TransformSystem)
3. Begin Viewport RenderTexture
   4. Begin 3D mode
   5. Draw grid
   6. RenderSystem.Render(scene)
   7. Draw selection highlights
   8. Draw gizmos
   9. End 3D mode
10. End RenderTexture
11. Begin main window
12. Draw RenderTexture
13. Draw UI panels (Hierarchy, Inspector)
14. End main window
```

### Coordinate System

Following raylib conventions:
- Right-handed coordinate system
- +X: Right
- +Y: Up
- +Z: Forward (towards camera)

### Dependencies

```go
require (
    github.com/gen2brain/raylib-go/raylib v0.55.1
    github.com/yohamta/donburi v1.15.7
)
```

**raylib-go**: Rendering, window management, input handling
**donburi**: ECS framework

## Development Workflow

### Building

```bash
go build -o build/editor main.go
```

### Running

```bash
./build/editor
```

### Project Structure

```
GoEcs/
├── main.go                 # Editor entry point
├── go.mod                  # Dependencies
├── pkg/
│   ├── ecs/                # Layer 2: ECS World
│   │   ├── components.go   # Component definitions
│   │   ├── world.go        # Donburi wrapper
│   │   ├── scene.go        # Scene API (mutation layer)
│   │   └── systems.go      # TransformSystem
│   ├── editor/             # Layer 3: Editor
│   │   ├── camera.go       # EditorCamera
│   │   ├── viewport.go     # Viewport rendering
│   │   ├── state.go        # EditorState
│   │   ├── picking.go      # Raycasting
│   │   ├── gizmo.go        # Transform gizmos
│   │   ├── hierarchy_panel.go
│   │   └── inspector_panel.go
│   ├── renderer/           # Layer 1: Rendering
│   │   └── render_system.go
│   ├── assets/             # Asset management
│   │   └── registry.go
│   └── math/               # Math utilities
│       └── raycast.go      # Ray-AABB intersection
└── build/                  # Build output (gitignored)
```

## Future Enhancements

### Phase 9: Asset System
- Load meshes from .obj files
- Load textures from image files
- Asset browser panel

### Phase 10: Scene Serialization
- Save/load scenes to JSON
- Entity archetype system

### Phase 11: Logging & Console
- In-engine console for commands
- Log messages visible in editor

### Phase 12: Undo/Redo
- Command pattern for mutations
- History stack management

### Phase 13: More Gizmos
- Rotation gizmo (arc handles)
- Scale gizmo (box handles)
- Gizmo mode switching (Q/W/E)

### Phase 14: Play Mode Runtime
- Play/Pause/Stop buttons
- State saving/restoration
- Gameplay systems integration

### Phase 15: Polish
- Grid snapping
- Keyboard shortcuts
- Preference system
- Dark theme refinement

## Performance Considerations

### Entity Iteration

**Bad** (O(n) for every query):
```go
for _, entity := range world.GetAllEntities() {
    if hasTransform(entity) && hasRenderMesh(entity) {
        // ...
    }
}
```

**Good** (Donburi query, optimized internally):
```go
q := query.NewQuery(filter.Contains(TransformComponent, RenderMeshComponent))
q.Each(world, func(entry *donburi.Entry) {
    // Only iterates entities that match
})
```

### Transform Cache

- World matrices computed once per frame in TransformSystem
- Cached by entity ID
- RenderSystem and Gizmo system read from cache
- Invalidated each frame before Update

### Viewport Rendering

- RenderTexture prevents unnecessary redraws
- Only rendered when viewport is visible
- Future: Only redraw on scene changes

## Common Patterns

### Creating an Entity

```go
// Create entity
entity := scene.CreateEntity("MyObject")

// Add components
scene.AddTransform(entity)
scene.AddRenderMesh(entity, "cube", "default")

// Set initial transform
pos := rl.NewVector3(0, 0, 0)
scale := rl.NewVector3(1, 1, 1)
rot := rl.NewQuaternion(0, 0, 0, 1)
scene.SetTransform(entity, pos, scale, rot)
```

### Querying Entities

```go
// Get all entities
entities := scene.GetAllEntities()

// Donburi query (in systems)
q := query.NewQuery(filter.Contains(TransformComponent))
q.Each(world.GetDonburiWorld(), func(entry *donburi.Entry) {
    transform := TransformComponent.Get(entry)
    // Process transform
})
```

### Component Access

```go
// Safe access with error checking
if transform, ok := scene.GetTransform(entity); ok {
    newPos := rl.Vector3Add(transform.Position, delta)
    scene.SetTransform(entity, newPos, transform.Scale, transform.Rotation)
}
```

## Debugging Tips

### Entity Selection Issues

If selection doesn't work:
1. Check ray generation (viewport coordinates)
2. Verify positive distance check (behind camera filter)
3. Check AABB bounds for meshes
4. Print hit distances to console

### Gizmo Interaction Issues

If gizmo doesn't respond:
1. Check sensitivity value (0.05 recommended)
2. Verify axis hover detection
3. Check mouse capture logic
4. Ensure transform updates go through Scene API

### Rendering Issues

If entities don't render:
1. Verify mesh/material IDs exist in registries
2. Check TransformSystem.Update() is called
3. Verify RenderSystem.Render() is inside 3D mode
4. Check entity has both Transform and RenderMesh components

---

**Document Version**: 3.0
**Last Updated**: 2025-11-15
**Engine Version**: Phase 8 Complete (Translation Gizmo)
