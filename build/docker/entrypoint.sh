#!/bin/sh

echo "Configuring Snow Local Environment..."

echo "Gin Starting"
cd /go/src/github.com/qit-team/snow
# 支持热更新
gin -p 80 -a 8000 -b build/bin/go-bin -t . -d .  &

#/etc/init.d/supervisord start

while true
do
    echo "hello world" > /dev/null
    sleep 6s
done


