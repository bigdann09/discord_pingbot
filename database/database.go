package database

import "gorm.io/gorm"

type Database struct {
	DB *gorm.DB
}

type Server struct {
	ID       int    `gorm:"primaryKey"`
	Hostname string `gorm:"string;not null,index:idx_name,unique"`
}
