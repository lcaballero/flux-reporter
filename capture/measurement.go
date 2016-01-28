package capture

import (
	"fmt"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

// A Capture is the re-useable portions of a Influx data point.
type Measurement struct {
	measurement string
	tags        map[string]string
	fields      map[string]interface{}
}

// New creates a new Capture with the given measurement name.
func New(name string) (*Measurement, error) {
	if name == "" {
		return nil, fmt.Errorf("A Measurement should have a non empty name: '%s'", name)
	}
	m := &Measurement{
		measurement: name,
		tags:        make(map[string]string),
		fields:      make(map[string]interface{}),
	}
	return m, nil
}

// Tag saves and/or overrides the key/value pair.
func (mm *Measurement) Tag(k, v string) *Measurement {
	if k != "" && v != "" {
		mm.tags[k] = v
	}
	return mm
}

// Field saves the key and value, and overrides any previous field with the
// given key.
func (mm *Measurement) Field(k string, v interface{}) *Measurement {
	if k != "" && v != nil {
		mm.fields[k] = v
	}
	return mm
}

// Point stamps out a new influx data point and COPIES the tags, and CLEARS
// the fields for this Capture so that it can be fields with new values.
// The new data point is stamped with the current time.
func (mm *Measurement) Point() (*influx.Point, error) {
	tags := mm.tags
	fields := mm.fields

	mm.tags = make(map[string]string)
	mm.fields = make(map[string]interface{})

	for k, v := range tags {
		mm.tags[k] = v
	}
	for k, v := range fields {
		mm.fields[k] = v
	}

	pt, err := influx.NewPoint(
		mm.measurement,
		mm.tags,
		mm.fields,
		time.Now(),
	)

	if err != nil {
		return nil, err
	} else {
		return pt, nil
	}
}
