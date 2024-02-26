package main

import (
	"errors"
	"flag"
	"fmt"
	//Библиотека для миграций
	migrator "github.com/golang-migrate/migrate/v4"
	//Драйвер для выполнения миграций
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	//Драйвер для получения файлов миграций
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string
	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is empty")
	}
	if migrationsPath == "" {
		panic("migrations path is empty")
	}
	m, err := migrator.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%v?x-migrations-table=%v", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrator.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied")
}
