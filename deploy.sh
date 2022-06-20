#!/bin/sh

git checkout .
git pull
pkill -f stepn
go build
nohup /root/stepn/stepn >> stepn.log &
tail -f /root/stepn/stepn.log