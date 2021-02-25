package public

import (
	"fmt"
	"testing"
)

const (
	STATUS_WAIT = 0
	STATUS_EXIT = 1
)

var sch =make(chan string,2)
func Test_RD(t *testing.T)  {
	redisfile := "./conf/redis.yaml"

	v := ViperConfig(redisfile)
	fmt.Println(v.Get("redis"))


	client,err := NewRedis(SV.Redis)
	if err!=nil{
		t.Log(err)
	}

	keyid := "20"
	//
	client.Rpush(keyid,"a1",0)
	client.Rpush(keyid,"a2",0)
	client.Rpush(keyid,"a4",0)
	client.Rpush(keyid,"a5",0)
	client.Rpush(keyid,"a6",0)
	//
	//
	data,err :=client.Lrange(keyid,0,-1)
	if err!=nil{
		t.Error(err)
	}
	t.Log(data)

}
