Distributed Concurrency job manager
=============================

Features
--------
 - run job by ttl
 - run job at time
 - repeat jobs by period
 - remote/local job storage (redis/file)
 - distributing job list for concurrency execution
 - decentralized structure
 
 
INSTALL
-------

#### From Ubuntu repository via apt-get

    # add PPA to apt source list:
    sudo add-apt-repository ppa:onm/shlanc
    sudo apt-get update
    
    #install package
    sudo apt-get install shlanc



#### From .deb package

Download package from [Launchpad PPA](https://launchpad.net/~onm/+archive/ubuntu/shlanc/+packages) and install:

        sudo dpkg -i path/to/package.deb

#### Copy binaries
Just copy binary file and create config file

- fix rights `chmod +x shlancd`
- run `./shlancd -c config.json` 


Configuring
-----------


#### Default config:
```
{
	"runMissed": true,

	"storage": {
		"type": "redis",
		"options": {
			"network":  "tcp",
			"address":  "127.0.0.1:6379",
			"key":      "hren"
		}
	},

	"client": {
		"type": "socket",
		"options":{
			"network":  "tcp",
			"address":  "127.0.0.1:6607"
		}
	},

	"executor":{
		"type": "local",
		"options":{
			"silent":   true,
			"async":    true
		}
	}
}
```


#### Parameters

- `runMissed` - run/skip missed jobs

**[STORAGE]**


- **`type:redis`** - use redis as job storage

- `options` - connection options for redis
    
	`options.network` - socket type: `tcp|udp|unix`
    
    `options.address` - socket address(TCP/IP: `127.0.0.1:6379`, unix: `"/path/to/redis.sock"`)
    
    `options.key` - prefix for keys in database

example:
 ```	
"storage": {
    "type": "redis",
    "options": {
        "network":  "tcp",
        "address":  "127.0.0.1:6379",
        "key":      "hren"
    }
}
```

- **`type:file`** - use local file as job storage(**not support job distribution and concurrency**)

    `options.path` - path to db-file


```
  "storage": {
    "type": "file",
    "options": {
      "path":   "/tmp/hren.db"
    }
  }
```
 

**[CLIENT]**

- **`type:socket`** - connect via socket
- `options` - connection options for socket
    
	`options.network` - socket type: `tcp|udp|unix`
    
    `options.address` - socket address(TCP/IP: `127.0.0.1:6379`, unix: `"/path/to/redis.sock"`)

```
"client": {
    "type": "socket",
    "options":{
        "network":  "tcp",
        "address":  "127.0.0.1:6607"
    }
}
```
**Just use telnet as client!**


**[EXECUTOR]**

Just use default config:
```
	"executor":{
		"type": "local",
		"options":{
			"silent":   true,
			"async":    true
		}
	}
```


Usage
--------------

- using shlanc cli:

```
username:~/$ shlanc -h
NAME:
   ShLANC-client - [SH]lanc [L]ike [A]s [N]ot [C]ron

USAGE:
   shlanc [global options] command [command options] [arguments...]

COMMANDS:
     add, a         Add job
     list, l        Show list of jobs
     remove, rm, r  Remove jobs by ID or time of start
     purge          Remove all jobs
     get, g         Get job by id
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  path to daemon config-file
   --debug                   show debug log
   --help, -h                show help
   --version, -v             print the version
```

- using telnet (need real connection):
```
username:~/$ $ telnet 127.0.0.1 6607
 Trying 127.0.0.1...
 Connected to 127.0.0.1.
 Escape character is '^]'.
 ShlaNc terminal connected OK
 Type "help" for show available commands
 >>help
 NAME:
    ShLANC-client - [SH]lanc [L]ike [A]s [N]ot [C]ron

 USAGE:
    shlancd [global options] command [command options] [arguments...]

 COMMANDS:
      add, a         Add job
      list, l        Show list of jobs
      remove, rm, r  Remove jobs by ID or time of start
      purge          Remove all jobs
      get, g         Get job by id
      exit, q        close connection
      help, h        Shows a list of commands or help for one command
 >>
```


License
-------

License: pick the one which suits you best:

- GPL v3 see <https://www.gnu.org/licenses/gpl.html>
- APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
