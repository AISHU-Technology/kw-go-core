package gormx

import (
	"gitee.com/tdxmkf123/gorm-driver-dameng/dameng"
	"gitee.com/tdxmkf123/gorm-driver-oracle/oracle"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	driverMysql      string = "mysql"
	driverPostgres   string = "postgres"
	driverSqlserver  string = "sqlserver"
	driverClickhouse string = "clickhouse"
	driveDm          string = "dm8"
	driveMariaBb     string = "mariadb"
	driveTiBb        string = "tidb"
	driveGoldenDB    string = "goldendb"
	driveOracle      string = "oracle"
)

var opens = map[string]func(string) gorm.Dialector{
	driverMysql:      mysql.Open,
	driveMariaBb:     mysql.Open,
	driveTiBb:        mysql.Open,
	driveGoldenDB:    mysql.Open,
	driverPostgres:   postgres.Open,
	driverSqlserver:  sqlserver.Open,
	driverClickhouse: clickhouse.Open,
	driveOracle:      oracle.Open,
	driveDm:          dameng.Open,
}
