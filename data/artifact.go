package data

import (
	"errors"
	"time"
)

// Artifact is some measure obtained from database
type Artifact struct {
	Metric   string
	Query    string
	Interval time.Duration
}

// Validate performs validation of artifact configuration
func (a Artifact) Validate() error {
	if len(a.Metric) == 0 {
		return errors.New("empty metric name")
	}
	if len(a.Query) == 0 {
		return errors.New("empty query")
	}
	return nil
}

// IntervalOrDefault returns configured refresh interval or default value
func (a Artifact) IntervalOrDefault(d time.Duration) time.Duration {
	if a.Interval <= 0 {
		// Interval not set or negative - using defaule
		return d
	}
	return a.Interval
}
