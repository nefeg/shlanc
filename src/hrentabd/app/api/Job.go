package api

import (
	"time"
	"math"
	"encoding/json"
	"log"
)

type Job interface {
	Ttl()               int64
	Index()             string
	Command()           string
	TimeStart()         time.Time

	SetTtl(ttl int64)
	SetCommand(command string)
	SetTimeStart(t time.Time)

	// repeatable
	IsPeriodic() bool
	SetPeriod(period int64)
	GetPeriod()(period int64)
	NextPeriod()


	Serialize() string
	UnSerialize(data string) Job
}


func NewJob(index string) Job{
	return Job( &commonJob{index:index} )
}


type commonJob struct {
	index       string
	command     string
	startTime   time.Time
	period      int64
}

func (c *commonJob)Ttl() int64{
	return c.startTime.Unix() - time.Now().Unix()
}

func (c *commonJob)SetTtl(ttl int64){
	c.startTime = time.Unix(time.Now().Unix() + int64(ttl), 0)
}


func (c *commonJob)Index() string{
	return c.index
}

func (c *commonJob)Command() string{
	return c.command
}

func (c *commonJob)SetCommand(command string) {
	c.command = command
}

// time start
func (c *commonJob)TimeStart() time.Time{
	return c.startTime
}

func (c *commonJob)SetTimeStart(t time.Time){
	c.startTime = t
}

// Repeatable
func (c *commonJob)IsPeriodic() bool{
	return c.period >0
}

func (c *commonJob)SetPeriod(period int64){
	c.period = period
}

func (c *commonJob)GetPeriod()(period int64){
	return c.period
}

func (c *commonJob)NextPeriod(){

	// current time > time start
	if c.IsPeriodic(){
		if diff := time.Now().Unix() - c.TimeStart().Unix(); diff>c.GetPeriod(){
			shiftTime := int64(math.Ceil(float64(diff) / float64(c.GetPeriod()))) * c.GetPeriod()
			c.SetTimeStart( time.Unix(c.TimeStart().Unix() + shiftTime, 0) )
		}else{

			c.SetTimeStart( time.Unix(c.TimeStart().Unix() + c.GetPeriod(), 0) )
		}
	}
}


// interface storage.Serializable
func (c *commonJob) Serialize() string{

	s, err := json.Marshal(c)
	if err != nil{
		log.Fatalln(err)
	}

	log.Println("[Job]Serialize: ", string(s), err)

	return string(s)
}

func (c *commonJob) UnSerialize(data string) (job Job){

	err := json.Unmarshal([]byte(data), c)
	if err != nil{
		log.Fatalln(err)
	}

	log.Println("[Job]Unserialize: ", c, err)

	return Job(c)
}

