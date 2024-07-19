# public_services
Public cmd,cron,service for VMs

## run in foreground
export PMON2_CONF=./pmon2_conf.yml
./public_service -f -c ./public_service.conf

## run in background
./public_service -c ./public_service.conf