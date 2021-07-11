package drivemysql

import (
	"fmt"
	"log"
	"os"

	config "github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Driversql struct {
	Host         string
	Port         string
	Database     string
	User         string
	Password     string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
	TimeZone     string
	print_log    bool
	Client       *gorm.DB
}
type MysqlDriver struct {
	d      *Driversql
	Client *gorm.DB
}

func (t *Driversql) SetConfig(Host string, Port string, Database string, User string, Password string, Charset string, MaxIdleConns int, MaxOpenConns int, print_log bool) {
	t.Host = Host
	t.Port = Port
	t.Database = Database
	t.User = User
	t.Password = Password
	t.Charset = Charset
	t.MaxIdleConns = MaxIdleConns
	t.MaxOpenConns = MaxOpenConns
	t.print_log = print_log
}
func (t *Driversql) Configure(ConfigPath string, ConfigName string) {
	config.AddConfigPath(ConfigPath)
	config.SetConfigName(ConfigName)
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if config.GetString("default.host") == "" {
		fmt.Println("falta host en el archivo " + ConfigName)
		os.Exit(1)
	}
	if config.GetString("default.database") == "" {
		fmt.Println("falta database en el archivo de config " + ConfigName)
		os.Exit(1)
	}
	if config.GetString("default.user") == "" {
		fmt.Println("falta user en el archivo de config " + ConfigName)
		os.Exit(1)
	}
	//set the config
	t.SetConfig(config.GetString("default.host"),
		config.GetString("default.port"),
		config.GetString("default.database"),
		config.GetString("default.user"),
		config.GetString("default.password"),
		config.GetString("default.charset"),
		config.GetInt("default.MaxIdleConns"),
		config.GetInt("default.MaxOpenConns"),
		config.GetBool("default.sql_log"))

	/*t.Host = config.GetString("default.host")
	t.Port = config.GetString("default.port")
	t.Database = config.GetString("default.database")
	t.User = config.GetString("default.user")
	t.Password = config.GetString("default.password")
	t.Charset = config.GetString("default.charset")
	t.MaxIdleConns = config.GetInt("default.MaxIdleConns")
	t.MaxOpenConns = config.GetInt("default.MaxOpenConns")
	t.print_log = config.GetBool("default.sql_log")*/
}

func (t *MysqlDriver) GetClient() *gorm.DB {
	return t.Client
}

func (t *MysqlDriver) ConfigureMySqlForFile(ConfigPath string, ConfigName string) error {
	t.configureForFile(ConfigPath, ConfigName)
	return t.newDriverMysql()
}
func (t *MysqlDriver) configureForFile(ConfigPath string, ConfigName string) {
	x := &Driversql{}
	x.Configure(ConfigPath, ConfigName)
	t.d = x
}
func (c *MysqlDriver) newDriverMysql() error {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", c.d.User, c.d.Password, c.d.Host, c.d.Port, c.d.Database, c.d.Charset)
	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		log.Panic(err)
		return err
	}
	c.Client = db
	return nil
}
