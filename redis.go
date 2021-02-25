package public

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)


type Redis struct {
	DB 			int			`yaml:"db"`
	Host 		string		`yaml:"host"`
	Port 		int 		`yaml:"port"`
	Password 	string		`yaml:"password"`
}

type RedisCl struct {
	Client *redis.Client
}

//实例化配置对象
func NewRedisCfg(path string)  {



}



//todo redis服务对象操作
func NewRedis(rc *Redis) (cli *RedisCl,err error) {
	rdhost := fmt.Sprintf("%s:%d",rc.Host,rc.Port)
	fmt.Printf("连接redis:%s,",rdhost)
	client := redis.NewClient(&redis.Options{
		Addr:     rdhost,
		Password: rc.Password, // no password set
		DB:       rc.DB,  // use default DB
		DialTimeout: time.Duration(time.Second*2),
	})
	_, err = client.Ping().Result()
	if err!=nil{
		fmt.Println(err.Error())
		return nil,err
	}
	clili :=&RedisCl{Client: client}
	return clili,nil
}

//转换过期时间
func (cli *RedisCl) makeExpTime(t int) (dt time.Duration) {

	if t !=0{
		dt = time.Duration(t) * time.Second
	}else {
		dt = time.Duration(REDIS_EXPTIME_DEFAULT)*time.Second
	}
	return dt
}

//set数据
func (cli *RedisCl) Set(k string,v interface{},ex int) (err error)  {
	expTime := cli.makeExpTime(ex)
	_,err = cli.Client.Set(k,v,expTime).Result()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	return nil
}

//get数据
func (cli *RedisCl) Get(k string) (data string,err error)  {
	data,err = cli.Client.Get(k).Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

//查询key是否存在
func (cli *RedisCl) keyExists(k string) bool  {
	_,err := cli.Client.Get(k).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		return true
	}
}

//删除key
func (cli *RedisCl) Del(k string) error  {
	_, err := cli.Client.Del(k).Result()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//rpush往列表push数据
func (cli *RedisCl) Rpush(k string,data interface{},exp int) (err error)  {
	_,err = cli.Client.RPush(k,data).Result()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	expTime := cli.makeExpTime(exp)
	cli.Client.Expire(k,expTime)
	return nil
}

//lrange查询列表数据
func (cli *RedisCl) Lrange(key string, start, stop int64) (data []string,err error)  {
	data,err = cli.Client.LRange(key,start,stop).Result()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	return
}

//删除列表数据
func (cli *RedisCl) LTrim(key string, start, stop int64)  {
	cli.Client.LTrim(key,start, stop)
}


//解密redis密码
func (rc *Redis) decryptPass() (dec string,err error) {
	aesAuth := NewAesAuth()
	dec,err = aesAuth.AesDecryptCFB(rc.Password)
	if err!=nil{
		return
	}
	return
}

//加密redis密码加密配置密码
func (rc *Redis) encryptPass() (err error)  {
	aesAuth := NewAesAuth()
	if rc.Password != ""{
		_,err := aesAuth.AesDecryptCFB(rc.Password)
		if err!=nil{
			fmt.Printf("加密明文")
			enc := aesAuth.AesEncryptCFB([]byte(rc.Password))
			rc.Password = enc
			return nil
		}
		return nil
	}
	return
}
