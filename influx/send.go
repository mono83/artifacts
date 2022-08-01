package influx

import (
	"bytes"
	"fmt"
	"github.com/mono83/artifacts/data"
	"net"
	"strconv"
	"strings"
)

// Send sends data to influxdb using udp
func Send(addr string, tables []data.ResultsTable) error {
	if len(tables) == 0 {
		return nil
	}
	// Dialing
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Sending data
	for _, t := range tables {
		if err := sendSingleTable(conn, t); err != nil {
			return err
		}
	}

	return nil
}

func sendSingleTable(conn net.Conn, t data.ResultsTable) error {
	b := bytes.NewBuffer(nil)
	for _, row := range t.Rows {
		b.WriteString(t.Metric)
		if len(t.Groups) > 0 {
			for k, v := range row.Groups {
				b.WriteRune(',')
				b.Write(Sanitize(t.Groups[k]))
				b.WriteRune('=')
				b.Write(Sanitize(v))
			}
		}
		b.WriteString(" value=")
		b.WriteString(strconv.FormatInt(row.Value, 10))
		b.WriteRune('\n')
	}
	fmt.Println(b)
	_, err := conn.Write(b.Bytes())
	return err
}

var sanitizeReplacement = byte('_')

// Sanitize function takes string value and removes values, that can be unsafe
// for remote receivers. In fact, it leaves only alpha numeric values, others are
// replaced with underscore
func Sanitize(value string) []byte {
	if len(value) == 0 {
		return []byte{}
	}

	bts := []byte(strings.TrimSpace(value))
	for i, v := range bts {
		if !(v == 46 || (v >= 48 && v <= 57) || (v >= 65 && v <= 90) || (v >= 97 && v <= 122)) {
			bts[i] = sanitizeReplacement
		}
	}

	return bts
}
