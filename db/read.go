package db

import (
	"database/sql"
	"errors"
	"github.com/mono83/artifacts/data"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
)

// Read reads artifact data from database
func Read(db *sql.DB, a data.Artifact) (*data.ResultsTable, error) {
	if db == nil {
		return nil, errors.New("nil database")
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}

	ray := xray.ROOT.Fork().With(args.SQL(a.Query))

	// Performing query
	ray.Debug("Performing query :sql")
	rows, err := db.Query(a.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Analyzing result
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var out = data.ResultsTable{
		Metric: a.Metric,
		Groups: cols[0 : len(cols)-1],
	}
	for rows.Next() {
		groups := make([]string, len(cols)-1)
		var value int64 = 0
		cells := make([]interface{}, len(cols))
		for i := 0; i < len(groups); i++ {
			cells[i] = &groups[i]
		}
		cells[len(cells)-1] = &value

		if err := rows.Scan(cells...); err != nil {
			return nil, err
		}

		var group map[string]string
		if len(groups) > 0 {
			group = make(map[string]string)
			for i, v := range groups {
				group[cols[i]] = v
			}
		}

		out.Add(value, groups...)
	}

	return &out, nil
}
