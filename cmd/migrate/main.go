package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/7Maliko7/forms-api/migration"
)

func main() {
	pgMigrationsDir := "database/postgres"
	pgConnString := os.Args[2]

	fmt.Println("Found migration files")
	fmt.Println("Postgres:")

	fileList, err := getAllFilenames(&migration.Database, pgMigrationsDir)
	if err != nil {
		log.Fatalln(err)
	}
	for i := range fileList {
		fmt.Println(fileList[i])
	}

	m, err := migrate.New("file://migration/"+pgMigrationsDir, pgConnString)
	if err != nil {
		log.Fatalln(err)
	}

	ver, dirt, _ := m.Version()
	fmt.Printf("DB Version: %v Dirty: %v\n", ver, dirt)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			migrateUp(m)
		case "down":
			migrateDown(m)
		}
	}

	ver, dirt, _ = m.Version()
	fmt.Printf("New DB Version: %v Dirty: %v\n", ver, dirt)
}

func getAllFilenames(fs *embed.FS, dir string) (out []string, err error) {
	if len(dir) == 0 {
		dir = "."
	}
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		fp := path.Join(dir, entry.Name())
		if entry.IsDir() {
			res, err := getAllFilenames(fs, fp)
			if err != nil {
				return nil, err
			}
			out = append(out, res...)
			continue
		}
		out = append(out, fp)
	}
	return
}

func migrateUp(m *migrate.Migrate) {
	log.Println("Migrate up started")

	if err := m.Up(); err != nil {
		log.Println(err)
	}

	log.Println("Migrate up finished")
}

func migrateDown(m *migrate.Migrate) {
	log.Println("Migrate down started")

	if err := m.Down(); err != nil {
		log.Println(err)
	}

	log.Println("Migrate down finished")
}
