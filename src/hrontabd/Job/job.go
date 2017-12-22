package Job

import (
	"hrontabd"
	"time"
	"encoding/json"
	"log"
	"github.com/gorhill/cronexpr"
)

type job struct {
	ID      string
	Cln     string // cron-line
	Cmd     string // command
	Com     string // comment
}

func New() *job{
	return &job{}
}

func (j *job)Id() string{
	return j.ID
}

func (j *job)Command() string{
	return j.Cmd
}

func (j *job)CronLine() string{
	return j.Cln
}

func (j *job)Comment() string{
	return j.Com
}

func (j *job)TimeStart() time.Time{

	return cronexpr.MustParse(j.Cln).Next(time.Now())
}


func (j *job)SetID(ID string){
	if j.Id() != ""{
		log.Panicf("[job] SetID(panic): try to change id %s --> %s \n", j.Id(), ID)
	}

	j.ID = ID
}

func (j *job)SetCommand(command string) {
	j.Cmd = command
}

func (j *job)SetCronLine(timeLine string){

	if _, e := cronexpr.Parse(timeLine); e != nil{
		log.Panicln("Invalid format for '-cron'")
	}

	j.Cln = timeLine
}

func (j *job)SetComment(comment string){
	j.Com = comment
}



// interface storage.Serializable
func (j *job) Serialize() string{

	s, err := json.Marshal(j)
	if err != nil{
		log.Fatalln(err)
	}

	log.Println("[Job]Serialize: ", string(s), err)

	return string(s)
}

func (j *job) UnSerialize(data string) (job hrontabd.Job){

	err := json.Unmarshal([]byte(data), j)
	if err != nil{
		log.Fatalln(err)
	}

	log.Println("[Job]Unserialize: ", j, err)

	return hrontabd.Job(j)
}

