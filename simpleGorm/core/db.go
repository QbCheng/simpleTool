package core

import (
	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"simpleTool/simpleGorm/customLogger"
	"sync"
	"time"
)

type DB struct {
	GormDB *gorm.DB
}

var (
	once     sync.Once
	defeatDB *DB
)

func GetDB() *gorm.DB {
	if defeatDB == nil {
		once.Do(func() {
			config := mysql2.Config{
				User:                 "root",
				Passwd:               "root",
				Net:                  "tcp",
				Addr:                 "192.168.31.202:3306",
				DBName:               "test",
				Loc:                  time.Local,
				Collation:            "utf8mb4_general_ci",
				MaxAllowedPacket:     4 << 20,
				CheckConnLiveness:    true,
				ParseTime:            true,
				AllowNativePasswords: true,
			}
			defeatDB = NewDB(config)
		})
	}
	return defeatDB.GormDB
}

func NewDB(config mysql2.Config) *DB {
	ret := &DB{}
	dsn := config.FormatDSN()
	var err error
	ret.GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: customLogger.NewLogger(customLogger.Config{
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: false,
		}),
	})
	if err != nil {
		return nil
	}
	return ret
}

func Migrator(models ...interface{}) error {
	return GetDB().Migrator().AutoMigrate(models...)
}
