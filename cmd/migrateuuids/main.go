package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/rs/zerolog/log"
)

var targetTables = map[string][]string{
	"users":                {"uuid", "profile_image_uuid"},
	"media":                {"uuid"},
	"breeds":               {"uuid"},
	"pets":                 {"uuid", "owner_uuid", "profile_image_uuid"},
	"base_posts":           {"uuid", "author_uuid"},
	"sos_posts":            {"thumbnail_uuid"},
	"sos_dates":            {"uuid"},
	"sos_posts_dates":      {"uuid", "sos_post_uuid", "sos_dates_uuid"},
	"sos_conditions":       {"uuid"},
	"sos_posts_conditions": {"uuid", "sos_post_uuid", "sos_condition_uuid"},
	"sos_posts_pets":       {"uuid", "sos_post_uuid", "pet_uuid"},
	"resource_media":       {"uuid", "media_uuid", "resource_uuid"},
}

func main() {
	log.Print("Running UUID migration script")

	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("error opening database")
	}

	for tableName, columns := range targetTables {
		log.Printf("Processing table %s", tableName)
		for _, column := range columns {
			log.Printf(" - Column %s", column)
		}
	}

	log.Print("Starting migration")
	ctx := context.Background()
	pndErr := migrateUUID(ctx, db)
	if pndErr != nil {
		log.Fatal().Err(pndErr.Err).Msg("error migrating UUIDs")
	}
	log.Print("Completed migration")
}

func migrateUUID(ctx context.Context, db *database.DB) *pnd.AppError {
	type Row struct {
		ID int
	}

	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}

	for tableName, columns := range targetTables {
		log.Printf("Processing table %s", tableName)

		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
		if err != nil {
			return pnd.ErrUnknown(fmt.Errorf("error counting rows from table %s due to %w", tableName, err))
		}
		log.Printf(" - Total rows: %d", count)

		rows, err := tx.Query("SELECT id FROM " + tableName + " ORDER BY created_at ASC")
		if err != nil {
			return pnd.ErrUnknown(
				fmt.Errorf("error selecting rows from table %s due to %w", tableName, err),
			)
		}
		defer rows.Close()

		for rows.Next() {
			var row Row
			err = rows.Scan(&row.ID)
			if err != nil {
				return pnd.ErrUnknown(
					fmt.Errorf("error scanning row from table %s with id %d due to %w", tableName, row.ID, err),
				)
			}
			for _, column := range columns {
				newUUID, err := uuid.NewV7()
				if err != nil {
					return pnd.ErrUnknown(
						fmt.Errorf("error generating new UUID for table %s column %s due to %w", tableName, column, err),
					)
				}

				log.Printf("SQL: UPDATE %s SET %s = %s WHERE id = %d", tableName, column, newUUID, row.ID)
				// _, err = tx.Exec("UPDATE "+tableName+" SET "+column+" = $1 WHERE id = $2", newUUID, row.ID)
				// if err != nil {
				// 	return pnd.ErrUnknown(
				// 		fmt.Errorf(
				// 			"error updating table %s column %s with new UUID %s for id %d due to %w",
				// 			tableName, column, newUUID, row.ID, err),
				// 	)
				// }
			}
		}
	}

	return nil
}
