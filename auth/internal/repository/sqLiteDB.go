package repository

import "gorm.io/gorm"

type GormSQLiteDB struct {
	db *gorm.DB
}
