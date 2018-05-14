package main

import (
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDB() (*gorm.DB, *gorm.DB, *gorm.DB) {
	blockdb, err := gorm.Open("sqlite3", filepath.Join(dbpath, "block.db"))
	if err != nil {
		panic(err)
	}

	accountdb, err := gorm.Open("sqlite3", filepath.Join(dbpath, "accounts.db"))
	if err != nil {
		panic(err)
	}

	witnessdb, err := gorm.Open("sqlite3", filepath.Join(dbpath, "witness.db"))
	if err != nil {
		panic(err)
	}
	// db.AutoMigrate(&SyncInfo{})

	return blockdb, accountdb, witnessdb
}