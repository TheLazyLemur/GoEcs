package assets

import (
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// MeshRegistry stores meshes by ID
type MeshRegistry struct {
	meshes map[string]rl.Mesh
	mu     sync.RWMutex
}

// NewMeshRegistry creates a new mesh registry
func NewMeshRegistry() *MeshRegistry {
	return &MeshRegistry{
		meshes: make(map[string]rl.Mesh),
	}
}

// Register adds a mesh to the registry
func (mr *MeshRegistry) Register(id string, mesh rl.Mesh) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.meshes[id] = mesh
}

// Get retrieves a mesh by ID
func (mr *MeshRegistry) Get(id string) (rl.Mesh, bool) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()
	mesh, ok := mr.meshes[id]
	return mesh, ok
}

// MaterialRegistry stores materials by ID
type MaterialRegistry struct {
	materials map[string]rl.Material
	mu        sync.RWMutex
}

// NewMaterialRegistry creates a new material registry
func NewMaterialRegistry() *MaterialRegistry {
	return &MaterialRegistry{
		materials: make(map[string]rl.Material),
	}
}

// Register adds a material to the registry
func (mr *MaterialRegistry) Register(id string, material rl.Material) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.materials[id] = material
}

// Get retrieves a material by ID
func (mr *MaterialRegistry) Get(id string) (rl.Material, bool) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()
	material, ok := mr.materials[id]
	return material, ok
}

// AssetManager holds all asset registries
type AssetManager struct {
	meshes    *MeshRegistry
	materials *MaterialRegistry
}

// NewAssetManager creates a new asset manager
func NewAssetManager() *AssetManager {
	return &AssetManager{
		meshes:    NewMeshRegistry(),
		materials: NewMaterialRegistry(),
	}
}

// GetMeshes returns the mesh registry
func (am *AssetManager) GetMeshes() *MeshRegistry {
	return am.meshes
}

// GetMaterials returns the material registry
func (am *AssetManager) GetMaterials() *MaterialRegistry {
	return am.materials
}

// CreateDefaultAssets creates some basic hardcoded assets for testing
func (am *AssetManager) CreateDefaultAssets() {
	fmt.Println("Creating default assets...")

	// Create default meshes
	cubeMesh := rl.GenMeshCube(1.0, 1.0, 1.0)
	sphereMesh := rl.GenMeshSphere(0.5, 16, 16)
	planeMesh := rl.GenMeshPlane(10.0, 10.0, 1, 1)

	am.meshes.Register("cube", cubeMesh)
	am.meshes.Register("sphere", sphereMesh)
	am.meshes.Register("plane", planeMesh)

	// Create default materials
	defaultMat := rl.LoadMaterialDefault()

	redMat := rl.LoadMaterialDefault()
	redImg := rl.GenImageColor(1, 1, rl.Red)
	redTex := rl.LoadTextureFromImage(redImg)
	rl.SetMaterialTexture(&redMat, rl.MapDiffuse, redTex)
	rl.UnloadImage(redImg)

	greenMat := rl.LoadMaterialDefault()
	greenImg := rl.GenImageColor(1, 1, rl.Green)
	greenTex := rl.LoadTextureFromImage(greenImg)
	rl.SetMaterialTexture(&greenMat, rl.MapDiffuse, greenTex)
	rl.UnloadImage(greenImg)

	blueMat := rl.LoadMaterialDefault()
	blueImg := rl.GenImageColor(1, 1, rl.Blue)
	blueTex := rl.LoadTextureFromImage(blueImg)
	rl.SetMaterialTexture(&blueMat, rl.MapDiffuse, blueTex)
	rl.UnloadImage(blueImg)

	am.materials.Register("default", defaultMat)
	am.materials.Register("red", redMat)
	am.materials.Register("green", greenMat)
	am.materials.Register("blue", blueMat)

	fmt.Println("Default assets created!")
}
