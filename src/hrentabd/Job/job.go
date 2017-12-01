package Job

import (
	"hrentabd"
	"time"
	"encoding/json"
	"log"
	"math"
)

type job struct {
	Ind     string
	Com     string
	Stm     time.Time
	Rep     int64
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
func (j *job)IsPeriodic() bool{
	return j.Rep >0
}

func (j *job)SetPeriod(period int64){
	j.Rep = period
}

func (j *job)GetPeriod()(period int64){
	return j.Rep
}

func (j *job)NextPeriod(){

	// current time > time start
	if j.IsPeriodic(){
		if diff := time.Now().Unix() - j.TimeStart().Unix(); diff>j.GetPeriod(){
			shiftTime := int64(math.Ceil(float64(diff) / float64(j.GetPeriod()))) * j.GetPeriod()
			j.SetTimeStart( time.Unix(j.TimeStart().Unix() + shiftTime, 0) )
		}else{

			j.SetTimeStart( time.Unix(j.TimeStart().Unix() + j.GetPeriod(), 0) )
		}
	}
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

