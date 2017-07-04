package gopg

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	flog "github.com/joelclouddistrict/gochecklist/log"
	"github.com/joelclouddistrict/gochecklist/store"
	GOPG "gopkg.in/pg.v5"
)

// GopgStore is a Storer implementation over PostgresSQL
type GopgStore struct {
	Pool *GOPG.DB
	Stat int
}

var _ store.Storer = (*GopgStore)(nil)

// New is a default Storer implementation based upon PostgresSQL
func New() (*GopgStore, error) {
	return &GopgStore{}, nil
}

// Dial perform the connection to the underlying database server
func (d *GopgStore) Dial(options store.Options) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}
	sslMode := true
	if v, ok := options["sslMode"]; !ok || v == "disable" {
		tlsConfig = nil
		sslMode = false
	}
	if v, ok := options["maxConns"]; !ok || v.(int) == 0 {
		options["maxConns"] = 20
	}
	// In minutes, time after which client closes idle connections
	if v, ok := options["idleTimeoutMin"]; !ok || v.(int) == 0 {
		options["idleTimeoutMin"] = 15
	}
	if v, ok := options["host"]; !ok || v.(string) == "" {
		options["host"] = "localhost"
	}
	if v, ok := options["port"]; !ok || v.(int) == 0 {
		options["port"] = 5432
	}

	addr := options["host"].(string) + ":" + strconv.Itoa(options["port"].(int))
	pgOptions := &GOPG.Options{
		TLSConfig:   tlsConfig,
		PoolSize:    options["maxConns"].(int),
		IdleTimeout: time.Duration(options["idleTimeoutMin"].(int)) * time.Minute,
		Addr:        addr,
		User:        options["user"].(string),
		Password:    options["password"].(string),
		Database:    options["dbname"].(string),
	}

	dsn := fmt.Sprintf("user=%s dbname=%s sslmode=%t address=%s", pgOptions.User, pgOptions.Database, sslMode, pgOptions.Addr)
	flog.Info(fmt.Sprintf("connecting to postgresql: %s", dsn))
	d.Pool = GOPG.Connect(pgOptions)

	d.Stat = store.CONNECTED

	s := d.Pool.Options()

	flog.Info(fmt.Sprintf("created GOPG pool store: mc=%d ", s.PoolSize))

	return nil
}

func (d *GopgStore) EnableQueryLogger() {
	GOPG.SetQueryLogger(log.New(os.Stdout, "[GO-PG] ", log.LstdFlags))
}

// Status return the current status of the underlying database
func (d *GopgStore) Status() (int, string) {
	return d.Stat, store.StatusStr[d.Stat]
}

// Close effectively close de database connection
func (d *GopgStore) Close() error {
	flog.Info("GoPg store CLOSING")
	err := d.Pool.Close()
	if err != nil {
		flog.Info(fmt.Sprintf("Error closing DB connection: %s", err))
	}
	d.Stat = store.DISCONNECTED
	return nil
}
