package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type Target struct {
	Table      string
	UUIDColumn string
	FKs        []FK
}

// Many To One 관계 업데이트
// table로부터 column을 가져오고, 그 column으로 SELECT를 해서
// reference table로부터 reference column 값을 가져와서 table의 UUID columnName을 업데이트
// e.g.
// SELECT id, column names, fks FROM table_name
// for each row
// SELECT reference_column, uuid FROM reference_table WHERE id = row.id
// UPDATE table_name SET uuid_column = uuid WHERE id = row
type FK struct {
	Column     string // FK 컬럼 이름
	UUIDColumn string // UUID 버전의 FK 컬럼 이름

	ReferencedTable  string
	ReferencedColumn string
}

var Targets = []Target{
	{"users", "uuid", []FK{
		{Column: "profile_image_id", UUIDColumn: "profile_image_uuid", ReferencedTable: "media", ReferencedColumn: "uuid"},
	}},
	{"media", "uuid", []FK{}},
	{"breeds", "uuid", []FK{}},
	{"pets", "uuid", []FK{
		{Column: "owner_id", UUIDColumn: "owner_uuid", ReferencedTable: "users", ReferencedColumn: "uuid"},
		{Column: "profile_image_id", UUIDColumn: "profile_image_uuid", ReferencedTable: "media", ReferencedColumn: "uuid"},
	}},
	{"base_posts", "uuid", []FK{
		{Column: "author_id", UUIDColumn: "author_uuid", ReferencedTable: "users", ReferencedColumn: "uuid"},
	}},
	{"sos_posts", "uuid", []FK{
		{Column: "thumbnail_id", UUIDColumn: "thumbnail_uuid", ReferencedTable: "media", ReferencedColumn: "uuid"},
	}},
	{"sos_dates", "uuid", []FK{}},
	{"sos_posts_dates", "uuid", []FK{
		{Column: "sos_post_id", UUIDColumn: "sos_post_uuid", ReferencedTable: "sos_posts", ReferencedColumn: "uuid"},
		{Column: "sos_dates_id", UUIDColumn: "sos_dates_uuid", ReferencedTable: "sos_dates", ReferencedColumn: "uuid"},
	}},
	{"sos_conditions", "uuid", []FK{}},
	{"sos_posts_conditions", "uuid", []FK{
		{Column: "sos_post_id", UUIDColumn: "sos_post_uuid", ReferencedTable: "sos_posts", ReferencedColumn: "uuid"},
		{
			Column: "sos_condition_id", UUIDColumn: "sos_condition_uuid", ReferencedTable: "sos_conditions",
			ReferencedColumn: "uuid",
		},
	}},
	{"sos_posts_pets", "uuid", []FK{
		{Column: "sos_post_id", UUIDColumn: "sos_post_uuid", ReferencedTable: "sos_posts", ReferencedColumn: "uuid"},
		{Column: "pet_id", UUIDColumn: "pet_uuid", ReferencedTable: "pets", ReferencedColumn: "uuid"},
	}},
	{
		"resource_media", "uuid", []FK{
			{Column: "media_id", UUIDColumn: "media_uuid", ReferencedTable: "media", ReferencedColumn: "uuid"},
			// Conditional FK -- Thankfully, we only have SOS post for resource
			{Column: "resource_id", UUIDColumn: "resource_uuid", ReferencedTable: "sos_posts", ReferencedColumn: "uuid"},
		},
	},
}

type MigrateOptions struct {
	ReadOnly bool
	Force    bool
	Log      bool
}

func migrate(ctx context.Context, db *database.DB, options MigrateOptions) *pnd.AppError {
	log.Println("Migrating UUIDs for tables")
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}
	// UUID 컬럼만 먼저 업데이트
	err = MigrateUUID(tx, options)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println()

	log.Println("Migrating FKs for tables")
	// Read fk constraints and update them as well
	tx, err = db.BeginTx(ctx)
	if err != nil {
		return err
	}
	err = MigrateFK(tx, options)
	if err != nil {
		return err
	}

	return nil
}

func MigrateUUID(tx *database.Tx, options MigrateOptions) *pnd.AppError {
	type Row struct {
		ID   int
		UUID *string
	}

	total := 0
	succeeded := 0
	skipped := 0

	for _, target := range Targets {
		log.Printf("Processing table %s\n", target.Table)

		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM " + target.Table).Scan(&count)
		if err != nil {
			return pnd.ErrUnknown(fmt.Errorf("error counting rows from table %s due to %w", target.Table, err))
		}
		log.Printf(" - Total rows: %d\n", count)
		total += count

		rows, err := tx.Query("SELECT id, uuid FROM " + target.Table + " ORDER BY created_at ASC")
		if err != nil {
			return pnd.ErrUnknown(
				fmt.Errorf("error selecting rows from table %s due to %w", target.Table, err),
			)
		}

		var rowData []Row
		for rows.Next() {
			var row Row
			err = rows.Scan(&row.ID, &row.UUID)
			if err != nil {
				return pnd.ErrUnknown(
					fmt.Errorf("error scanning row from table %s with id %d due to %w", target.Table, row.ID, err),
				)
			}
			rowData = append(rowData, row)
		}
		rows.Close()

		for _, row := range rowData {
			if !options.Force && row.UUID != nil {
				skipped++
				log.Printf(" - Skipping row %d with UUID %s\n", row.ID, *row.UUID)
				continue
			}

			newUUID, err := uuid.NewV7()
			if err != nil {
				return pnd.ErrUnknown(fmt.Errorf("error generating UUID due to %w", err))
			}

			if !options.ReadOnly {
				_, err = tx.Exec("UPDATE "+target.Table+" SET "+target.UUIDColumn+" = $1 WHERE id = $2", newUUID, row.ID)
				if err != nil {
					return pnd.ErrUnknown(
						fmt.Errorf("error updating UUID column in table %s with id %d due to %w", target.Table, row.ID, err),
					)
				}
				log.Printf(" - Updated row %d with UUID %s\n", row.ID, newUUID)
			}

			succeeded++
		}
	}

	log.Printf("Total rows: %d, Succeeded: %d, Skipped: %d\n", total, succeeded, skipped)

	return nil
}

// FK 컬럼을 연관된 테이블의 uuid 컬럼을 조회해 업데이트
func MigrateFK(tx *database.Tx, options MigrateOptions) *pnd.AppError {
	for _, target := range Targets {
		log.Printf("Processing table %s\n", target.Table)

		for _, fk := range target.FKs {
			succeeded := 0
			skipped := 0

			rows, err := tx.Query("SELECT id, " + fk.Column + " FROM " + target.Table)
			if err != nil {
				return pnd.ErrUnknown(
					fmt.Errorf("error selecting rows from table %s due to %w", target.Table, err),
				)
			}

			type RowData struct {
				ID      int
				FKValue *int
			}

			var rowData []RowData
			for rows.Next() {
				var row struct{ ID int }
				var fkValue *int

				err = rows.Scan(&row.ID, &fkValue)
				if err != nil {
					return pnd.ErrUnknown(
						fmt.Errorf("error scanning row from table %s with id %d due to %w", target.Table, row.ID, err),
					)
				}

				rowData = append(rowData, RowData{row.ID, fkValue})
			}
			rows.Close()

			total := len(rowData)
			log.Printf(" - FK %s, Total rows: %d\n", fk.Column, len(rowData))

			for _, row := range rowData {
				// Skip if FK is not set
				if row.FKValue == nil {
					log.Printf(" - Skipping row %d with FK %s due to null value\n", row.ID, fk.Column)
					skipped++
					continue
				}

				var uuidValue *string
				err = tx.QueryRow(
					"SELECT "+fk.ReferencedColumn+" FROM "+fk.ReferencedTable+" WHERE id = $1",
					row.FKValue).Scan(&uuidValue)
				if err != nil && err.Error() != "sql: no rows in result set" {
					return pnd.ErrUnknown(
						fmt.Errorf("error selecting UUID from table %s with id %d due to %w", fk.ReferencedTable, row.FKValue, err),
					)
				}

				// Skip if UUID is not set
				if uuidValue == nil {
					// Perhaps you forgot to run MigrateUUID() before MigrateFK()
					log.Printf(" - Skipping row %d with FK %s due to missing UUID\n", row.ID, fk.Column)
					skipped++
					continue
				}

				if !options.ReadOnly {
					_, err = tx.Exec("UPDATE "+target.Table+" SET "+fk.UUIDColumn+" = $1 WHERE id = $2", uuidValue, row.ID)
					if err != nil {
						return pnd.ErrUnknown(
							fmt.Errorf("error updating UUID column in table %s with id %d due to %w", target.Table, row.ID, err),
						)
					}
				}

				succeeded++
			}

			log.Printf(" - FK %s, Total rows: %d, Succeeded: %d, Skipped: %d\n", fk.Column, total, succeeded, skipped)
		}
	}

	return nil
}
