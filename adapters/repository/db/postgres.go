package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Postgres struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
	ClientCert string `yaml:"client_cert"`
	ClientKey  string `yaml:"client_key"`
	ServerCert string `yaml:"server_cert"`
}

func NewPostgres(p Postgres, wg *sync.WaitGroup, loglevel bool) *Model {

	var (
		db    *bun.DB
		sqldb *sql.DB
		err   error
		//dsn       string
		tlsConfig *tls.Config
		tlog      *log.Logger
	)

	tlog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//dsn = fmt.Sprintf(
	//	"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	//	p.Host, p.Port, p.User, p.Password, p.Database)

	tlsConfig = ConfTLS(p.ClientCert, p.ClientKey, p.ServerCert)
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(p.Host+":"+p.Port),
		pgdriver.WithTLSConfig(tlsConfig),
		pgdriver.WithUser(p.User),
		pgdriver.WithPassword(p.Password),
		pgdriver.WithDatabase(p.Database),
	)

	sqldb = sql.OpenDB(pgconn)
	if sqldb == nil {
		tlog.Fatal("failed to open database connection")
	}

	db = bun.NewDB(sqldb, pgdialect.New())

	//sqldb, err = db.DB()

	err = sqldb.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// trap SIGINT / SIGTERM to exit cleanly
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Shutting down DB...")
		_ = sqldb.Close()
		fmt.Println("BD connection closed gracefully")
		wg.Done()
	}()

	return &Model{
		db:       db,
		sqldb:    sqldb,
		loglevel: loglevel,
		errorLog: tlog,
		Postgres: p,
		wg:       wg,
	}

}

func ConfTLS(clientCert string, clientKey string, serverCert string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Println("failed to load client certificate: %v", err)
	}

	CACert, err := os.ReadFile(serverCert)
	if err != nil {
		log.Println("failed to load server certificate: %v", err)
	}

	CACertPool := x509.NewCertPool()
	CACertPool.AppendCertsFromPEM(CACert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            CACertPool,
		InsecureSkipVerify: true,
	}
}
