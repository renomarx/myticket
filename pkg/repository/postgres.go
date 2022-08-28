package repository

import (
	"fmt"
	"os"

	"github.com/renomarx/myticket/pkg/core/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type postgresDB struct {
	db *sqlx.DB
}

func NewPostgresDB() *postgresDB {
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		logrus.Fatalln(err)
	}
	return &postgresDB{
		db: db,
	}
}

func (db *postgresDB) CreateTicket(ticket *model.Ticket) error {
	res, err := db.db.NamedExec(`
		INSERT INTO tickets (
			body,
			status,
			error_details,
			created_at,
			updated_at
		)
		VALUES (
			:body,
			:status,
			:error_details,
			:created_at,
			:updated_at
		) RETURNING id
	`,
		ticket)
	if err != nil {
		logrus.Error(err)
		model.AppMetrics.IncDatabaseErrors("Error creating ticket")
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		// Should not happen, as postgres supports it
		logrus.Error(err)
		return err
	}
	if rowsAffected == 0 {
		err = fmt.Errorf("0 ticket created")
		model.AppMetrics.IncDatabaseErrors(err.Error())
		return err
	}
	return nil
}

func (db *postgresDB) UpdateTicket(ticket *model.Ticket) error {
	res, err := db.db.NamedExec(`
		UPDATE tickets SET
			status = :status,
			updated_at = :updated_at
		WHERE id = :id
	`,
		ticket)
	if err != nil {
		logrus.Error(err)
		model.AppMetrics.IncDatabaseErrors("Error updating ticket")
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		// Should not happen, as postgres supports it
		logrus.Error(err)
		return err
	}
	if rowsAffected == 0 {
		err = fmt.Errorf("0 ticket updated")
		model.AppMetrics.IncDatabaseErrors(err.Error())
		return err
	}
	return nil
}

func (db *postgresDB) SaveProducts(products []model.Product) error {
	res, err := db.db.NamedExec(`
			INSERT INTO products (
				product_id,
				name,
				price,
				created_at,
				updated_at
			)
			VALUES (
				:product_id,
				:name,
				:price,
				:created_at,
				:updated_at
			) ON CONFLICT (product_id) DO UPDATE SET
					name = EXCLUDED.name,
					price = EXCLUDED.price,
					updated_at = EXCLUDED.updated_at
		`,
		products)
	if err != nil {
		logrus.Error(err)
		model.AppMetrics.IncDatabaseErrors("Error upserting products")
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected != int64(len(products)) {
		err = fmt.Errorf("Expected %d products updated, got %d", len(products), 0)
		model.AppMetrics.IncDatabaseErrors("Error upserting products")
		return err
	}
	return nil
}
