package config

import (
	"errors"
	"time"
)

// Configuration contains application configuration data
type Configuration struct {
	MySQLDSN     string `yaml:"mysql"`
	MySQLTable   string `yaml:"table"`
	InfluxDBAddr string `yaml:"influx"`
	Schedule     int    `yaml:"scheduleSeconds"`
}

// Validate performs validation of configuration
func (c Configuration) Validate() error {
	if len(c.MySQLDSN) == 0 {
		return errors.New("empty MySQL DSN")
	}
	if len(c.InfluxDBAddr) == 0 {
		return errors.New("empty InfluxDB address")
	}

	return nil
}

// MySQLTableOrDefault returns name of table where artifacts configuration resides
func (c Configuration) MySQLTableOrDefault() string {
	if len(c.MySQLTable) == 0 {
		return "__artifacts"
	}
	return c.MySQLTable
}

// ScheduleOrDefault returns schedule interval
func (c Configuration) ScheduleOrDefault() time.Duration {
	if c.Schedule <= 1 {
		return 2 * time.Minute
	}

	return time.Second * time.Duration(c.Schedule)
}
