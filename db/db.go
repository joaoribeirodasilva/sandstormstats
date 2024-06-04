package db

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joaoribeirodasilva/sandstormstats/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Db struct {
	Conn   *gorm.DB
	logger *slog.Logger
	config *conf.DbConf
}

func New(conf *conf.DbConf, logger *slog.Logger) *Db {
	d := new(Db)
	d.logger = logger
	d.config = conf
	return d
}

func (d *Db) Connect() {

	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		d.config.User,
		d.config.Password,
		d.config.Host,
		d.config.Port,
		d.config.Database,
		d.config.Extra,
	)

	d.logger.Debug(fmt.Sprintf("connecting to database '%s' at '%s:%d'", d.config.Database, d.config.Host, d.config.Port))

	d.Conn, err = gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        dsn,
	}), &gorm.Config{})

	if err != nil {
		d.logger.Error(fmt.Sprintf("failed to connecto to database '%s' at '%s:%d'", d.config.Database, d.config.Host, d.config.Port))
		os.Exit(1)
	}

	d.logger.Info(fmt.Sprintf("database '%s' successfully", d.config.Database))
}
