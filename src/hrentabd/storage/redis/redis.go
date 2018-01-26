package redis

import (
	"github.com/mediocregopher/radix.v2/redis"
	"log"
	"strconv"
	"errors"
)

var errIndexExist = errors.New("index already exist")

type storageRedis struct {

	network     string
	addr        string
	storageKey  string
	storageLock string
	storageVer  string

	storage     *redis.Client
}

func New(network, addr, storageKey string) *storageRedis{

	s := &storageRedis{network:network, addr:addr}
	s.storageKey    = storageKey + ".db"
	s.storageLock   = s.storageKey + ".lock"
	s.storageVer    = s.storageKey + ".version"

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

	var version string
	if version = f.Version(); version == "0" {version = f.incVersion()}

	log.Println("[storage.redis]Connect: ", f.isConnected())
	log.Println("[storage.redis]Version: ", version)

	return isConnected
}

func (f *storageRedis) Disconnect(){
	f.storage.Close()
}

func (f *storageRedis) isConnected() bool{
	return f.storage != nil
}


func (f *storageRedis) Pull(index string) (record string){

	defer f.unLock(index)

	log.Println("[storage.redis]Pull: ", index)
	if f.lock(index){

		if record,_ = f.storage.Cmd("HGET", f.storageKey, index).Str(); record == ""{
			log.Println("[storage.redis]Pull: no data for index", index)
		}

		f.storage.Cmd("HDEL", f.storageKey, index)

		f.incVersion()

	}else{
		log.Println("[storage.redis]Pull: lock fail for", index)
	}

	return record
}

func (f *storageRedis) Push(index string, record string, force bool) (result bool, err error){

	log.Println("[storage.redis]Push: ", index, record, force)

	var resp int

	if force {
		resp, err = f.storage.Cmd("HSET", f.storageKey, index, record).Int()
	}else{
		resp, err = f.storage.Cmd("HSETNX", f.storageKey, index, record).Int()
	}

	if result = resp > 0; !result{
		err = errIndexExist

	}else{
		f.incVersion()
	}

	log.Println("[storage.redis]Push: ", "result:", result)
	log.Println("[storage.redis]Push: ", "error:", err)

	return result, err
}


func (f *storageRedis) List() (data map[string]string){

	data, _ = f.storage.Cmd("HGETALL", f.storageKey).Map()

	return data
}

func (f *storageRedis) Flush(){
	f.incVersion()
	f.storage.Cmd("DEL", f.storageKey)
}


func (f *storageRedis) Version() (version string){

	intVersion, _ := f.storage.Cmd("GET", f.storageVer).Int()

	version = strconv.Itoa(intVersion)

	return version
}

func (f *storageRedis) incVersion() (version string){

	oldVersion := f.Version()

	intVersion, _ := f.storage.Cmd("INCR", f.storageVer).Int()

	version = strconv.Itoa(intVersion)

	log.Println("[storage.redis]Version: ", "update:", oldVersion,"-->",intVersion)

	return version
}


func (f *storageRedis) lock(index string) bool{

	l,_ := f.storage.Cmd("HSETNX", f.storageLock, index, 1).Int()

	return l==1
}

func (f *storageRedis) unLock(index string) {

	f.storage.Cmd("HDEL", f.storageLock, index)
}