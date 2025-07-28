package zapgorm

import (
	"context"
	//"database/sql"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"

	
)

type ZapLogger struct {
	LogLevel      logger.LogLevel
	Logger        *zap.Logger
	SlowThreshold time.Duration
}

// Error implements logger.Interface.
func (z *ZapLogger) Error(ctx context.Context,msg string,args ...interface{}) {
	//panic("unimplemented")
	if z.LogLevel < logger.Error {
        return
    }
    z.Logger.Sugar().Errorw(msg,args...)
}

// Info implements logger.Interface.
func (z *ZapLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	//panic("unimplemented")
	if z.LogLevel< logger.Info {
		return
	}
	z.Logger.Sugar().Infow(msg, args...)
}

func (z *ZapLogger) Warn(ctx context.Context,msg string,arg ...interface{}) {
	//panic("unimplemented")
	if z.LogLevel< logger.Warn {
		return
	}
	z.Logger.Sugar().Warnw(msg,arg...)
}

// Trace implements logger.Interface.
func (z *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	//panic("unimplemented")
    elapsed := time.Since(begin)
	sql, rowsAffected := fc()

	if err != nil && z.LogLevel >= logger.Error {
			z.Logger.Sugar().Errorw("gorm trace error",
				"elapsed", elapsed,
				"sql", sql,
				"rowsAffected", rowsAffected,
				"error", err,
			)
			return
		}
		
	
	if elapsed> z.SlowThreshold && z.SlowThreshold != 0 {
		if z.LogLevel >= logger.Warn {
			z.Logger.Sugar().Warnw("gorm trace slow",
				"elapsed", elapsed,
				"sql", sql,
				"rowsAffected", rowsAffected,
			)
		}
		return
	}

 if z.LogLevel >= logger.Info {
		z.Logger.Sugar().Infow("gorm trace",
			"elapsed", elapsed,
			"sql", sql,
			"rowsAffected", rowsAffected,
		)
	}

}

// Warn implements logger.Interface.


//constructor

func NewZapLogger(zapLogger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		Logger:        zapLogger,
		LogLevel:      logger.Info, // default log level
		SlowThreshold: 200 * time.Millisecond,
	}
}

func (z *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *z
	newLogger.LogLevel = level
	return &newLogger
}
