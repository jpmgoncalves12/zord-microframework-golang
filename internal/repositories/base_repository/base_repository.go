package base_repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
)

type BaseRepository[dom Domain] interface {
	Get(field string, value string) (*dom, error)
	Create(data dom) error
	List(limit int, offset int) (*[]dom, error)
	Edit(data dom, field string, value string) (int, error)
	Delete(field string, values string) error
	Count() (int64, error)
}

type Domain interface {
	Schema() string
}

type BaseRepo[dom Domain] struct {
	Mysql  *sqlx.DB
	fields []string
}

func NewBaseRepository[dom Domain](mysql *sqlx.DB) *BaseRepo[dom] {
	do := new(dom)
	listFields := []string{}
	fields := structs.Fields(&do)
	for _, field := range fields {
		tag := field.Tag("db")
		listFields = append(listFields, tag)
	}

	return &BaseRepo[dom]{
		Mysql:  mysql,
		fields: listFields,
	}
}

func (repo *BaseRepo[Domain]) Commit(tx *sqlx.Tx) error {
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo *BaseRepo[Domain]) Rollback(tx *sqlx.Tx, err error) error {
	rollbErr := tx.Rollback()
	if rollbErr != nil {
		return rollbErr
	}
	return err
}

func (repo *BaseRepo[Domain]) Get(field string, value string) (*Domain, error) {
	var Data Domain
	row := repo.Mysql.QueryRowx(
		fmt.Sprintf(
			"SELECT %s FROM %s WHERE %s = ?",
			strings.Join(repo.fields, ", "),
			Data.Schema(),
			field,
		),
		value,
	)
	err := row.StructScan(&Data)
	if err != nil {
		return nil, err
	}
	return &Data, nil
}

func (repo *BaseRepo[Row]) Create(d Row) error {
	tx, err := repo.Mysql.Beginx()
	if err != nil {
		return err
	}

	namedValues := []string{}
	for _, field := range repo.fields {
		namedValues = append(namedValues, ":"+field)
	}
	res, err := tx.NamedExec(
		fmt.Sprintf(
			`INSERT INTO %s (%s) VALUES (%s)`,
			d.Schema(),
			strings.Join(repo.fields, ", "),
			strings.Join(namedValues, ", "),
		),
		d,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if affected < 1 {
		return errors.New("error on create, rows not affected")
	}
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo *BaseRepo[Row]) List(limit int, offset int) (*[]Row, error) {
	var data []Row
	var value Row
	rows, err := repo.Mysql.Queryx(
		fmt.Sprintf(
			`SELECT %s FROM %s LIMIT %v OFFSET %v`,
			strings.Join(repo.fields, ", "),
			value.Schema(),
			limit,
			offset,
		),
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.StructScan(&value)
		if err != nil {
			return nil, err
		}
		data = append(data, value)
	}
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *BaseRepo[Row]) Edit(d Row, field string, value string) (int, error) {
	tx, err := repo.Mysql.Beginx()
	if err != nil {
		return 0, err
	}

	namedValues := []string{}
	for _, field := range repo.fields {
		namedValues = append(namedValues, "`"+field+"`"+" = :"+field)
	}

	query, args, bindErr := tx.BindNamed(
		fmt.Sprintf(
			`UPDATE %s SET %s WHERE %s = '%s'`,
			d.Schema(),
			strings.Join(namedValues, ", "),
			field,
			value,
		),
		&d,
	)

	if bindErr != nil {
		return 0, repo.Rollback(tx, bindErr)
	}

	res, execErr := repo.Mysql.Exec(query, args...)
	if execErr != nil {
		return 0, repo.Rollback(tx, execErr)
	}

	affected, rowsAffErr := res.RowsAffected()
	if rowsAffErr != nil {
		return int(affected), repo.Rollback(tx, rowsAffErr)
	}

	if affected < 1 {
		repo.Rollback(tx, nil)
		return int(affected), nil
	}

	commitErr := repo.Commit(tx)
	if commitErr != nil {
		return int(affected), repo.Rollback(tx, commitErr)
	}

	return int(affected), nil
}

func (repo *BaseRepo[Row]) Delete(field string, value string) error {
	var data Row
	exec, err := repo.Mysql.Exec(
		fmt.Sprintf(
			`DELETE FROM %s WHERE %s = '%s'`,
			data.Schema(),
			field,
			value,
		),
	)
	if err != nil {
		return err
	}
	rowsAffected, err := exec.RowsAffected()
	if rowsAffected < 1 {
		return errors.New("nothing deleted")
	}
	return err
}

func (repo *BaseRepo[Row]) Count() (int64, error) {
	var count int64
	var data Row
	err := repo.Mysql.Get(&count, "SELECT count(1) FROM "+data.Schema())
	return count, err
}
