package editor

import (
	"github.com/TheLazyLemur/SpaceImpact/pkg/ecs"
)

// EditorMode represents the current mode of the editor
type EditorMode int

const (
	ModeEdit EditorMode = iota
	ModePlay
	ModePause
)

// String returns the string representation of the mode
func (m EditorMode) String() string {
	switch m {
	case ModeEdit:
		return "EDIT"
	case ModePlay:
		return "PLAY"
	case ModePause:
		return "PAUSE"
	default:
		return "UNKNOWN"
	}
}

// EditorState holds all editor-specific state
type EditorState struct {
	Mode      EditorMode
	Scene     *ecs.Scene
	Selected  []ecs.EntityID
	Viewports []*Viewport

	// Undo/Redo stacks (placeholder for now)
	UndoStack []interface{}
	RedoStack []interface{}

	// UI state
	HierarchyWidth  float32
	InspectorWidth  float32
	ShowHierarchy   bool
	ShowInspector   bool
}

// NewEditorState creates a new editor state
func NewEditorState(scene *ecs.Scene) *EditorState {
	return &EditorState{
		Mode:            ModeEdit,
		Scene:           scene,
		Selected:        make([]ecs.EntityID, 0),
		Viewports:       make([]*Viewport, 0),
		UndoStack:       make([]interface{}, 0),
		RedoStack:       make([]interface{}, 0),
		HierarchyWidth:  250.0,
		InspectorWidth:  300.0,
		ShowHierarchy:   true,
		ShowInspector:   true,
	}
}

// SelectEntity sets the selected entity (single selection)
func (es *EditorState) SelectEntity(id ecs.EntityID) {
	es.Selected = []ecs.EntityID{id}
}

// AddToSelection adds an entity to the selection (multi-select)
func (es *EditorState) AddToSelection(id ecs.EntityID) {
	// Check if already selected
	for _, selectedID := range es.Selected {
		if selectedID == id {
			return
		}
	}
	es.Selected = append(es.Selected, id)
}

// ClearSelection clears all selected entities
func (es *EditorState) ClearSelection() {
	es.Selected = make([]ecs.EntityID, 0)
}

// IsSelected checks if an entity is currently selected
func (es *EditorState) IsSelected(id ecs.EntityID) bool {
	for _, selectedID := range es.Selected {
		if selectedID == id {
			return true
		}
	}
	return false
}

// GetFirstSelected returns the first selected entity (or nil if none)
func (es *EditorState) GetFirstSelected() ecs.EntityID {
	if len(es.Selected) > 0 {
		return es.Selected[0]
	}
	return nil
}

// HasSelection returns true if any entities are selected
func (es *EditorState) HasSelection() bool {
	return len(es.Selected) > 0
}

// SetMode changes the editor mode
func (es *EditorState) SetMode(mode EditorMode) {
	es.Mode = mode
}
