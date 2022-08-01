package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultsTableWithoutColumn(t *testing.T) {
	r := ResultsTable{Metric: "foo", Groups: []string{"A", "B"}}
	r.Add(24, "a1", "b1")
	r.Add(614, "a2", "b2")

	r0 := r.WithoutColumn(1)
	r1 := r.WithoutColumn(0)

	if assert.Len(t, r0.Groups, 1) {
		assert.Equal(t, "A", r0.Groups[0])
		assert.Equal(t, "a1", r0.Rows[0].Groups[0])
		assert.Equal(t, "a2", r0.Rows[1].Groups[0])
		assert.Equal(t, int64(24), r0.Rows[0].Value)
		assert.Equal(t, int64(614), r0.Rows[1].Value)

		rz := r0.WithoutColumn(0)
		assert.Len(t, rz.Groups, 0)
	}
	if assert.Len(t, r1.Groups, 1) {
		assert.Equal(t, "B", r1.Groups[0])
		assert.Equal(t, "b1", r1.Rows[0].Groups[0])
		assert.Equal(t, "b2", r1.Rows[1].Groups[0])
		assert.Equal(t, int64(24), r1.Rows[0].Value)
		assert.Equal(t, int64(614), r1.Rows[1].Value)
	}
}

func TestSortByGroupValues(t *testing.T) {
	r := ResultsTable{Metric: "foo", Groups: []string{"Gate", "Status", "Enabled"}}
	r.Add(5, "1", "activated", "true")
	r.Add(3, "1", "in_progress", "true")
	r.Add(8, "1", "activated", "false")
	r.Add(2, "2", "investigation", "true")
	r.Add(2, "2", "activated", "true")
	r.SortByGroupValues()

	assert.Equal(t, []string{"1", "activated", "false"}, r.Rows[0].Groups)
	assert.Equal(t, []string{"1", "activated", "true"}, r.Rows[1].Groups)
	assert.Equal(t, []string{"1", "in_progress", "true"}, r.Rows[2].Groups)
	assert.Equal(t, []string{"2", "activated", "true"}, r.Rows[3].Groups)
	assert.Equal(t, []string{"2", "investigation", "true"}, r.Rows[4].Groups)
}

func TestMergeSorted(t *testing.T) {
	r := ResultsTable{Metric: "foo", Groups: []string{"A"}}
	r.Add(23, "bar")
	r.Add(1, "bar")
	r.Add(10, "foo")
	r.Add(10, "foo")
	r.Add(10, "foo")

	s := r.mergeSorted()
	if assert.Len(t, s.Rows, 2) {
		assert.Equal(t, "bar", s.Rows[0].Groups[0])
		assert.Equal(t, "foo", s.Rows[1].Groups[0])

		assert.Equal(t, int64(24), s.Rows[0].Value)
		assert.Equal(t, int64(30), s.Rows[1].Value)
	}
}

func TestRecombine(t *testing.T) {
	r := ResultsTable{Metric: "foo", Groups: []string{"Gate", "Status", "Enabled"}}
	r.Add(5, "1", "activated", "true")
	r.Add(3, "1", "in_progress", "true")
	r.Add(8, "1", "activated", "false")
	r.Add(2, "2", "investigation", "true")
	r.Add(2, "2", "activated", "true")
	r.Add(6, "1", "activated", "true")

	if recombined := r.Recombine(); len(recombined) == 7 {
		if len(recombined[0].Rows) == 4 {
			assert.Equal(t, ResultsRow{Value: 8, Groups: []string{"activated", "false"}}, recombined[0].Rows[0])
			assert.Equal(t, ResultsRow{Value: 13, Groups: []string{"activated", "true"}}, recombined[0].Rows[1])
			assert.Equal(t, ResultsRow{Value: 3, Groups: []string{"in_progress", "true"}}, recombined[0].Rows[2])
			assert.Equal(t, ResultsRow{Value: 2, Groups: []string{"investigation", "true"}}, recombined[0].Rows[3])
		}

		if len(recombined[1].Rows) == 2 {
			assert.Equal(t, ResultsRow{Value: 8, Groups: []string{"false"}}, recombined[1].Rows[0])
			assert.Equal(t, ResultsRow{Value: 18, Groups: []string{"true"}}, recombined[1].Rows[1])
		}

		if len(recombined[2].Rows) == 1 {
			assert.Equal(t, ResultsRow{Value: 26}, recombined[2].Rows[0])
		}

		if len(recombined[3].Rows) == 3 {
			assert.Equal(t, ResultsRow{Value: 21, Groups: []string{"activated"}}, recombined[3].Rows[0])
			assert.Equal(t, ResultsRow{Value: 3, Groups: []string{"in_progress"}}, recombined[3].Rows[1])
			assert.Equal(t, ResultsRow{Value: 2, Groups: []string{"investigation"}}, recombined[3].Rows[2])
		}

		if len(recombined[4].Rows) == 3 {
			assert.Equal(t, ResultsRow{Value: 8, Groups: []string{"1", "false"}}, recombined[4].Rows[0])
			assert.Equal(t, ResultsRow{Value: 14, Groups: []string{"1", "true"}}, recombined[4].Rows[1])
			assert.Equal(t, ResultsRow{Value: 4, Groups: []string{"2", "true"}}, recombined[4].Rows[2])
		}

		if len(recombined[5].Rows) == 2 {
			assert.Equal(t, ResultsRow{Value: 22, Groups: []string{"1"}}, recombined[5].Rows[0])
			assert.Equal(t, ResultsRow{Value: 4, Groups: []string{"2"}}, recombined[5].Rows[1])
		}

		if len(recombined[6].Rows) == 4 {
			assert.Equal(t, ResultsRow{Value: 19, Groups: []string{"1", "activated"}}, recombined[6].Rows[0])
			assert.Equal(t, ResultsRow{Value: 3, Groups: []string{"1", "in_progress"}}, recombined[6].Rows[1])
			assert.Equal(t, ResultsRow{Value: 2, Groups: []string{"2", "activated"}}, recombined[6].Rows[2])
			assert.Equal(t, ResultsRow{Value: 2, Groups: []string{"2", "investigation"}}, recombined[6].Rows[3])
		}
	}
}
