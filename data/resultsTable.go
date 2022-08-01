package data

import (
	"sort"
)

// ResultsRow contains information on single row of data
type ResultsRow struct {
	Groups []string // Group values
	Value  int64
}

// WithoutColumn returns row with one group column removed
func (r ResultsRow) WithoutColumn(c int) ResultsRow {
	return ResultsRow{Value: r.Value, Groups: sliceRemove(r.Groups, c)}
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

// WithoutColumn returns results table with one column removed
func (r ResultsTable) WithoutColumn(c int) ResultsTable {
	// Rebuilding group values
	var rows []ResultsRow
	for _, row := range r.Rows {
		rows = append(rows, row.WithoutColumn(c))
	}

	return ResultsTable{
		Metric: r.Metric,
		Groups: sliceRemove(r.Groups, c),
		Rows:   rows,
	}
}

// SortByGroupValues performs table data sorting by group values
func (r *ResultsTable) SortByGroupValues() {
	sort.Sort(r)
}

// Recombine recombines column data
func (r *ResultsTable) Recombine() []ResultsTable {
	var skip [][]string
	return r.recombineSorted(&skip)
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

func (r *ResultsTable) recombineSorted(skip *[][]string) (out []ResultsTable) {
	for x := 0; x < len(r.Groups); x++ {
		truncated := r.WithoutColumn(x)
		todo := true
		for _, s := range *skip {
			if sliceEquals(s, truncated.Groups) {
				todo = false
				break
			}
		}
		if todo {
			*skip = append(*skip, truncated.Groups)
			truncated.SortByGroupValues()
			merged := truncated.mergeSorted()
			out = append(out, merged)
			if len(merged.Groups) > 0 {
				out = append(out, merged.recombineSorted(skip)...)
			}
		}
	}
	return
}

// Len is sort.Interface implementation
func (r ResultsTable) Len() int { return len(r.Rows) }

// Swap is sort.Interface implementation
func (r *ResultsTable) Swap(i, j int) { r.Rows[i], r.Rows[j] = r.Rows[j], r.Rows[i] }

// Less is sort.Interface implementation
func (r *ResultsTable) Less(i, j int) bool {
	ri, rj := r.Rows[i].Groups, r.Rows[j].Groups
	for k := range ri {
		if ri[k] == rj[k] {
			continue
		}
		return ri[k] < rj[k]
	}
	return false
}

func sliceRemove(slice []string, s int) (out []string) {
	out = append(out, slice[:s]...)
	out = append(out, slice[s+1:]...)
	return
}

func sliceEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
