package gorm_app

import (
	"fmt"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/config"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Phase 2: Database connection setup (GORM initialization, DB-from-context pattern). See resources/phase-02-database-di.md

func Open(cfg *config.Config) (*gorm.DB, func(), error) {
	db, err := gorm.Open(mysqlDriver.Open(buildDSN(cfg.DB)), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return db, cleanup, nil
}

func buildDSN(cfg *config.DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}
