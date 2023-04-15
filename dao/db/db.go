package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DbIns *DB

type DB struct {
	*gorm.DB
}

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s&timeout=%s&readTimeout=%s&writeTimeout=%s",
		conf.Handler.GetString("mysql.auth"),
		conf.Handler.GetString("mysql.password"),
		conf.Handler.GetString("mysql.addr"),
		conf.Handler.GetString("mysql.database"),
		conf.Handler.GetString("mysql.charset"),
		conf.Handler.GetBool("mysql.parseTime"),
		conf.Handler.GetString("mysql.loc"),
		conf.Handler.GetString("mysql.timeout"),
		conf.Handler.GetString("mysql.readTimeout"),
		conf.Handler.GetString("mysql.writeTimeout"),
	)
	fmt.Println(dsn)
	RAILWAY_ENVIRONMENT := os.Getenv("RAILWAY_ENVIRONMENT")
	fmt.Println("RAILWAY_ENVIRONMENT", RAILWAY_ENVIRONMENT)
	MYSQL_URL := os.Getenv("MYSQL_URL")
	MYSQLDATABASE := os.Getenv("MYSQLDATABASE")
	MYSQLHOST := os.Getenv("MYSQLHOST")
	MYSQLPASSWORD := os.Getenv("MYSQLPASSWORD")
	MYSQLPORT := os.Getenv("MYSQLPORT")
	MYSQLUSER := os.Getenv("MYSQLUSER")
	fmt.Println("MYSQL_URL, MYSQLDATABASE, MYSQLHOST, MYSQLPASSWORD, MYSQLPORT, MYSQLUSER", MYSQL_URL, MYSQLDATABASE, MYSQLHOST, MYSQLPASSWORD, MYSQLPORT, MYSQLUSER)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.LogLevel(conf.Handler.GetInt("mysql.log_mode")))})
	if err != nil {
		tlog.Handler.Fatalf(nil, consts.SLTagMysqlFail, "DB init fail")
		log.Panicf("DB mysql init fail : %s", err)
	}
	sqldb, _ := db.DB()
	sqldb.SetMaxIdleConns(conf.Handler.GetInt("mysql.max_idle_conns"))
	sqldb.SetMaxOpenConns(conf.Handler.GetInt("mysql.max_open_conns"))
	sqldb.SetConnMaxLifetime(time.Duration(conf.Handler.GetInt("mysql.conn_max_lifetime")) * time.Minute)
	cancelPing, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	sqldb.PingContext(cancelPing)
	// 自动建表改表
	db.AutoMigrate(&Permission{})
	db.AutoMigrate(&RechargeKey{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&RolePermission{})
	db.AutoMigrate(&ThirdAuth{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserIntegral{})
	db.AutoMigrate(&UserRole{})
	db.AutoMigrate(&VerificationCode{})
	db.AutoMigrate(&InviteRelation{})
	db.AutoMigrate(&Feedback{})
	DbIns = &DB{
		db,
	}
}

func GetDbIns(db *gorm.DB) *gorm.DB {
	if db == nil {
		return DbIns.DB
	} else {
		return db
	}
}
