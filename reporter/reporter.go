package reporter

import (
	"fmt"
	"time"

	flux "github.com/influxdata/influxdb/client/v2"
	"github.com/lcaballero/hitman"
)


type Reporter struct {
	batch flux.BatchPoints
	accum  chan flux.BatchPoints
	conf flux.BatchPointsConfig
	points chan *flux.Point
}

func NewReporter(conf flux.BatchPointsConfig, stats chan flux.BatchPoints, points chan *flux.Point) (*Reporter, error) {
	bp, err := flux.NewBatchPoints(conf)
	if err != nil {
		return nil, err
	}
	r := &Reporter{
		accum:  stats,
		batch: bp,
		conf: conf,
		points: points,
	}
	return r, nil
}

func (a *Reporter) Name() string {
	return "Accumulator"
}

func (a *Reporter) Start() hitman.KillChannel {
	done := hitman.NewKillChannel()
	tripwire := 1 * time.Second
	tic := time.NewTicker(tripwire).C
	last := time.Now().UnixNano()
	_1000ms := time.Duration(1000)
	fuzz := time.Duration(10)

	go func() {
		for {
			select {
			case cleaner := <-done:
				cleaner.WaitGroup.Done()
				fmt.Println("Closing reporter")
				return
			case pt := <-a.points:
				a.batch.AddPoint(pt)
			case <-tic:
				now := time.Now().UnixNano()
				diff := now - last
				delta := time.Duration(diff) / time.Millisecond
				isOverdue := _1000ms < (delta + fuzz)

				if !isOverdue {
					continue
				}
				fmt.Println(len(a.accum), len(a.points), len(a.batch.Points()))
				bp := a.batch
				newPts, err := flux.NewBatchPoints(a.conf)
				if err != nil {
					fmt.Println(err)
				}
				a.batch = newPts
				a.accum <- bp
				last = time.Now().UnixNano()
			}
		}
	}()

	return done
}
