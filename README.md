# simple-ami

```/public``` and ```/application``` directories are included here just because ```go run main.go``` process will be created in same directory, just for purposes of this, usually ```build``` and run in ```/var/www/...```

```/core``` directory is not related to AMI example/implementation, it's here to be used as wrapper for ```net/http``` and ```gorilla/websocket``` instead of using external module or in case of ```gorilla/websocket``` using external module directly and generally decluttering of the AMI code

```/server``` is not part of ```/core``` because it is used on applicative implementation

**all important stuff in** ```/app```
