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

Just copy binary file and create config file

- fix rights `chmod +x hrentabd`
- run `./hrentabd config.json` 


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

```
$ telnet 127.0.0.1 6607
Trying 127.0.0.1...
Connected to 127.0.0.1.
Escape character is '^]'.
HrenTab terminal connected OK
>>\h
Unknown command
  exit (\q) - meta-command, send exit error
	usage:
	  quit (\q)

  list (\l) - show list of jobs
	usage:
	  list (\l) -index <index>
	  list (\l) -ts <timestamp>
	  list (\l) --help

  add (\a) - add job into job list
	usage:
	  add (\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -ttl <ttl>
	  add (\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -ts <timestamp>
	  add (\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -tm <2006-01-02T15:04:05Z07:00>
	  add (\a) --help

  rm (\r) - remove job from job list
	usage:
	  rm (\r) -index <index>
	  rm (\r) -ts <timestamp>
	  rm (\r) --all
	  rm (\r) --help


>>
```



License
-------

License: pick the one which suits you best:

- GPL v3 see <https://www.gnu.org/licenses/gpl.html>
- APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
