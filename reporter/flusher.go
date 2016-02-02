package reporter

import (
	"bytes"
	"fmt"

	flux "github.com/influxdata/influxdb/client/v2"
	"github.com/lcaballero/hitman"
)

type Flusher struct {
	points chan flux.BatchPoints
}

func NewFlusher(pts chan flux.BatchPoints) (*Flusher, error) {
	f := &Flusher{
		points: pts,
	}
	return f, nil
}

func (a *Flusher) Name() string {
	return "flusher of stats to influx"
}

func (a *Flusher) Start() hitman.KillChannel {
	killSwitch := hitman.NewKillChannel()

	go func() {
		for {
			select {
			case cleaner := <-killSwitch:
				cleaner.WaitGroup.Done()
				return
			case batch := <-a.points:
				pts := batch.Points()
				for n,p := range pts {
					a.flush(p, n)
				}
			}
		}
	}()

	return killSwitch
}

func (a *Flusher) flush(pt *flux.Point, k int) {
	//	fmt.Println("flushing data points to influx", k)
	return
	name := pt.Name()
	tags := pt.Tags()
	fields := pt.Fields()
	nanos := pt.Time()

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(name)

	n := len(tags)
	if n > 0 {
		buf.WriteString(" ")
		for k, v := range tags {
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v)
			n--
			if n > 0 {
				buf.WriteString(",")
			} else {
				break
			}
		}
	}

	d := len(fields)
	if d > 0 {
		buf.WriteString(" ")
		s := "v"
		for k, _ := range fields {
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(s)
			d--
			if d > 0 {
				buf.WriteString(",")
			} else {
				break
			}
		}
	}

	buf.WriteString(" ")
	buf.WriteString(fmt.Sprintf("%di", nanos.UnixNano()))

//	fmt.Println(buf.String())
}
