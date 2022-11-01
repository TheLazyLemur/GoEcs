package main

import (
	cmp "GoEcs/components"
	ent "GoEcs/entities"
	sys "GoEcs/systems"

	"github.com/EngoEngine/ecs"
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	w := ecs.World{}

	var enemyable *sys.Enemyable
	var playerable *sys.Playerable

	w.AddSystemInterface(&sys.PlayerSystem{
		Entities: make(map[uint64]ent.PlayerEntity),
	}, playerable, nil)

	w.AddSystemInterface(&sys.EnemySystem{
		Entities: make(map[uint64]ent.EnemyEntity),
	}, enemyable, nil)


	player := ent.PlayerEntity{
		BasicEntity: ecs.NewBasic(),
		MovementComponent: &cmp.MovementComponent{
			Speed: 100,
		},
		HealthComponent: &cmp.HealthComponent{
			Max:     100,
			Current: 100,
		},
        SpaceComponent: cmp.NewSpaceComponent(),
	}

	entity := ent.EnemyEntity{
		BasicEntity: ecs.NewBasic(),
		AiComponent: &cmp.AiComponent{
			Speed: 100,
		},
		HealthComponent: &cmp.HealthComponent{
			Max:     100,
			Current: 100,
		},
	}

	entity2 := ent.EnemyEntity{
		BasicEntity: ecs.NewBasic(),
		AiComponent: &cmp.AiComponent{
			Speed: 100,
		},
		HealthComponent: &cmp.HealthComponent{
			Max:     200,
			Current: 200,
		},
	}

	w.SortSystems()

	w.AddEntity(&player)
	w.AddEntity(&entity)
	w.AddEntity(&entity2)

	rl.InitWindow(800, 450, "Space Impact")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		w.Update(rl.GetFrameTime())
		rl.ClearBackground(rl.Black)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
