# public_services
Public cmd,cron,service for VMs

## build
go get github.com/cubefs/cubefs/util/config  
go get github.com/go-git/go-git/v5  
go get github.com/jacobsa/daemonize  
go get github.com/ntt360/pmon2/app  
go get github.com/ntt360/pmon2/app/god  
go get github.com/ntt360/pmon2/app/model  
go get github.com/ntt360/pmon2/app/output  
go get github.com/ntt360/pmon2/client/proxy  
go get github.com/robfig/cron/v3  
go get github.com/shirou/gopsutil/process  
go build

## run
### run in foreground
export PMON2_CONF=./pmon2_conf.yml  
./public_service -f -c ./public_service.conf

### run in background
./public_service -c ./public_service.conf

### pmon2
https://github.com/ntt360/pmon2
可以使用pmon2进行服务管理
export PMON2_CONF=/home/public_services/pmon2_conf.yml
```
Usage:
  pmon2 [command]

Available Commands:
  del         del process by id or name
  desc        print the process detail message
  exec        run one binary golang process file
  help        Help about any command
  ls          list all processes
  reload      reload some process
  start       start some process by id or name
  stop        stop running process
  log         display process log by id or name
  logf        display process log dynamic by id or name
  version     show pmon2 version
```
