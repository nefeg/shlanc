package ctrl

import "ctrl/controls"

var ComConf []Command = []Command{

	Command(&controls.ComHalt{controls.New("halt",    `\t`)}),
	Command(&controls.ComQuit{controls.New("exit",    `\q`)} ),
	Command(&controls.ComList{controls.New("list",    `\l`)} ),
	Command(&controls.ComAdd{controls.New("add",    `\a`)} ),
	Command(&controls.ComRemove{controls.New("rm",  `\r`)}),
	//Command(&Com.ComQuit{Com.New("save",    `\s`)}),
}