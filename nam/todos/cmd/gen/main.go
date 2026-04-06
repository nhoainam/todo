package main

import (
	"gorm.io/gen"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/infra/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(
		entity.Todo{},
	)

	g.Execute()
}
