#!/bin/bash
set -e

DATE_COMMAND=$(which date)
TIME_STAMP=`${DATE_COMMAND} '+%Y-%m-%d.%H:%M:%S'`
CURRENT_DIR=`pwd`
SERVER_HOME=`cd ..;pwd`
export SERVER_HOME

function default_(){
  echo "starting redis server...."
  ./redis-server ../configs/redis.conf
  echo "redis started successfully!"

  echo "Starting main server....."
  ./server.bin ${SERVER_HOME}
  echo "Main server started successfully...!!"
  echo $! > server.pid
}

function start_(){
    echo "starting redis server..!"
    if test -f "../configs/redis.conf"; then
        ./redis-server ../configs/redis.conf
    else
       ./redis-server ../configs/redis.default.conf
    fi
    echo "redis started successfully!"

    nohup ./server.bin ${SERVER_HOME} > ../logs/nohup.log 2>&1&
    echo $! > server.pid
    echo "server started successfully!"
}

function stop_(){
    if [ -f redis.pid ]; then
        kill -9 `cat redis.pid`
        echo "redis stopped successfully!"
        rm -rf redis.pid
    fi

    if [ -f caddy.pid ]; then
        kill -9 `cat caddy.pid`
        echo "caddy stopped successfully!"
        rm -rf caddy.pid
    fi

    if [ -f server.pid ]; then
        kill -9 `cat server.pid`
        echo "server stopped successfully!"
        rm -rf server.pid
    fi
}

case "$1" in
        "")
           default_
           ;;

        start)
            start_
            ;;

        stop)
            stop_
            ;;

        status)
            process=$(ps -ef | grep server.bin | grep -v grep)
            if [ "$process" ]; then
             echo "server is up and running."
            else
             echo "server is not running at the moment."
            fi
            ;;
        restart)
            stop_
            start_
            ;;
        *)
            echo $"Usage: $0 {start|stop|restart|status}"
            exit 1
esac

