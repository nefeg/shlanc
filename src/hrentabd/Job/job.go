package Job

import (
	"hrentabd"
	"time"
	"encoding/json"
	"log"
)

type job struct {
	Ind     string
	Com     string
	Stm     time.Time
	Rep     bool
}

func New(index string) hrentabd.Job{
	return hrentabd.Job( &job{Ind:index} )
}



func (j *job)Ttl() int64{
	return j.Stm.Unix() - time.Now().Unix()
}

func (j *job)SetTtl(ttl int64){
	j.Stm = time.Unix(time.Now().Unix() + int64(ttl), 0)
}


func (j *job)Index() string{
	return j.Ind
}

func (j *job)Command() string{
	return j.Com
}

func (j *job)SetCommand(command string) {
	j.Com = command
}

// time start
func (j *job)TimeStart() time.Time{
	return j.Stm
}

func (j *job)SetTimeStart(t time.Time){
	j.Stm = t
}

// Repeatable
func (j *job)IsRepeatable() bool{
	return j.Rep
}

func (j *job)SetRepeatable(repeat bool){
	j.Rep = repeat
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

func (j *job) UnSerialize(data string) (job hrentabd.Job){

	err := json.Unmarshal([]byte(data), j)
	if err != nil{
		log.Fatalln(err)
	}

	log.Println("[Job]Unserialize: ", j, err)

	return hrentabd.Job(j)
}

