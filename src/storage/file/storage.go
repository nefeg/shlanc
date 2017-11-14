package file

import (
	"os"
	"log"
	"bufio"
	"io"
	"errors"
	. "storage"
)

type storage struct{

	storagePath string
	storage     *os.File
	items       map[Index]*item
}


var errIndexExist = errors.New("index already exist")

func NewFileStorage(path string) *storage{

	storage := &storage{storagePath:path, items:map[Index]*item{}}
	storage.Connect()
	storage.loadItems()

	return storage
}


func (f *storage) Connect() (isConnected bool){

	if !f.isConnected(){

		if file, err := os.OpenFile(f.storagePath, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644); err == nil{
			f.storage = file
		}else{
			log.Panicln(err)
		}
	}

	isConnected = f.isConnected()

	log.Println("[storage.file]Connect: ", f.isConnected())

	return isConnected
}

func (f *storage) Disconnect(){
	f.commit()
	f.storage.Close()
}

func (f *storage) isConnected() bool{
	return int(f.storage.Fd()) > 0
}

func (f *storage) Get(index Index) (record Record){
	if item, isset := f.items[index]; isset{
		record = Record(item.Data())
	}

	return record
}

func (f *storage) Add(index Index, record Record, force bool) (result bool, err error){

	fi := NewFileItem(index, record)

	if !f.hasIndex( index ) || force{

		f.addItem(fi)

		if _, err := f.commit(); err != nil{// rollback
			f.rmItem(fi)
		}else{
			result = true
		}

	}else{
		err = errIndexExist
	}

	if err == nil{
		f.addItem(fi)
	}

	return result, err
}

func (f *storage) Rm(index Index) (result bool, err error){

	if item, isset := f.items[index]; isset{

		f.rmItem(item)

		if _,err = f.commit(); err != nil{ // rollback
			f.addItem(item)
		}else{
			result = true
		}

	}

	return result, err
}

func (f *storage) List() (data map[Index]Record){

	data = map[Index]Record{}
	for i,d := range f.items{
		data[i] = d.Data()
	}

	return data
}

func (f *storage) Flush() {

	for _,i := range f.items{
		f.rmItem(i)
	}

	f.commit()
}



func (f *storage) hasIndex(index Index) bool{

	_, isset := f.items[index]
	return isset
}

func (f *storage) commit() (size int, err error){

	var dump string
	for _,item := range f.items{
		dump += item.ToString()
	}


	f.storage.Truncate(0)
	f.storage.Seek(0, 0)

	if dump != "" {
		size, err = f.storage.WriteString( dump )
	}

	// log.Println("[storage.file]commit: ", size, err, dump, ofs)

	return size, err
}

func (f *storage) flush(){}


func (f *storage) addItem(fi *item){
	f.items[fi.Index()] = fi
}

func (f *storage) rmItem(fi *item){
	delete(f.items, fi.Index())
}

// load items from file
func (f *storage) loadItems(){

	var c = 0
	rd := bufio.NewReader(io.Reader(f.storage))
	for{
		if l,e := rd.ReadString('\n'); e == nil{
			fi := &item{}
			fi.FromString(l)

			if f.hasIndex(fi.Index()){
				log.Panicln("[storage.file]loadItems: Duplicated index - ", fi.Index())
			}

			f.addItem( fi )
			c++

		}else if e == io.EOF{
			break

		}else{// unexpected error
			log.Panicln(e)
		}
	}

	log.Println("[storage.file]loadItems: Records loaded - ", c)
}
