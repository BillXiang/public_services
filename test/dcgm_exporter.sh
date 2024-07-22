#!/bin/bash
cd /tmp/public_services/test
tar -zxvf dcgm-exporter-3.3.6-3.4.2.tar.gz
cd dcgm-exporter-3.3.6-3.4.2
make binary
./cmd/dcgm-exporter/dcgm-exporter -f ./etc/default-counters.csv