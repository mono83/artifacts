package db

import (
	"database/sql"
	"errors"
	"github.com/mono83/artifacts/data"
	"time"
)

/*
CREATE TABLE `__artifacts` (
  `name` varchar(1024) CHARACTER SET latin1 NOT NULL,
  `query` varchar(2048) COLLATE utf8mb4_unicode_ci NOT NULL,
  `interval` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci*/

// ReadFromConfigTable reads artifacts configuration from database table
func ReadFromConfigTable(db *sql.DB, table string) ([]data.Artifact, error) {
	if db == nil {
		return nil, errors.New("nil database")
	}

	rows, err := db.Query("SELECT `name`,`query`,`interval` FROM " + table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []data.Artifact
	var name, query string
	var interval int
	for rows.Next() {
		if err = rows.Scan(&name, &query, &interval); err != nil {
			return nil, err
		}

		out = append(out, data.Artifact{
			Metric:   name,
			Query:    query,
			Interval: time.Second * time.Duration(interval),
		})
	}
	return out, nil
}
