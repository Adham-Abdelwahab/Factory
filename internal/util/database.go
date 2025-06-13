package util

import (
	"context"
	"os"
	"reflect"
	"unsafe"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type database struct {
	Ctx  context.Context
	Conn *pgxpool.Pool
}

var Database database

func InitializeDatabase() {
	ctx := context.Background()
	db := os.Getenv("DATABASE_URL")
	if db, err := pgxpool.New(ctx, db); err != nil {
		panic("failed to connect to the database")
	} else {
		Database = database{ctx, db}
	}
}

func scanner(structure any) []any {
	var pointers []any

	value := reflect.ValueOf(structure)
	pointer := unsafe.Pointer(value.Pointer())
	fields := reflect.TypeOf(value.Elem().Interface())

	if value.Elem().Kind() != reflect.Struct {
		return []any{structure}
	}

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		offset := unsafe.Add(pointer, field.Offset)

		switch field.Type.String() {
		case "int":
			pointers = append(pointers, (*int)(offset))
		case "string":
			pointers = append(pointers, (*string)(offset))
		case "bool":
			pointers = append(pointers, (*bool)(offset))
		}
	}

	return pointers
}

func (db *database) Query(sql string, args ...any) (pgx.Rows, error) {
	return db.Conn.Query(db.Ctx, sql, args...)
}

func (*database) ForEach(rows pgx.Rows, scan any, method func() error) error {
	_, err := pgx.ForEachRow(rows, scanner(scan), method)
	return err
}

func (db *database) QueryRow(row any, sql string, args ...any) error {
	return db.Conn.QueryRow(db.Ctx, sql, args...).Scan(row)
}

func (db *database) Close() {
	db.Ctx.Done()
	db.Conn.Close()
}
