package cli

import (
	"errors"
	"fmt"
	sapi "shlancd/app/api"


)

var ErrCmdArgs  = errors.New("ERR: expected argument")




func view(job sapi.Job) string{

	if job == nil{
		return ""
	}

	var period string
	if job.IsPeriodic(){
		period = fmt.Sprint(job.Period())
	}else{
		period = "null"
	}

	return fmt.Sprintln(
		job.At().String(),"\t",
		job.Index(),"\t",
		period,"\t",
		"\""+job.Command()+"\"", "\t",
		job.Ttl(),"\t",
	)
}


