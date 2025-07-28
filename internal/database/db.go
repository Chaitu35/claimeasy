package database

import (
	"log"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/chaitu35/claimeasy/zapgorm"
	
)

var (
	DB *gorm.DB
	zapLogger *zap.Logger
	
)

func getLogLevelByEnv(env string) logger.LogLevel{
	switch env{
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

func initZapLogger(env string) *zap.Logger{
 // var logger *zap.Logger
  var err error
 if(env == "development" || env == "test") {
		zapLogger,err= zap.NewDevelopment()
	}else{
	     zapLogger,err= zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("failed to create zap logger: %v", err)
	}
	return zapLogger

}

func Connection(){
     dsn := os.Getenv("DB_DSN")
	 env:= os.Getenv("ENV")
	 if dsn == "" {	
		log.Fatal("DB_DSN environment variable is not set")
	 }	
	 if env == "" {
		log.Fatal("ENV environment variable is not set")	
	}
	zapLogger = initZapLogger(env)
	logLevel:= getLogLevelByEnv(env)
	gormLogger := zapgorm.NewZapLogger(zapLogger).
	LogMode(logLevel)



	 
	//defer zapLogger.Sync()




	
	db,err:=gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
if err != nil {
	zapLogger.Fatal("failed to connect to database", zap.Error(err))
}
DB=db
zapLogger.Info("connected to database", zap.String("dsn", dsn), zap.String("env", env))
_=zapLogger.Sync() // flushes buffer, if any
}