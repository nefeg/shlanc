package storage

import (
	"os"
	"log"
	"bufio"
	"io"
	"errors"
	"storage/file"
)


type storageFile struct{

	storagePath string
	storage     *os.File
	items       map[string]file.Item
}


var errIndexExist = errors.New("index already exist")

func NewStorageFile(path string) *storageFile{

	storage := &storageFile{storagePath:path, items:map[string]file.Item{}}
	storage.Connect()
	storage.loadItems()

	return storage
}


func (f *storageFile) Connect() (isConnected bool){

	if !f.isConnected(){

		if fh, err := os.OpenFile(f.storagePath, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644); err == nil{
			f.storage = fh
		}else{
			log.Panicln(err)
		}
	}

	isConnected = f.isConnected()

	log.Println("[storage.file]Connect: ", f.isConnected())

	return isConnected
}

func (f *storageFile) Disconnect(){
	f.commit()
	f.storage.Close()
}

func (f *storageFile) isConnected() bool{
	return int(f.storage.Fd()) > 0
}

func (f *storageFile) Get(index string) (record string){
	if item, isset := f.items[index]; isset{
		record = item.Data()
	}

	return record
}

func (f *storageFile) Add(index string, record string, force bool) (result bool, err error){

	fi := file.NewItem(index, record)

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

func (f *storageFile) Rm(index string) (result bool, err error){

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

func (f *storageFile) List() (data map[string]string){

	data = map[string]string{}
	for i,d := range f.items{
		data[i] = d.Data()
	}

	return data
}

func (f *storageFile) Flush() {

	for _,i := range f.items{
		f.rmItem(i)
	}

	f.commit()
}



func (f *storageFile) hasIndex(index string) bool{

	_, isset := f.items[index]
	return isset
}

func (f *storageFile) commit() (size int, err error){

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

func (f *storageFile) flush(){}


func (f *storageFile) addItem(fi file.Item){
	f.items[fi.Index()] = fi
}

func (f *storageFile) rmItem(fi file.Item){
	delete(f.items, fi.Index())
}

// load items from file
func (f *storageFile) loadItems(){

	var c = 0
	rd := bufio.NewReader(io.Reader(f.storage))
	for{
		if l,e := rd.ReadString('\n'); e == nil{

			fi := file.NewItemFromString(l)


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
