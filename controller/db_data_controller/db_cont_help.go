package db_data_controller

import (
	"errors"
	"github.com/go-redis/redis"
	"goRedisAdmin/global/global_redis"
	"time"
)

type DbDataHelpModel struct {
	DbNum  int    `json:"db_num"`
	VType  string `json:"type"`
	Key    string `json:"key"`
	Val    string `json:"val"`
	Expire int    `json:"expire"`
}

type DbDataHelpCont struct {
	redis *redis.Client
	DbDataHelpModel
}

func NewDbDataHelpController(v *DbDataHelpModel) (*DbDataHelpCont, error) {
	if v.DbNum < 0 || v.DbNum > 15 {
		return nil, errors.New("dbNum must be between 0 and 15")
	}
	cont := new(DbDataHelpCont)
	cont.DbDataHelpModel = *v
	rd, err := cont.getRedisClient()
	if err != nil {
		return nil, err
	}
	cont.redis = rd
	return cont, nil
}

func (d *DbDataHelpCont) getRedisClient() (*redis.Client, error) {
	return global_redis.GetRedisClient(d.DbNum)
}

func (d *DbDataHelpCont) CloseClient() {
	_ = d.redis.Close()
}

func (d *DbDataHelpCont) AddString() error {
	_, err := d.redis.Set(d.Key, d.Val, time.Duration(int64(d.Expire))*time.Second).Result()
	return err
}

func (d *DbDataHelpCont) AddList() error {
	_, err := d.redis.LPush(d.Key, d.Val).Result()
	if err != nil {
		return err
	}
	return setExpire(d)
}

func (d *DbDataHelpCont) AddSet() error {
	_, err := d.redis.SAdd(d.Key, d.Val).Result()
	if err != nil {
		return err
	}
	return setExpire(d)
}

func setExpire(d *DbDataHelpCont) error {
	if d.Expire > 0 {
		_, err := d.redis.Expire(d.Key, time.Duration(int64(d.Expire))*time.Second).Result()
		return err
	} else {
		return nil
	}
}
