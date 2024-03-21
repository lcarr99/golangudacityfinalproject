package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Migration interface {
	GetName() string
	Up(connection *sql.DB) error
	Down(connection *sql.DB) error
}

type CreateCustomersTable struct{}

func (cct CreateCustomersTable) GetName() string {
	return "create_customers_table"
}

func (cct CreateCustomersTable) Up(connection *sql.DB) error {
	_, err := connection.Exec("CREATE TABLE IF NOT EXISTS customers (id INTEGER auto_increment, name varchar(100), role varchar(50), email varchar(100), phone varchar(20), contacted tinyint, PRIMARY KEY (id))")
	return err
}

func (cct CreateCustomersTable) Down(connection *sql.DB) error {
	_, err := connection.Exec("DROP TABLE IF EXISTS customers")
	return err
}

type InsertFakeCustomers struct{}

func (cct InsertFakeCustomers) GetName() string {
	return "insert_fake_customers"
}

func (cct InsertFakeCustomers) Up(connection *sql.DB) error {
	_, err := connection.Exec("INSERT INTO customers (name, `role`, email, phone, contacted) VALUES('John Doe', 'CEO', 'john.doe@test.com', '12345678910', 1);")

	if err != nil {
		return err
	}

	_, err = connection.Exec("INSERT INTO customers (name, `role`, email, phone, contacted) VALUES('Jane Doe', 'CTO', 'jane.doe@test.com', '12345678940', 1);")

	if err != nil {
		return err
	}

	_, err = connection.Exec("INSERT INTO customers(name, `role`, email, phone, contacted) VALUES('Jack Doe', 'Developer', 'jack.doe@test.com', '12345678900', 0);")

	return err
}

func (cct InsertFakeCustomers) Down(connection *sql.DB) error {
	return nil
}

var migrations = []Migration{
	CreateCustomersTable{},
	InsertFakeCustomers{},
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Please ensure you have a .env set")
	}

	var address = fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	config := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   address,
		DBName: os.Getenv("DB_DATABASE")}

	connection, err := sql.Open("mysql", config.FormatDSN())

	if err != nil {
		log.Fatalf("Error when connecting to the database")
	}

	connection.Exec("CREATE TABLE IF NOT EXISTS migrations(name varchar(100), executed_datetime datetime)")

	var lastMigration Migration

	for _, migration := range migrations {
		var count int
		migrationSearch := connection.QueryRow("SELECT count(*) as count FROM migrations WHERE name = ?", migration.GetName())

		migrationSearch.Scan(&count)

		if count > 0 {
			continue
		}

		err := migration.Up(connection)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		currentTime := time.Now()

		connection.Exec("INSERT INTO migrations(name, executed_datetime) VALUES (?, ?)", migration.GetName(), currentTime.Format(time.DateTime))

		lastMigration = migration
	}

	if lastMigration != nil {
		fmt.Printf("Migrated up to %s \n", lastMigration.GetName())
	} else {
		fmt.Println("Up to date...")
	}
}
