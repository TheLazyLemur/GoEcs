# GoEcs Engine + Editor Development Roadmap

## Overview

This roadmap outlines the phased development approach for building a Go-based game engine with an integrated editor from scratch. Each phase builds upon the previous, culminating in a fully-featured editor with ECS-based game world.

**Current Status**: Phase 8 Complete ✅
**Next Phase**: Phase 9 (Asset System)

---

## ✅ Phase 0: Project Foundation

**Status**: Complete
**Commit**: Initial setup

### Objectives
- Set up Go project structure
- Initialize go.mod with dependencies
- Create package hierarchy

### Deliverables
- [x] Project structure (`pkg/ecs`, `pkg/editor`, `pkg/renderer`, `pkg/assets`, `pkg/math`)
- [x] go.mod with raylib-go v0.55.1
- [x] Basic main.go entry point
- [x] .gitignore configuration

---

## ✅ Phase 1-2: Core ECS Foundation

**Status**: Complete (Later refactored to Donburi)
**Commits**: Initial ECS, Donburi refactor

### Objectives
- Implement Donburi ECS wrapper
- Define core component types
- Create Scene API wrapper
- Build basic systems

### Deliverables
- [x] `pkg/ecs/world.go` - Donburi World wrapper
- [x] `pkg/ecs/components.go` - Transform, RenderMesh, Name, Camera components
- [x] `pkg/ecs/scene.go` - Scene API (single mutation entry point)
- [x] `pkg/ecs/systems.go` - TransformSystem with matrix caching
- [x] Integration with Donburi v1.15.7

### Technical Notes
- Originally implemented custom ECS
- Refactored to Donburi ECS per spec requirements
- EntityID is `*donburi.Entry`
- Components use `ComponentType[T]` for type safety
- Scene API isolates editor from direct World access

---

## ✅ Phase 3: Editor Camera & Viewport

**Status**: Complete
**Commit**: "feat: Phase 3 - Editor camera and viewport"

### Objectives
- FPS-style editor camera
- RenderTexture-based viewport
- Viewport bounds and screen-space mapping

### Deliverables
- [x] `pkg/editor/camera.go` - EditorCamera with WASD + mouse look
- [x] `pkg/editor/viewport.go` - Viewport with RenderTexture pipeline
- [x] Camera controls (WASD, Q/E, right-click look, scroll speed)
- [x] Viewport rendering integration in main loop

### Technical Notes
- EditorCamera is NOT an ECS component (editor layer)
- RenderTexture allows post-processing and multi-viewport support
- Viewport-relative coordinates used for picking

### Controls
```
WASD     - Move horizontally
Q/E      - Move up/down
RMB Drag - Look around
Scroll   - Adjust move speed
```

---

## ✅ Phase 4-6: Editor UI Panels

**Status**: Complete
**Commit**: "feat: Phase 4-6 - Editor state and UI panels"

### Objectives
- Central editor state management
- Hierarchy panel (entity list)
- Inspector panel (component editor)
- Selection system

### Deliverables
- [x] `pkg/editor/state.go` - EditorState (mode, selection, viewports)
- [x] `pkg/editor/hierarchy_panel.go` - Entity list with selection
- [x] `pkg/editor/inspector_panel.go` - Component display
- [x] Selection management (single/multi-select)
- [x] Panel drawing with custom raylib primitives

### Technical Notes
- Current UI uses custom drawing (not raygui yet)
- Selection supports Ctrl+Click for multi-select
- Inspector displays Name, Transform, RenderMesh components
- Panels use scissor mode for scrolling

### Known Limitations
- Inspector is read-only (no editing yet)
- No raygui integration yet (planned for later)

---

## ✅ Phase 7: Viewport Picking & Selection

**Status**: Complete
**Commit**: "feat: Phase 7 - Viewport picking and selection"

### Objectives
- Raycasting for entity selection
- Viewport-relative mouse picking
- Selection highlights

### Deliverables
- [x] `pkg/math/raycast.go` - Ray-AABB intersection, ray generation
- [x] `pkg/editor/picking.go` - Viewport picking system
- [x] Selection highlights rendering
- [x] Click-to-select in viewport

### Technical Notes
- Ray generated from viewport-relative mouse position
- Positive distance check filters objects behind camera
- Uses mesh AABB bounds (transformed to world space)
- Integration with hierarchy panel selection

### Bug Fixes
- Fixed ray generation (correct projection matrix)
- Fixed behind-camera filtering (distance > 0 check)
- Fixed AABB transformation

---

## ✅ Phase 8: Translation Gizmo

**Status**: Complete
**Commit**: "feat: Phase 8 - Translation gizmo"

### Objectives
- 3D manipulation handles for entity movement
- Axis-constrained dragging
- Visual feedback for active axis

### Deliverables
- [x] `pkg/editor/gizmo.go` - Gizmo system with translation mode
- [x] Three-axis handles (X=Red, Y=Green, Z=Blue)
- [x] Hover detection via raycasting
- [x] Click-and-drag interaction
- [x] Transform updates via Scene API

### Technical Notes
- Sensitivity: 0.05 units per pixel
- Axis constraint: Maps 2D mouse delta to 3D axis movement
- Visual highlight for active axis
- Integrated with viewport rendering

### Bug Fixes
- Increased sensitivity from 0.01 to 0.05
- Fixed Z-axis direction for intuitive forward/back

---

## 🔄 Phase 9: Asset System

**Status**: Pending
**Next Up**

### Objectives
- Load meshes from files (.obj, .gltf)
- Load textures from images (.png, .jpg)
- Asset browser panel
- Asset registry improvements

### Planned Deliverables
- [ ] `pkg/assets/mesh_loader.go` - Mesh file loading
- [ ] `pkg/assets/texture_loader.go` - Texture file loading
- [ ] `pkg/editor/asset_browser.go` - Asset browser panel
- [ ] File watching for hot-reload
- [ ] Asset metadata system

### Technical Plan
- Use raylib's built-in LoadModel/LoadTexture
- Asset browser shows thumbnail grid
- Drag-and-drop from browser to viewport
- Asset path resolution (relative to project root)

### Success Criteria
- Load .obj mesh and see it in viewport
- Load .png texture and apply to material
- Browse assets in dedicated panel

---

## 📋 Phase 10: Scene Serialization

**Status**: Pending

### Objectives
- Save scenes to JSON files
- Load scenes from JSON files
- Scene file management

### Planned Deliverables
- [ ] `pkg/ecs/serialization.go` - Scene save/load
- [ ] JSON schema for scene format
- [ ] File → Save/Load menu (or shortcuts)
- [ ] Entity archetype system

### Technical Plan
```json
{
  "name": "MyScene",
  "entities": [
    {
      "name": "Cube",
      "components": {
        "Transform": { "position": [0, 0, 0], "rotation": [0, 0, 0, 1], "scale": [1, 1, 1] },
        "RenderMesh": { "mesh": "cube", "material": "default" }
      }
    }
  ]
}
```

### Success Criteria
- Save scene with multiple entities
- Load scene and see entities restored
- Handle missing assets gracefully

---

## 📋 Phase 11: Logging System & Console

**Status**: Pending

### Objectives
- In-engine console for commands
- Log message system
- Console panel in editor

### Planned Deliverables
- [ ] `pkg/logging/logger.go` - Logging system
- [ ] `pkg/editor/console_panel.go` - Console UI
- [ ] Log levels (Debug, Info, Warn, Error)
- [ ] Console commands (create entity, delete entity, etc.)

### Technical Plan
- Ring buffer for log messages
- Console panel with input field
- Command parser with autocomplete
- Scrollable log output

### Example Commands
```
create cube
delete entity_5
set cube.position 1 2 3
help
clear
```

---

## 📋 Phase 12: Undo/Redo System

**Status**: Pending

### Objectives
- Command pattern for all mutations
- Undo/Redo stack
- Keyboard shortcuts (Ctrl+Z, Ctrl+Y)

### Planned Deliverables
- [ ] `pkg/editor/command.go` - Command interface
- [ ] `pkg/editor/command_history.go` - History stack
- [ ] Commands for all Scene API operations
- [ ] Undo/Redo UI indicators

### Technical Plan
```go
type Command interface {
    Execute(scene *Scene) error
    Undo(scene *Scene) error
    Description() string
}

type SetTransformCommand struct {
    entityID    EntityID
    oldTransform Transform
    newTransform Transform
}
```

### Success Criteria
- Move entity, undo, see it return to original position
- Redo after undo
- Undo stack persists across multiple operations

---

## 📋 Phase 13: Rotation & Scale Gizmos

**Status**: Pending

### Objectives
- Rotation gizmo with arc handles
- Scale gizmo with box handles
- Gizmo mode switching (Q/W/E)

### Planned Deliverables
- [ ] Rotation gizmo (three arc handles)
- [ ] Scale gizmo (six box handles + center for uniform scale)
- [ ] Gizmo mode switching UI
- [ ] Keyboard shortcuts (Q=Select, W=Translate, E=Rotate, R=Scale)

### Technical Plan

**Rotation Gizmo**:
- Three colored arcs (X=Red, Y=Green, Z=Blue)
- Drag to rotate around axis
- Display angle during rotation

**Scale Gizmo**:
- Six boxes on axis endpoints
- Center box for uniform scaling
- Drag to scale along axis

### Success Criteria
- Rotate entity using rotation gizmo
- Scale entity using scale gizmo
- Switch between gizmos with keyboard

---

## 📋 Phase 14: Play Mode Runtime

**Status**: Pending

### Objectives
- Play/Pause/Stop buttons
- State saving/restoration
- Separate edit mode from play mode
- Gameplay system integration

### Planned Deliverables
- [ ] Play mode state management
- [ ] Scene snapshot for restoration
- [ ] Play/Pause/Stop UI buttons
- [ ] Disable editor tools during play mode

### Technical Plan
```go
type EditorMode int
const (
    ModeEdit EditorMode = iota
    ModePlay
    ModePause
)
```

**Play Flow**:
1. User clicks Play
2. Save scene snapshot (deep copy)
3. Switch to ModePlay
4. Run gameplay systems
5. User clicks Stop
6. Restore scene from snapshot
7. Switch to ModeEdit

### Success Criteria
- Click Play, entities move/interact
- Click Stop, entities return to edit state
- Pause/Resume works correctly

---

## 📋 Phase 15: Polish & Enhancement

**Status**: Pending

### Objectives
- Quality-of-life improvements
- Performance optimizations
- User experience refinement

### Planned Deliverables
- [ ] Grid snapping (hold Ctrl while dragging)
- [ ] Keyboard shortcuts panel (F1 for help)
- [ ] Preference system (settings file)
- [ ] Recent files menu
- [ ] Dark theme refinement
- [ ] Performance profiling tools
- [ ] Entity duplication (Ctrl+D)
- [ ] Multi-viewport support
- [ ] Camera bookmarks
- [ ] Custom grid size

### Technical Improvements
- Spatial partitioning for large scenes
- Frustum culling
- LOD system
- Asset streaming
- Memory profiling

### Success Criteria
- Editor feels responsive with 1000+ entities
- Workflow is smooth and intuitive
- Settings persist across sessions

---

## Completed Milestones

### ✅ Milestone 1: Basic Editor (Phases 0-3)
- Project structure ✅
- ECS foundation ✅
- Editor camera ✅
- Viewport rendering ✅

### ✅ Milestone 2: Selection & Inspection (Phases 4-7)
- Editor state management ✅
- UI panels ✅
- Viewport picking ✅
- Selection system ✅

### ✅ Milestone 3: Transform Manipulation (Phase 8)
- Translation gizmo ✅
- Basic entity editing ✅

### ✅ Milestone 3.5: ECS Architecture Refactor
- Migrated from custom ECS to Donburi ✅
- Scene API isolation ✅
- Type-safe component access ✅

---

## Upcoming Milestones

### 📋 Milestone 4: Asset Pipeline (Phases 9-10)
- File loading
- Scene serialization
- Asset management

### 📋 Milestone 5: Editor Features (Phases 11-13)
- Logging & console
- Undo/Redo
- Full gizmo suite

### 📋 Milestone 6: Runtime & Polish (Phases 14-15)
- Play mode
- Final polish
- Release preparation

---

## Development Principles

### Code Quality
- ✅ All code must compile without warnings
- ✅ Commit after each phase completion
- ✅ Clear commit messages describing changes
- ✅ No commented-out code in commits

### Testing Strategy
- Manual testing after each phase
- Interactive testing in editor
- Test with multiple entities (stress test)
- Test edge cases (empty scenes, missing assets)

### Performance Targets
- 60 FPS with 100 entities
- 30 FPS with 1000 entities
- < 100ms startup time
- < 16ms frame time (edit mode)

---

## Version History

### v0.8.0 (Current)
- Phase 8 Complete: Translation Gizmo
- Donburi ECS integration
- Selection and picking
- Editor camera and viewport

### v0.3.0
- Phase 3 Complete: Editor camera

### v0.2.0
- Phase 1-2 Complete: ECS foundation (custom)

### v0.1.0
- Phase 0 Complete: Project structure

---

## Contributing Guidelines

### Branch Naming
- Feature branches: `claude/ecs-editor-roadmap-{sessionID}`
- Hotfix branches: `hotfix/description`

### Commit Format
```
<type>: <description>

<optional body>

<optional footer>
```

**Types**: feat, fix, refactor, docs, chore, test

### Phase Completion Checklist
- [ ] All planned deliverables implemented
- [ ] Code compiles without errors
- [ ] Manual testing performed
- [ ] Commit with descriptive message
- [ ] Push to remote branch
- [ ] Update this roadmap document

---

**Document Version**: 1.0
**Last Updated**: 2025-11-15
**Current Phase**: 8/15 Complete

