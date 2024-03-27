package orm

import (
	"fmt"
	"gitee.com/tdxmkf123/gorm-driver-dameng/dameng"
	"github.com/AISHU-Technology/kw-go-core/orm/gormx"
	"testing"
	"time"
)

type CompanyInfo struct {
	CompanyId   string        `gorm:"column:COMPANY_ID"`
	CnName      string        `gorm:"column:CN_NAME"`
	ManageRange dameng.MyClob `gorm:"column:MANAGE_RANGE"`
}

type User struct {
	ID        string    `gorm:"type:varchar(50);primary_key"`
	UserName  string    `gorm:"column:USER_NAME;size:60;not null"`
	Age       int64     `gorm:"column:F_AGE;size:10;not null"`
	Email     string    `gorm:"column:EMAIL;unique_index"`
	Content   string    `gorm:"column:f_content;type:varchar(255)"`
	Address   string    `gorm:"column:ADDRESS"`
	Remark    string    `gorm:"column:f_remark;size:2000"`
	CreatedAt time.Time `gorm:"column:CREATED_AT;type:DATETIME"`
}

func TestDm8(t *testing.T) {
	o := gormx.GormConf{
		DriverName:  "dm8",
		Host:        "10.2.174.222",
		Port:        52360,
		Config:      "compatibleMode=mysql",
		Username:    "anydata",
		Password:    "anydata001",
		TablePrefix: "ANYDATA.t_",
		LogMode:     "warn",
	}
	ormx := gormx.InitGormConf(o)
	ormx.AutoMigrate(&User{})
	var companyList []User
	//err = ormx.Raw("SELECT COMPANY_ID,CN_NAME,MANAGE_RANGE FROM DW.COMPANY_INFO WHERE COMPANY_ID in('C3DDBD2F17554E8A838DB706C139B883') ").Scan(&companyList).Error //MANAGE_RANGE ,'4A5A6EA9B47445D48CB30683BEE68C4A'
	//err = ormx.Raw("select t.ID,t.Title,t.Content,t.rowid from dw.table1 t where id=1 ").Scan(&companyList).Error //MANAGE_RANGE ,'4A5A6EA9B47445D48CB30683BEE68C4A'
	// 读取
	//ormx.Select("COMPANY_ID,CN_NAME").Where("COMPANY_ID = ?", "C3DDBD2F17554E8A838DB706C139B888").Find(&companyList) // 查询id为1的product
	//ormx.First(&companyList, "COMPANY_ID = ?", "C3DDBD2F17554E8A838DB706C139B883").Select("ID") // 查询code为l1212的product
	//err = ormx.Exec("update  dw.table1 set title='标题333',content='内容333' where id=3").Error
	ormx.Select("*").Where("id = ?", "1").Find(&companyList) // 查询id为1的product
	for i, v := range companyList {
		fmt.Printf("====输出==i==%d,输出值id==%s==Name==%s==CreatedAt==%v===Remark===%s===Address==%s", i, v.ID, v.UserName, v.CreatedAt, v.Remark, v.Address)
	}
}

func TestMysql(t *testing.T) {
	o := gormx.GormConf{
		DriverName: "mariadb",
		Host:       "10.4.71.138",
		Port:       3330,
		//Config:     "compatibleMode=mysql",
		DBName:      "anydata",
		Username:    "anyshare",
		Password:    "eisoo.com123",
		TablePrefix: "anydata.t_",
	}
	ormx := gormx.InitGormConf(o)
	ormx.AutoMigrate(&User{})
	var companyList []User
	//err = ormx.Raw("SELECT COMPANY_ID,CN_NAME,MANAGE_RANGE FROM DW.COMPANY_INFO WHERE COMPANY_ID in('C3DDBD2F17554E8A838DB706C139B883') ").Scan(&companyList).Error //MANAGE_RANGE ,'4A5A6EA9B47445D48CB30683BEE68C4A'
	//err = ormx.Raw("select t.ID,t.Title,t.Content,t.rowid from dw.table1 t where id=1 ").Scan(&companyList).Error //MANAGE_RANGE ,'4A5A6EA9B47445D48CB30683BEE68C4A'
	// 读取
	//ormx.Select("COMPANY_ID,CN_NAME").Where("COMPANY_ID = ?", "C3DDBD2F17554E8A838DB706C139B888").Find(&companyList) // 查询id为1的product
	//ormx.First(&companyList, "COMPANY_ID = ?", "C3DDBD2F17554E8A838DB706C139B883").Select("ID") // 查询code为l1212的product
	//err = ormx.Exec("update  dw.table1 set title='标题333',content='内容333' where id=3").Error
	ormx.Select("*").Where("ID = ?", "1").Find(&companyList) // 查询id为1的product
	for i, v := range companyList {
		fmt.Printf("====输出==i==%d,输出值id==%s==Name==%s==CreatedAt==%v===Remark===%s===Address==%s", i, v.ID, v.UserName, v.CreatedAt, v.Remark, v.Address)
	}
}
