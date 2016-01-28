package start

import (
	"time"

	"fmt"
	"math/rand"

	flux "github.com/influxdata/influxdb/client/v2"
	"github.com/lcaballero/flux-reporter/capture"
	"github.com/lcaballero/flux-reporter/reporter"
	"github.com/lcaballero/hitman"
//	"github.com/vrecan/death"
//	"syscall"
)

func Start() {
	conf := flux.BatchPointsConfig{
		Precision:        "",
		Database:         "",
		RetentionPolicy:  "",
		WriteConsistency: "",
	}

	stats := make(chan flux.BatchPoints, 10)
	points := make(chan *flux.Point, 10000)

	targets := hitman.NewTargets()
	targets.AddOrPanic(reporter.NewFlusher(stats))
	targets.AddOrPanic(reporter.NewReporter(conf,stats,points))

	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)
	Fire(points)

	<-time.After(300 * time.Second)
	targets.Close()
}

func Fire(add chan *flux.Point) {
	increment := time.Millisecond * 20
	m, err := capture.New("client.request")
	if err != nil {
		panic(err)
	}
	m.Tag("service", "fire")
	m.Tag("ip", "127.0.0.1")

	go func() {
		for {
			n := time.Duration(rand.Intn(10))
			r := n * increment
			<-time.After(r)
			m.Field("value", r)
			m.Field("now", r)

			pt, err := m.Point()
			if err != nil {
				fmt.Println(err)
			}
			add <- pt
		}
	}()
}
