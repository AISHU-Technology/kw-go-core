package gnorm

type GnConf struct {
	Host         []string `json:",default=127.0.0.1:9669"`
	Username     string   `json:",default=root"`   // username
	Password     string   `json:",default=nebula"` // password
	MinConnsSize int      `json:",default=5"`      // unit 毫秒
	MaxConnSize  int      `json:",default=10"`
	Timeout      int      `json:",default=300000"`
	IdleTime     int      `json:",default=180000"`
}

//func InitGnConf(e GnConf) *norm.DB {
//dalector := dialectors.MustNewNebulaDialector(dialectors.DialectorConfig{
//	Addresses: []string{"127.0.0.1:9669"},
//	Timeout:   time.Second * 5,
//	Space:     "test",
//	Username:  "test",
//	Password:  "test",
//})
//db := norm.MustOpen(dalector, norm.Config{})
//return db
//}
