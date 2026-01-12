package shared

import (
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func Connect(url string, logger *zap.Logger) *sqlx.DB {
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		logger.Sugar().Fatalf("db open failed: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Sugar().Fatalf("db pinged: %v", err)
	}

	logger.Sugar().Info("DB connected successful")
	return db
}
