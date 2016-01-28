package capture

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMeasurement(t *testing.T) {

	Convey("name, tags, fields and timestamp should be propagated to new Point.", t, func() {
		measurement := "measure.this"
		machine := "machine01"

		m, _ := New(measurement)
		m.Field("value", 10)
		m.Tag("service", machine)
		pt, err := m.Point()

		So(err, ShouldBeNil)

		name := pt.Name()
		fields := pt.Fields()
		tags := pt.Tags()

		So(name, ShouldEqual, measurement)
		So(tags["service"], ShouldEqual, machine)
		So(fields["value"], ShouldEqual, 10)
	})

	Convey("Name should be propagated to new Point.", t, func() {
		name := "measure.this"
		m, _ := New(name)
		m.Field("value", 10)
		pt, err := m.Point()

		So(err, ShouldBeNil)
		So(pt.Name(), ShouldEqual, name)
	})

	Convey("A new Measurement without fields should cause an error.", t, func() {
		m, err := New("name")
		So(err, ShouldBeNil)

		_, err = m.Point()
		So(err, ShouldNotBeNil)
	})

	Convey("A named Measurement should be fine.", t, func() {
		m, err := New("so.named")
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
	})

	Convey("A UN-named Measurement should produce an error.", t, func() {
		m, err := New("")
		So(err, ShouldNotBeNil)
		So(m, ShouldBeNil)
	})
}
