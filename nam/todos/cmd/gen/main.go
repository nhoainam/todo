package main

import (
	"gorm.io/gen"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/persistence/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/infra/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(
		model.Todo{},
	)

	g.Execute()
}
