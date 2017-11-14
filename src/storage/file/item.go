package file

import (
	"bytes"
	"log"
	. "storage"
)

type item struct{

	index   Index
	data    Record
}

func NewFileItem(index Index, data Record) *item{

	return &item{index,data}
}


func (fi *item) FromString(itemString string) {

	if parts := bytes.Split([]byte(itemString), []byte{0}); len(parts) != 2{
		log.Panicln("[storage.file.item]FromString: Corrupted data")

	}else{
		fi.index    = Index(parts[0])
		fi.data     = Record(parts[1])
	}
}

func (fi *item) ToString() (itemString string){
	return string(fi.index) + string(byte(0)) + string(fi.data) + "\n"
}

func (fi *item) Index() (index Index){
	return fi.index
}

func (fi *item) Data() (data Record){
	return fi.data
}
