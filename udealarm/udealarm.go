package udealarm

import (
	"fmt"
	"github.com/Go-routine-4595/ude-alert/domain"
	"github.com/rs/zerolog"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/Go-routine-4595/ude-alert/adapters/repository/db"
	simulationpackage "github.com/Go-routine-4595/ude-alert/adapters/simulation"
	"github.com/Go-routine-4595/ude-alert/service"
	"gopkg.in/yaml.v3"
)

type FreqItem struct {
	Frequency int `yaml:"frequency"`
	MaxPeak   int `yaml:"max_peak"`
}

const defaultConfigFile = "config.yaml"

type Config struct {
	Postgresql db.Postgres `yaml:"service"`
	Freq       FreqItem    `yaml:"frequency"`
}

func StartSim(conf string) {
	var (
		wg       *sync.WaitGroup
		database service.Storer
		svc      domain.IService
		sim      *simulationpackage.DataGen
	)

	if conf == "" {
		conf = defaultConfigFile
	}
	cfg := openFile(conf)

	wg = &sync.WaitGroup{}
	wg.Add(1)
	if cfg.Postgresql.Host != "" {
		database = db.NewPostgres(cfg.Postgresql, wg, true)
	}

	// create our service logic
	svc = service.NewService(database)

	// new simulator
	wg.Add(1)
	sim = simulationpackage.NewDataGen(cfg.Freq.Frequency, cfg.Freq.MaxPeak, svc, wg)
	sim.Start()
	wg.Wait()

}

func AddEquipment(conf string, data string) error {
	var (
		wg       *sync.WaitGroup
		err      error
		database service.Storer
		svc      domain.IService
		zlog     zerolog.Logger
	)

	zlog = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	cfg := openFile(conf)

	wg = &sync.WaitGroup{}
	wg.Add(1)
	if cfg.Postgresql.Host != "" {
		database = db.NewPostgres(cfg.Postgresql, wg, true)
	}

	// create our service logic
	svc = service.NewService(database)

	err = svc.AddEquipment([]byte(data))
	if err != nil {
		zlog.Error().Err(err).Msg("error adding equipment")
	}

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	wg.Wait()
	return err

}

func openFile(s string) Config {
	f, err := os.Open(s)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
	}

	return cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
