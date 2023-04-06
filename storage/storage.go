package storage

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"service_template/logger"
)

type DBID struct {
	ID uint
}

type DBIDs []DBID

func (a DBIDs) toList(ctx context.Context) (ret []uint) {
	log := logger.FromContext(ctx).WithField("m", "toList")
	log.Debugf("toList:: ")

	for _, s := range a {
		ret = append(ret, s.ID)
	}

	return
}

func (a DBIDs) toString(ctx context.Context) (ret string) {
	log := logger.FromContext(ctx).WithField("m", "toString")
	log.Debugf("toString:: ")

	for i, s := range a {
		ret = ret + strconv.Itoa(int(s.ID))
		if i < len(a)-1 {
			ret += ","
		}
	}

	return
}

type Storage struct {
	DB *gorm.DB
}

func (a *Storage) DBLog(ctx context.Context, isSet bool) {
	log := logger.FromContext(ctx).WithField("m", "DBLog")
	log.Debugf("DBLog:: isSet: %v", isSet)

	a.DB.LogMode(true)
}

func (a *Storage) InitPostgress(ctx context.Context,
	host string, port int, dbName string, user string, password string, logOut io.Writer) error {
	log := logger.FromContext(ctx).WithField("m", "InitPostgress")
	log.Debugf("InitPostgress:: host:%v, port:%v", host, port)

	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbName, user, password))
	if err != nil {
		return err
	}
	a.DB = db

	return nil
}
