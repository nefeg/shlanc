package controls

import (
	"errors"
	"log"
)


var ArgsBindingError error = errors.New("invalid number of arguments")
var ArgsGetError error = errors.New("unknown argument %s")

type Arguments struct{
	container map[string]string
}

func (a *Arguments) this() *Arguments{
	if a.container == nil{
		a.container = map[string]string{}
	}

	return a
}

func (a *Arguments) Bind(values, names []string) (*Arguments, error){

	var err error

	if len(values) != len(names){
		err = ArgsBindingError

	}else{
		for index, value := range values{
			a.SetX(names[index], value)
		}
	}

	return a,err
}

func (a *Arguments) GetX(name string) (string){

	value, isset := a.this().container[name]

	if !isset{
		log.Fatalf(ArgsGetError.Error(), name)
	}

	return value
}

func (a *Arguments) SetX(name, value string) *Arguments{

	a.this().container[name] = value

	return a
}

func (a *Arguments) Num() int{
	return len(a.this().container)
}
