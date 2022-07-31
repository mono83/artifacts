package data

import "sort"

// ResultsRow contains information on single row of data
type ResultsRow struct {
	Groups []string // Group values
	Value  int64
}

// Without column returns row with one group column removed
func (r ResultsRow) WithoutColumn(c int) ResultsRow {
	return ResultsRow{Value: r.Value, Groups: remove(r.Groups, c)}
}

// SameGroups return true if current and other results rows
// contains same group values
func (r ResultsRow) SameGroups(other ResultsRow) bool {
	if len(r.Groups) == len(other.Groups) {
		for i, v := range r.Groups {
			if v != other.Groups[i] {
				return false
			}
		}
		return true
	}
	return false
}

// ResultsTable contains resulting data obtained from database source
type ResultsTable struct {
	Metric string
	Groups []string // Group names
	Rows   []ResultsRow
}

// Add places row data into results table
func (r *ResultsTable) Add(value int64, groups ...string) {
	r.Rows = append(r.Rows, ResultsRow{Value: value, Groups: groups})
}

// Without column returns results table with one column removed
func (r ResultsTable) WithoutColumn(c int) ResultsTable {
	// Rebuilding group values
	var rows []ResultsRow
	for _, row := range r.Rows {
		rows = append(rows, row.WithoutColumn(c))
	}

	return ResultsTable{
		Metric: r.Metric,
		Groups: remove(r.Groups, c),
		Rows:   rows,
	}
}

// SortByGroupValues performs table data sorting by group values
func (r *ResultsTable) SortByGroupValues() {
	sort.Slice(r.Rows, func(i, j int) bool {
		ri, rj := r.Rows[i].Groups, r.Rows[j].Groups
		for k, v := range ri {
			if v < rj[k] {
				return true
			}
		}
		return false
	})
}

// Recombine recombines column data
func (r *ResultsTable) Recombine() ResultsTable {
	// Sorting
	r.SortByGroupValues()

	return r.recombineSorted()
}

func (r ResultsTable) mergeSorted() ResultsTable {
	var rows []ResultsRow
	for i, row := range r.Rows {
		if i == 0 || !row.SameGroups(rows[len(rows)-1]) {
			rows = append(rows, row)
		} else {
			rows[len(rows)-1].Value += row.Value
		}
	}
	return ResultsTable{
		Metric: r.Metric,
		Groups: r.Groups,
		Rows:   rows,
	}
}

func (r ResultsTable) recombineSorted() ResultsTable {
	for x := 0; x < len(r.Groups); x++ {
		truncated := r.WithoutColumn(x)

	}
}

func remove(slice []string, s int) (out []string) {
	out = append(out, slice[:s]...)
	out = append(out, slice[s+1:]...)
	return
}
