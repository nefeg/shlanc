package file

import (
	"bytes"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

type Item interface{

	ToString() (itemString string)

	Index() (index string)
	Data() (data string)
}


type item struct{

	index   string
	data    string
}


func NewItem(index string, data string) Item{

	return Item( &item{index,data} )
}

func NewItemFromString(itemString string) Item{

	fi := &item{}

	if parts := bytes.Split([]byte(itemString), []byte{0}); len(parts) != 2{
		slog.Infoln("[storage.file.item] FromString: Corrupted data")

	}else{
		fi.index    = string(parts[0])
		fi.data     = string(parts[1])
	}

	return Item(fi)
}


func (fi *item) ToString() (itemString string){
	return string(fi.index) + string(byte(0)) + string(fi.data) + "\n"
}

func (fi *item) Index() (index string){
	return fi.index
}

func (fi *item) Data() (data string){
	return fi.data
}
