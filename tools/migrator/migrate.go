package migrator

import (
	"go-skeleton/application/domain/dummy"

	//{{codeGen3}}
	"go-skeleton/pkg/database"
)

type migrations struct {
	Id    int    `gorm:"primarykey" json:"id"`
	Name  string `validate:"required" json:"name"`
	Batch int    `gorm:"primarykey" json:"batch"`
}

type Migrator struct {
	db *database.MySql
}

func NewMigrator(db *database.MySql) *Migrator {
	return &Migrator{
		db: db,
	}
}

func (m *Migrator) Migrate() {
	m.db.Db.AutoMigrate(&migrations{})

	dummy.NewMigration(m.db).Migrate()
	//{{codeGen4}}
}

func (m *Migrator) RollBack() {
	m.db.Db.AutoMigrate(&migrations{})

	dummy.NewMigration(m.db).RollBack()
	//{{codeGen5}}
}
