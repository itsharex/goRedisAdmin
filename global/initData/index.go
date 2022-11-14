package initData

import (
	"errors"
	"fmt"
	"goRedisAdmin/global/global_redis"
	"goRedisAdmin/global/global_write_ip"
	"gopkg.in/ini.v1"
	"log"
	"strings"
)

//IniRead ini reader
var IniRead *ini.File

func init() {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Panicln(err)
	}
	IniRead = cfg
}

func Initialization() {
	global_redis.SetRDConfig(IniRead.Section("redis"))
	//白名单
	initWriteListIp(IniRead.Section("whitelist_ip"))
	fmt.Println(global_write_ip.WriteListIp)
}

func initWriteListIp(config *ini.Section) {
	ips := config.Key("allow_ip").String()
	if ips == "" {
		panic(errors.New("ip白名单为空"))
	}
	ipStr := strings.Split(strings.TrimSpace(ips), ",")
	global_write_ip.WriteListIp = make(map[string]string)
	for _, ip := range ipStr {
		global_write_ip.WriteListIp[ip] = ip
	}
}
