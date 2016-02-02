package start

import (
	"time"

	"fmt"
	"math/rand"

	flux "github.com/influxdata/influxdb/client/v2"
	"github.com/lcaballero/flux-reporter/capture"
	"github.com/lcaballero/flux-reporter/reporter"
	"github.com/lcaballero/hitman"
	"github.com/vrecan/death"
	"syscall"
	"flag"
	"runtime/pprof"
	"log"
	"os"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func Start() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	conf := flux.BatchPointsConfig{
		Precision:        "",
		Database:         "",
		RetentionPolicy:  "",
		WriteConsistency: "",
	}

	stats := make(chan flux.BatchPoints, 10000)
	points := make(chan *flux.Point, 10000)

	targets := hitman.NewTargets()
	targets.AddOrPanic(reporter.NewFlusher(stats))
	targets.AddOrPanic(reporter.NewReporter(conf,stats,points))

	closer := NewCloser(targets).AddAll(Fire(points, 200))

	death.NewDeath(syscall.SIGINT, syscall.SIGTERM).WaitForDeath(closer)
}

type closer struct {
	targets hitman.Targets
}
func NewCloser(tgs hitman.Targets) *closer {
	return &closer{
		targets: tgs,
	}
}
func (c *closer) Close() error {
	return c.targets.Close()
}
func (c *closer) AddAll(more... hitman.Target) *closer {
	for _,t := range more {
		c.targets.Add(t)
	}
	return c
}

type Fired struct {
	add chan *flux.Point
	numIncoming int
}

func Fire(add chan *flux.Point, n int) hitman.Target {
	return &Fired{
		add: add,
		numIncoming: n,
	}
}

func (f *Fired) Name() string {
	return "new fired"
}

func (f *Fired) Start() hitman.KillChannel {
	add := f.add
	m, err := capture.New("client.request")
	if err != nil {
		panic(err)
	}
	m.Tag("service", "fire")
	m.Tag("ip", "127.0.0.1")
	tic := time.NewTicker(20*time.Millisecond).C
	numIncoming := int32(f.numIncoming)
	kill := hitman.NewKillChannel()

	go func() {
		for {
			select {
			case cleaner := <-kill:
				cleaner.WaitGroup.Done()
				return
			case now := <-tic:
				n := rand.Int31n(numIncoming)
				r := now.UnixNano()
				for i := int32(0); i < n; i++ {
					m.Field("value", r)
					m.Field("now", r)

					pt, err := m.Point()
					if err != nil {
						fmt.Println(err)
					}
					add <- pt
				}
			}
		}
	}()

	return kill
}
