package data

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"undersea/im-api/conf"
)

func NewMysql(conf conf.Conf) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s", conf.MySQL.Username,
		conf.MySQL.Password, conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.DbName, conf.MySQL.Charset, conf.MySQL.Zone)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("NewMysql->gorm.Open err,%v", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		err = fmt.Errorf("NewMysql->db.DB() err,%v", err)
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		err = fmt.Errorf("NewMysql->sqlDB.Ping() err,%v", err)
		return
	}

	sqlDB.SetMaxOpenConns(conf.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(conf.MySQL.MaxIdleConns)
	return
}
