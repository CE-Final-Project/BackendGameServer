package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Config struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbName"`
	Password string `yaml:"password"`
}

const (
	maxConn         = 50
	maxConnIdleTime = 1 * time.Minute
	maxConnLifetime = 3 * time.Minute
)

func NewPostgresDatabase(cfg *Config) (*sqlx.DB, error) {

	//dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?sslmode=disable&parseTime=true",
	//	cfg.User,
	//	cfg.Password,
	//	cfg.Host,
	//	cfg.Port,
	//	cfg.DBName,
	//)

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	db, err := sqlx.Connect(cfg.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConn)
	db.SetConnMaxIdleTime(maxConnIdleTime)
	db.SetConnMaxLifetime(maxConnLifetime)

	return db, nil
}
