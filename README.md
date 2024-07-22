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
