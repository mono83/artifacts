package influx

import (
	"bytes"
	"github.com/mono83/artifacts/data"
	"net"
	"strconv"
	"strings"
)

// Send sends data to influxdb using udp
func Send(addr string, results []data.Result) error {
	if len(results) == 0 {
		return nil
	}
	// Building data to send
	b := bytes.NewBuffer(nil)
	for _, result := range results {
		b.WriteString(result.Metric)
		if len(result.Group) > 0 {
			for k, v := range result.Group {
				b.WriteRune(',')
				b.Write(Sanitize(k))
				b.WriteRune('=')
				b.Write(Sanitize(v))
			}
		}
		b.WriteString(" value=")
		b.WriteString(strconv.FormatInt(result.Value, 10))
		b.WriteRune('\n')
	}

	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(b.Bytes())
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
