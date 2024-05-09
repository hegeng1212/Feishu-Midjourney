#!/bin/bash

cd midjourney

# 判断是否存在包含关键字的进程
if pgrep -f "midjourney-run" >/dev/null; then
    echo "存在包含关键字的进程，开始杀死..."
    pkill -f "midjourney-run"
    echo "进程已被杀死."
else
    echo "不存在包含关键字的进程."
fi

go build -o midjourney-run

pwd

nohup ./midjourney-run > /tmp/midjourney-run.log 2>&1 &
