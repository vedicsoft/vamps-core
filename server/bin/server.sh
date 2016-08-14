#!/bin/bash
set -e

DATE_COMMAND=$(which date)
TIME_STAMP=`${DATE_COMMAND} '+%Y-%m-%d.%H:%M:%S'`
CURRENT_DIR=`pwd`
SERVER_HOME=`cd ..;pwd`
export CONF_FILE="fgfg"
export SERVER_HOME

function default_(){
  echo "Starting main server....."
  ./server.bin
  echo "Main server started successfully...!!"
  echo $! > server.pid
}

function start_(){
    echo "Starting main server....."
    nohup ./server.bin > ../logs/nohup.log 2>&1&
    echo $! > server.pid
    echo "server started successfully!"
}

function stop_(){
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
            default_
            echo $"Usage: $0 {start|stop|restart|status}"
            exit 1
esac

