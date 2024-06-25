#!/bin/bash

DECLARATION_PATH="/Users/k.song/src/golang/as3-benchmark/one-tenant"

for i in 1 10 20 30 40 50 60 70 80 90 100 110 120 130 140 150 160 170 180 190 200
do
  echo "ADD $i"
  go run cmd/as3-benchmark/main.go --ops=add --declaration=$DECLARATION_PATH/declaration-$i.json --bigip-host=10.155.223.51 --bigip-username=admin --bigip-password=admin
  sleep 1
  echo "DEL $i"
  go run cmd/as3-benchmark/main.go --ops=del --declaration=$DECLARATION_PATH/declaration-$i.json --bigip-host=10.155.223.51 --bigip-username=admin --bigip-password=admin
  sleep 1
done
