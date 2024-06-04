package conf

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	default_log_path    = "./log"
	default_db_host     = "localhost"
	default_db_port     = 3306
	default_db_database = "sandstorm_stats"
	default_db_user     = ""
	default_db_password = ""
)

type LogConf struct {
	Path string
}

type DbConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Extra    string
}

type Conf struct {
	log    LogConf
	db     DbConf
	logger *slog.Logger
}

func New(logger *slog.Logger) *Conf {

	c := new(Conf)
	c.logger = logger
	c.log.Path = default_log_path
	c.db.Host = default_db_host
	c.db.Port = default_db_port
	c.db.Database = default_db_database
	c.db.User = default_db_user
	c.db.Password = default_db_password

	return c
}

func (c *Conf) Read() {

	godotenv.Load()

	temp := os.Getenv("LOG_DIRECTORY")
	c.log.Path = strings.TrimSpace(temp)

	temp = os.Getenv("DATABASE_HOST")
	temp = strings.TrimSpace(temp)
	if temp != "" {
		c.db.Host = temp
	}

	temp = os.Getenv("DATABASE_PORT")
	temp = strings.TrimSpace(temp)
	tempI, err := strconv.ParseInt(temp, 10, 32)
	if err == nil && tempI > 0 {
		c.db.Port = int(tempI)
	}

	temp = os.Getenv("DATABASE_NAME")
	temp = strings.TrimSpace(temp)
	if temp != "" {
		c.db.Database = temp
	}

	temp = os.Getenv("DATABASE_USER")
	temp = strings.TrimSpace(temp)
	if temp != "" {
		c.db.User = temp
	}

	temp = os.Getenv("DATABASE_PASSWORD")
	temp = strings.TrimSpace(temp)
	if temp != "" {
		c.db.Password = temp
	}

	temp = os.Getenv("DATABASE_EXTRA")
	temp = strings.TrimSpace(temp)
	if temp != "" {
		c.db.Extra = temp
	}
}

func (c *Conf) GetLogConf() *LogConf {

	return &c.log
}

func (c *Conf) GetDbConf() *DbConf {

	return &c.db
}
