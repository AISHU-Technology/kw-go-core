package gormx

import (
	"errors"
	"fmt"
	"github.com/AISHU-Technology/kw-go-core/utils"
	"github.com/acmestack/gorm-plus/gplus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
)

type GormConf struct {
	DriverName      string `json:",default=mysql,options=mysql|postgres|sqlserver|clickhouse|oracle|dm8|mariadb|tidb|goldendb"`
	Host            string `json:",default=localhost"` // address
	Port            int    `json:",default=3306"`      // port
	Config          string `json:",optional"`          // extra config such as mysql:charset=utf8mb4&parseTime=True&loc=Local或者postgres:sslmode=disable TimeZone=Asia/Shangh或者clickhouse:read_timeout=10&write_timeout=20 达梦开启compatibleMode=mysql
	DBName          string `json:",default=anydata"`   // orm name
	Username        string `json:",default=root"`      // username
	Password        string `json:",default=root"`      // password
	LogMode         string `json:",default=error"`     // open gorm's global logger
	ConnMaxIdleTime int    `json:",default=100"`       // unit 毫秒
	ConnMaxLifetime int    `json:",default=500"`
	MaxIdleConns    int    `json:",default=10"`
	MaxOpenConns    int    `json:",default=100"`
	TablePrefix     string `json:",optional"` //表前缀
}

/*
GormConf 配置文件
*/
func InitGormConf(e GormConf) *gorm.DB {
	dbType := os.Getenv("DB_TYPE")
	if utils.IsNotBlank(dbType) {
		dbType = strings.ToLower(dbType)
		e.DriverName = dbType
	}
	driver, ok := opens[e.DriverName]
	if !ok {
		panic(errors.New("orm dialect is not supported"))
	}
	dsn := dbHandler(e)
	gdb, _ := gorm.Open(driver(dsn), &gorm.Config{
		PrepareStmt: true,
		QueryFields: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,          //表名后面不加s
			TablePrefix:   e.TablePrefix, // 表前缀
		},
		Logger: logger.New(writer{}, logger.Config{
			SlowThreshold:             2 * time.Second, // 慢 SQL 阈值
			Colorful:                  true,            // Ignore ErrRecordNotFound error for logger
			IgnoreRecordNotFoundError: false,           // Disable color
			LogLevel:                  getLevel(e.LogMode),
		})})
	sqlDb, sqlErr := gdb.DB()
	if sqlErr != nil {
		panic(sqlErr)
	}
	sqlDb.SetMaxIdleConns(e.MaxIdleConns)                             //设置最大的空闲连接数
	sqlDb.SetMaxOpenConns(e.MaxOpenConns)                             //设置最大连接数
	sqlDb.SetConnMaxLifetime(time.Duration(e.ConnMaxLifetime) * 1000) //可重用链接得最大时间
	sqlDb.SetConnMaxIdleTime(time.Duration(e.ConnMaxIdleTime) * 1000) //越短连接过期的次数就会越频繁
	sqlDb.Ping()
	gplus.Init(gdb)
	return gdb
}

func postgresDSN(e GormConf) string {
	var config = "sslmode=disable TimeZone=Asia/Shanghai"
	if utils.IsNotBlank(e.Config) {
		config = e.Config
	}
	dSn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %s", e.Host, e.Username, e.Password, e.DBName, e.Port, config)
	user, _ := os.LookupEnv("RDSUSER")
	if utils.IsNotBlank(user) {
		dSn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s",
			os.Getenv("RDSPORT"),
			os.Getenv("RDSUSER"),
			os.Getenv("RDSPASS"),
			os.Getenv("RDSDBNAME"),
			os.Getenv("RDSHOST"),
			config)
	}
	return dSn
}

func mysqlDSN(e GormConf) string {
	var config = "charset=utf8mb4&parseTime=True&loc=Local"
	if utils.IsNotBlank(e.Config) {
		config = e.Config
	}
	dSn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", e.Username, e.Password, e.Host, e.Port, e.DBName, config)
	user, _ := os.LookupEnv("RDSUSER")
	if utils.IsNotBlank(user) {
		dSn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
			os.Getenv("RDSUSER"),
			os.Getenv("RDSPASS"),
			os.Getenv("RDSHOST"),
			os.Getenv("RDSPORT"),
			os.Getenv("RDSDBNAME"),
			config)
	}
	return dSn
}

func sqlServerDSN(e GormConf) string {
	dSn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?orm=%s", e.Username, e.Password, e.Host, e.Port, e.DBName)
	user, _ := os.LookupEnv("RDSUSER")
	if utils.IsNotBlank(user) {
		dSn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?orm=%s",
			os.Getenv("RDSUSER"),
			os.Getenv("RDSPASS"),
			os.Getenv("RDSHOST"),
			os.Getenv("RDSPORT"),
			os.Getenv("RDSDBNAME"))
	}
	return dSn
}

func clickHouseDSN(e GormConf) string {
	var config = "read_timeout=10&write_timeout=20"
	if utils.IsNotBlank(e.Config) {
		config = e.Config
	}
	dSn := fmt.Sprintf("tcp://%s:%d?orm=%s&username=%s&password=%s&%s", e.Host, e.Port, e.DBName, e.Username, e.Password, config)
	user, _ := os.LookupEnv("RDSUSER")
	if utils.IsNotBlank(user) {
		dSn = fmt.Sprintf("tcp://%s:%s?orm=%s&username=%s&password=%s&%s",
			os.Getenv("RDSHOST"),
			os.Getenv("RDSPORT"),
			os.Getenv("RDSDBNAME"),
			os.Getenv("RDSUSER"),
			os.Getenv("RDSPASS"),
			config)
	}
	return dSn
}

func dmDSN(e GormConf) string {
	var config = "timeout=20s&autocommit=true&readTimeout=10s&genKeyNameCase=2&doSwitch=1"
	if utils.IsNotBlank(e.Config) {
		config = e.Config
	}
	// dm://sysdba:dameng123!@193.100.100.221:5236?autoCommit=true
	dSn := fmt.Sprintf("dm://%s:%s@%s:%d?schema=%s&%s", e.Username, e.Password, e.Host, e.Port, e.DBName, config)
	user, _ := os.LookupEnv("RDSUSER")
	if utils.IsNotBlank(user) {
		dSn = fmt.Sprintf("dm://%s:%s@%s:%s?schema=%s&%s",
			os.Getenv("RDSUSER"),
			os.Getenv("RDSPASS"),
			os.Getenv("RDSHOST"),
			os.Getenv("RDSPORT"),
			os.Getenv("RDSDBNAME"),
			config)
	}
	return dSn
}

func oracleDSN(e GormConf) string {
	// ZTK/sirc1234@193.100.100.43:1521/ORCL
	dSn := fmt.Sprintf("%s/%s@%s:%d/%s", e.Username, e.Password, e.Host, e.Port, e.DBName)
	user, _ := os.LookupEnv("RDSUSER")
	if utils.IsNotBlank(user) {
		dSn = fmt.Sprintf("%s/%s@%s:%s/%s",
			os.Getenv("RDSUSER"),
			os.Getenv("RDSPASS"),
			os.Getenv("RDSHOST"),
			os.Getenv("RDSPORT"),
			os.Getenv("RDSDBNAME"))
	}
	return dSn
}

func tidbDSN(e GormConf) string {
	var config = "charset=utf8mb4&parseTime=True&loc=Local"
	if utils.IsNotBlank(e.Config) {
		config = e.Config
	}
	dSn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", e.Username, e.Password, e.Host, e.Port, e.DBName, config)
	user, _ := os.LookupEnv("RDSUSER")
	if utils.IsNotBlank(user) {
		dSn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
			os.Getenv("RDSUSER"),
			os.Getenv("RDSPASS"),
			os.Getenv("RDSHOST"),
			os.Getenv("RDSPORT"),
			os.Getenv("RDSDBNAME"),
			config)
	}
	return dSn
}

func dbHandler(e GormConf) string {
	var dsn string
	switch e.DriverName {
	case driverMysql:
		dsn = mysqlDSN(e)
	case driverPostgres:
		dsn = postgresDSN(e)
	case driverSqlserver:
		dsn = sqlServerDSN(e)
	case driverClickhouse:
		dsn = clickHouseDSN(e)
	case driveDm:
		dsn = dmDSN(e)
	case driveTiBb:
		dsn = tidbDSN(e)
	case driveOracle:
		dsn = oracleDSN(e)
	default:
		dsn = mysqlDSN(e)
	}
	return dsn
}

func getLevel(logMode string) logger.LogLevel {
	var level logger.LogLevel
	switch logMode {
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	case "error":
		level = logger.Error
	default:
		level = logger.Error
	}
	return level
}

// 自定义一个 Writer
type writer struct {
	logger.Writer
}

func (l writer) Printf(message string, data ...any) {
	log.Println(message, data)
	// TODO上报日志到AR
}
