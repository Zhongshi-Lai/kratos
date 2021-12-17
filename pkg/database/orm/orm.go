package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"strings"
	"time"

	"kratos/pkg/ecode"
	"kratos/pkg/log"
	xtime "kratos/pkg/time"

	// database driver
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Config mysql config.
type Config struct {
	DSN         string         // data source name.
	Active      int            // pool
	Idle        int            // pool
	IdleTimeout xtime.Duration // connect max life time.
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info(strings.Repeat("%v ", len(v)), v...)
}

func init() {
	gorm.ErrRecordNotFound = ecode.NothingFound
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *gorm.DB) {
	db, err := gorm.Open(mysql.New(mysql.Config{DSN:c.DSN}), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		log.Error("orm: open error(%v)", err)
		panic(err)
	}
	sql, err := db.DB()
	if err != nil {
		log.Error("mysql: connPool error(%v)", err)
		panic(err)
	}
	sql.SetMaxIdleConns(c.Idle)
	sql.SetMaxOpenConns(c.Active)
	sql.SetConnMaxLifetime(time.Duration(c.IdleTimeout))

	return
}
