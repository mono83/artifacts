package db

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

// Connect establishes connection to MySQL server on given DSN
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err == nil {
		err = db.Ping()
	}
	return db, err
}

// ConnectStdMySQL establishes connection to MySQL server
func ConnectStdMySQL(addr, login, password string) (*sql.DB, error) {
	c := &mysql.Config{
		User:                 login,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 addr,
		AllowNativePasswords: true,
	}
	return Connect(c.FormatDSN())
}
