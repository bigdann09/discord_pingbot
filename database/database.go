package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

type Server struct {
	ID       int    `gorm:"primaryKey"`
	Hostname string `gorm:"string;not null,index:idx_name,unique"`
}

// connectToDB establishes a connection to the PostgreSQL database using GORM.
func Connect(dsn string) (*Database, error) {
	if dsn == "" {
		return nil, gorm.ErrInvalidDB
	}

	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Server{}) // auto-migrate to database
	return &Database{DB: db}, err
}

// findServerByID retrieves a server by its ID from the database.
func (db Database) FindAllServers() ([]Server, error) {
	var servers []Server
	result := db.DB.Table("servers").Scan(&servers)
	if result.Error != nil {
		return nil, result.Error
	}
	return servers, nil
}
