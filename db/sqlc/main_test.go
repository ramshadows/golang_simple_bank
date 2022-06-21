package db

import (
	"database/sql"
	"log"
	"os"

	"simple_bank/utils"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// load app configs
	config, err := utils.LoadConfigs("../..") //"." means parent folder

	if err != nil {
		log.Fatal("Cannot load app configs", err)
	}

	// Create a Postgres DB connection
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	testQueries = New(testDB)

	// Run the test
	os.Exit(m.Run())
}
