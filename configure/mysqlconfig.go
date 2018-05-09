package configure

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //init mysql
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	// DBAddr mysql地址
	DBAddr = "127.0.0.1:3306"
	// DBUser 用户名
	DBUser = "tiger"
	// DBPasswd 密码
	DBPasswd = "tigerisnotcat"
	// DBName 数据库名
	DBName = "neochain"
	// DBMaxIdleTime 数据库连接池最大保持空闲时间ms
	DBMaxIdleTime = 1000
	// DBMaxIdle 数据库连接池最大保持空闲个数
	DBMaxIdle = 1
	// DBMaxOverflow 最大容许溢出数量
	DBMaxOverflow = 50
)

//InitMysql mysql的初始化
func InitMysql() (*sqlx.DB, error) {

	dbURI := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8",
		DBUser,
		DBPasswd,
		DBAddr,
		DBName)

	DB, err := sqlx.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	DB.SetConnMaxLifetime(time.Millisecond * time.Duration(DBMaxIdleTime))
	DB.SetMaxIdleConns(DBMaxIdle)
	DB.SetMaxOpenConns(DBMaxOverflow + DBMaxIdle)

	err = DB.Ping()
	if err != nil {
		return nil, err
	}
	return DB, nil
}
