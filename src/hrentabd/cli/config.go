package cli

import "hrentabd/cli/Com"

var commandConfig = []Cmd{

		//Cmd( &Halt{New("halt",    `\t`)} ),
		Cmd( &Com.Quit{Com.New("exit",    `\q`, "  exit (\\q) - meta-command, send exit error")} ),
		Cmd( &Com.List{Com.New("list",    `\l`, "  list (\\l) - show list of jobs")} ),
		Cmd( &Com.Add{Com.New("add",    `\a`, "  add (\\a) - add job into job list")} ),
		Cmd( &Com.Remove{Com.New("rm",  `\r`, "  rm (\\r) - remove job from job list")} ),
	}
