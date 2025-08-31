package database

import (
	"log"
	"os"

	"github.com/chaitu35/claimeasy/zapgorm"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB        *gorm.DB
	zapLogger *zap.Logger
)

func getLogLevelByEnv(env string) logger.LogLevel {
	switch env {
	case "development":
		return logger.Info
	case "test":
		return logger.Warn
	case "production":
		return logger.Error
	default:
		return logger.Silent
	}
}

func initZapLogger(env string) *zap.Logger {
	var err error
	if env == "development" || env == "test" {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("failed to create zap logger: %v", err)
	}
	return zapLogger
}

func Connection() {
	// 1. Detect environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// 2. Load .env file dynamically
	// (optional if config already loaded, but safe)
	// _ = godotenv.Load(".env." + env)

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	zapLogger = initZapLogger(env)
	logLevel := getLogLevelByEnv(env)
	gormLogger := zapgorm.NewZapLogger(zapLogger).LogMode(logLevel)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		zapLogger.Fatal("failed to connect to database", zap.Error(err))
	}

	DB = db
	zapLogger.Info("connected to database", zap.String("dsn", dsn), zap.String("env", env))
	_ = zapLogger.Sync()
}
