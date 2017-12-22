package cli

import "cli/Com"

var commandConfig = []Cmd{

		Cmd( &Com.Quit{Com.New("exit",    `\q`, "  exit (\\q) - meta-command, send exit error")} ),
		Cmd( &Com.List{Com.New("list",    `\l`, "  list (\\l) - show list of jobs")} ),
		Cmd( &Com.Get{Com.New("get",    `\g`, "  get (\\g) - get job")} ),
		Cmd( &Com.Add{Com.New("add",    `\a`, "  Add job into job list")} ),
		Cmd( &Com.Remove{Com.New("rm",  `\r`, "  rm (\\r) - remove job from job list")} ),
	}
