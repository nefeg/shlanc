package api

import (
	"time"
	"math"
	"encoding/json"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

type Job interface {
	Ttl()               int64
	Index()             string
	Command()           string
	At()                time.Time
	Period()            int64


	TtlX(ttl int64)
	CommandX(command string)
	AtX(t time.Time)
	PeriodX(period int64)

	// repeatable
	IsPeriodic() bool
	NextPeriod()

	Serialize() string
	UnSerialize(data string) Job
}


func NewJob(index string) Job{
	return Job( &commonJob{index:index} )
}


type commonJob struct {
	index   string
	command string
	at      time.Time
	period  int64
}

type commonJobTemplate struct {
	Index   string      `json:"index"`
	Command string      `json:"command"`
	At      time.Time   `json:"at"`
	Period  int64       `json:"period"`
}

func (c *commonJob) Ttl() int64{
	return c.at.Unix() - time.Now().Unix()
}

func (c *commonJob)TtlX(ttl int64){
	c.at = time.Unix(time.Now().Unix() + int64(ttl), 0)
}


func (c *commonJob)Index() string{
	return c.index
}

func (c *commonJob)Command() string{
	return c.command
}

func (c *commonJob)CommandX(command string) {
	c.command = command
}

// time start
func (c *commonJob)At() time.Time{
	return c.at
}

func (c *commonJob)AtX(t time.Time){
	c.at = t
}

// Repeatable
func (c *commonJob)IsPeriodic() bool{
	return c.period >0
}

func (c *commonJob)PeriodX(period int64){
	c.period = period
}

func (c *commonJob)Period()(period int64){
	return c.period
}

func (c *commonJob)NextPeriod(){

	// current time > time start
	if c.IsPeriodic(){
		if diff := time.Now().Unix() - c.at.Unix(); diff>c.Period(){
			shiftTime := int64(math.Ceil(float64(diff) / float64(c.Period()))) * c.Period()
			c.AtX( time.Unix(c.at.Unix() + shiftTime, 0) )
		}else{

			c.AtX( time.Unix(c.at.Unix() + c.Period(), 0) )
		}
	}
}


func (c *commonJob) MarshalJSON() ([]byte, error) {
	return json.Marshal(commonJobTemplate{
		c.index,
		c.command,
		c.at,
		c.period,
	})
}

func (c *commonJob) UnmarshalJSON(b []byte) error {

	temp := &commonJobTemplate{}

	if err := json.Unmarshal(b, temp); err != nil {
		return err
	}

	c.index = temp.Index
	c.CommandX(temp.Command)
	c.AtX(temp.At)
	c.PeriodX(temp.Period)

	return nil
}



// interface storage.Serializable
func (c *commonJob) Serialize() string{

	s, err := c.MarshalJSON()
	if err != nil{
		slog.FatalF("[api.Job] Serialize: %v\n", err)
	}

	slog.DebugF("[api.Job] Serialize: %s [%v]\n", string(s), err)

	return string(s)
}

func (c *commonJob) UnSerialize(data string) (job Job){

	slog.DebugLn("[api.Job] UnSerialize: ", data)

	err := c.UnmarshalJSON([]byte(data))
	if err != nil{
		slog.FatalF("[api.Job] %v", err)
	}

	slog.DebugLn("[api.Job] UnSerialize: ", c, err)

	return Job(c)
}

