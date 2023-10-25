package dummy

import (
	"go-skeleton/pkg/database"
)

type migration struct {
	db *database.MySql
}

func NewMigration(db *database.MySql) *migration {
	return &migration{
		db: db,
	}
}

func (db *migration) Migrate() {
	db.db.Db.Migrator().CreateTable(&Dummy{})
}

func (db *migration) RollBack() {
	db.db.Db.Migrator().DropTable(&Dummy{})
}
