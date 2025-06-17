package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/lorendsnow/updater/cmd"
)

// UpdateService periodically downloads csv files from the City's website and
// updates the database.
//
// The City's website notes that every year is updated monthly to account for
// updates on cases from prior years, so we need to download the data for each
// year, not just the current year.
//
// Since there are no unique keys we can extract from the data, we also can't
// really run an upsert type statement into a single table. So we use a
// blue/green deployment strategy within the database, cycling between two
// different tables and maintaining which one was most recently updated for the
// repository to check before querying.
//
// In the future this should probably push updates to the repository instead of
// the repository pulling the table to use from the database. This could be tied
// into a cache used by the repository, or via a message/event type of service.
type UpdateService struct {
	CheckEvery string
	BlueTable  *Table
	GreenTable *Table
	Db         *sql.DB
}

// Table represents one of the two blue/green tables the UpdateService will
// update, holding the table name and its last update datetime
type Table struct {
	Name        string
	LastUpdated time.Time
}

// NewUpdateService creates a new UpdateService with the given update interval.
//
// The UpdateService will check for updates every updateEvery duration, and
// will use the blue and green tables to store the data.
func NewUpdateService(config *cmd.Config) *UpdateService {
	return &UpdateService{
		CheckEvery: config.Service.CheckInterval,
		BlueTable:  &Table{Name: config.Service.BlueTable},
		GreenTable: &Table{Name: config.Service.GreenTable},
	}
}

// LastUpdatedTable returns the name of the table that was most recently updated.
//
// This is used by the repository to determine which table to query.
func (s *UpdateService) LastUpdatedTable() string {
	if s.BlueTable.LastUpdated.After(s.GreenTable.LastUpdated) {
		return s.BlueTable.Name
	}

	return s.GreenTable.Name
}

// ConnectToDatabase connects to the database using the given configuration.
func (s *UpdateService) ConnectToDatabase(config *cmd.Config) error {
	dbConfig := mysql.Config{
		User:   config.Database.Username,
		Passwd: config.Database.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", config.Database.Host, config.Database.Port),
		DBName: config.Database.Name,
	}

	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return err
	}

	// Ping the database to make sure we have a real connection.
	if err := db.Ping(); err != nil {
		return err
	}

	s.Db = db

	return nil
}
