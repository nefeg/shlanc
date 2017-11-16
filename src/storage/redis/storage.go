package redis

import (
	"errors"
	"github.com/mediocregopher/radix.v2/redis"
	"log"
)

type storage struct {

	network string
	addr    string
	storage *redis.Client
}

var storageKey      = "hrentab"
var errIndexExist   = errors.New("index already exist")

func NewRedisStorage(network, addr string) *storage{

	s := &storage{network:network, addr:addr}
	s.Connect()

	return s
}

func (f *storage) Connect() (isConnected bool){

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

func (f *storage) Disconnect(){
	f.storage.Close()
}

func (f *storage) isConnected() bool{
	return f.storage != nil
}


func (f *storage) Get(index string) (record string){

	record,_ = f.storage.Cmd("HGET", storageKey, index).Str()
	//if str,err := f.storage.Cmd("HGET", storageKey, index).Str(); err != redis.ErrRespNil{
	//	record = Record(str)
	//}

	return record
}

func (f *storage) Add(index string, record string, force bool) (result bool, err error){

	log.Println("[storage.redis]Add: ", index, record, force)

	if f.Get(index) == "" || force{

		var resp int
		resp, err = f.storage.Cmd("HSET", storageKey, index, record).Int()

		result = resp > 0

	}else{
		err = errIndexExist
	}

	log.Println("[storage.redis]Add: ", "result:", result)
	log.Println("[storage.redis]Add: ", "error:", err)

	return result, err
}

func (f *storage) Rm(index string) (result bool, err error){

	log.Println("[storage.redis]Rm: ", index)

	if r, err := f.storage.Cmd("HDEL", storageKey, index).Int(); err == nil {

		result = r>0
	}

	log.Println("[storage.redis]Rm: ", "result:", result)
	log.Println("[storage.redis]Rm: ", "error:", err)

	return result, err
}

func (f *storage) List() (data map[string]string){

	data, _ = f.storage.Cmd("HGETALL", storageKey).Map()

	return data
}

func (f *storage) Flush(){
	f.storage.Cmd("DEL", storageKey)
}