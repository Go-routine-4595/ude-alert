package simulationpackage

import (
	"fmt"
	"github.com/Go-routine-4595/ude-alert/domain"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Consumer interface {
	readData() int
}

type DataGen struct {
	sim     Consumer
	svc     domain.IService
	maxPeak int
	freq    int
	wg      *sync.WaitGroup
	log     zerolog.Logger
}

func NewDataGen(freq int, maxPeak int, svc domain.IService, wg *sync.WaitGroup) *DataGen {

	var sim *DataGen

	if freq == 0 {
		freq = 1
	}
	if maxPeak == 0 {
		maxPeak = 1
	}

	sim = &DataGen{
		maxPeak: maxPeak,
		freq:    freq,
		svc:     svc,
		sim:     NewSim(freq, maxPeak),
		wg:      wg,
		log:     zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.DebugLevel).With().Timestamp().Logger(),
	}

	return sim
}

func (d *DataGen) Start() {
	go d.start()
}

func (d *DataGen) start() {

	var (
		count  int
		ticker *time.Ticker
		err    error
	)

	ticker = time.NewTicker(time.Duration(d.freq) * time.Second)

	// trap SIGINT / SIGTERM to exit cleanly
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		ticker.Stop()
		fmt.Println("Shutting down Simulation...")
		d.wg.Done()
	}()

	err = d.svc.LoadEquipment(d.maxPeak)
	if err != nil {
		d.log.Fatal().Err(err).Msg("Fatal")
	}

	for {
		select {
		case <-ticker.C:
			// Gen the data here
			count = d.sim.readData()
			d.log.Debug().Int("count", count).Msg("Read data")
			_ = d.svc.UpdateEquipment(count)
		}
	}
}
