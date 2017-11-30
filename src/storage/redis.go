package storage

import (
	"github.com/mediocregopher/radix.v2/redis"
	"log"
)

type storageRedis struct {

	network string
	addr    string
	prefix  string
	storage *redis.Client
}

func NewStorageRedis(network, addr, prefix string) *storageRedis{

	s := &storageRedis{network:network, addr:addr, prefix:prefix}
	s.Connect()

	return s
}

func (f *storageRedis) Connect() (isConnected bool){

	if !f.isConnected(){

		if conn, err := redis.Dial(f.network, f.addr) ; err == nil{
			f.storage = conn
		}else{
			log.Panicln(err)
		}
	}

	isConnected = f.isConnected()

	log.Println("[storage.redis]Connect: ", f.isConnected())

	return isConnected
}

func (f *storageRedis) Disconnect(){
	f.storage.Close()
}

func (f *storageRedis) isConnected() bool{
	return f.storage != nil
}


func (f *storageRedis) Get(index string) (record string){

	record,_ = f.storage.Cmd("HGET", f.prefix, index).Str()
	//if str,err := f.storage.Cmd("HGET", storageKey, index).Str(); err != redis.ErrRespNil{
	//	record = Record(str)
	//}

	return record
}

func (f *storageRedis) Add(index string, record string, force bool) (result bool, err error){

	log.Println("[storage.redis]Add: ", index, record, force)

	if f.Get(index) == "" || force{

		var resp int
		resp, err = f.storage.Cmd("HSET", f.prefix, index, record).Int()

		result = resp > 0

	}else{
		err = errIndexExist
	}

	log.Println("[storage.redis]Add: ", "result:", result)
	log.Println("[storage.redis]Add: ", "error:", err)

	return result, err
}

func (f *storageRedis) Rm(index string) (result bool, err error){

	log.Println("[storage.redis]Rm: ", index)

	if r, err := f.storage.Cmd("HDEL", f.prefix, index).Int(); err == nil {

		result = r>0
	}

	log.Println("[storage.redis]Rm: ", "result:", result)
	log.Println("[storage.redis]Rm: ", "error:", err)

	return result, err
}

func (f *storageRedis) List() (data map[string]string){

	data, _ = f.storage.Cmd("HGETALL", f.prefix).Map()

	return data
}

func (f *storageRedis) Flush(){
	f.storage.Cmd("DEL", f.prefix)
}