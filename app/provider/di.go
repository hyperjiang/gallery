package provider

import (
	"log"
	"sync"

	"github.com/hyperjiang/gallery-service/app/config"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Container - the structure of dependency injecter
type Container struct {
	store  sync.Map
	onceDB sync.Once
}

// DI keys
const (
	KeyConfig = "config"
	KeyDB     = "db"
	KeyLog    = "log"
)

// the global DI container
var di *Container

// New - return a new instance of Container
func New() *Container {
	return &Container{
		store: sync.Map{},
	}
}

// DI - return the global DI container
func DI() *Container {
	if di == nil {
		di = New()
	}
	return di
}

// Config - get global configs
func (di *Container) Config() *config.Config {
	if v, ok := di.store.Load(KeyConfig); ok {
		return v.(*config.Config)
	}

	config := config.LoadAll()
	di.store.Store(KeyConfig, config)
	return config
}

// Log - get a logger
func (di *Container) Log() *zap.SugaredLogger {
	if v, ok := di.store.Load(KeyLog); ok {
		return v.(*zap.SugaredLogger)
	}

	logger := Logger([]string{"stderr", di.Config().Server.LogDir + "/service.log"})
	di.store.Store(KeyLog, logger)
	return logger
}

// InitDB - initialize DB connections
func (di *Container) InitDB() error {
	var err error

	di.onceDB.Do(func() {
		var db *sqlx.DB
		db, err = openDB(&di.Config().Database.Main)
		if err != nil {
			log.Println("Fail to open DB")
			return
		}
		di.store.Store(KeyDB, db)
	})

	return err
}

// CloseDB - close DB connections
func (di *Container) CloseDB() {
	db := di.DB()
	if db != nil {
		db.Close()
	}
}

// Bind - bind a value to a key
func (di *Container) Bind(k interface{}, v interface{}) *Container {
	di.store.Store(k, v)
	return di
}

// Unbind - unbind a value from a key
func (di *Container) Unbind(k interface{}) *Container {
	di.store.Delete(k)
	return di
}

// DB - get the  db connection
func (di *Container) DB() *sqlx.DB {
	client, ok := di.store.Load(KeyDB)
	if client == nil || !ok {
		return nil
	}
	return client.(*sqlx.DB)
}
