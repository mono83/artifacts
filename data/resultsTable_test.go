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
