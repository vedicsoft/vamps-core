#!/bin/bash
set -e

DATE_COMMAND=$(which date)
TIME_STAMP=`${DATE_COMMAND} '+%Y-%m-%d.%H:%M:%S'`
CURRENT_DIR=`pwd`
SERVER_HOME=`cd ..;pwd`
export SERVER_HOME
export JWT_PRIVATE_KEY_PATH=../resources/security/private.key
export JWT_PUBLIC_KEY_PATH=../resources/security/public.key

export CADDYPATH=$SERVER_HOME/configs/.caddy
CADDY_CONF_FILE=$SERVER_HOME/webapps/Caddyfile

# token expiration time in hours
export JWT_EXPIRATION_DELTA=72

function default_(){
  echo "starting redis server...."
  ./redis-server ../configs/redis.conf
  echo "redis started successfully!"

  echo "Starting caddy server....."
  nohup ./caddy --conf="$CADDY_CONF_FILE" > ../logs/nohup.log 2>&1&
  echo "Main server started successfully...!!"

  echo "Starting main server....."
  ./server.bin $SERVER_HOME
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
    nohup ./server.bin $SERVER_HOME > ../logs/nohup.log 2>&1&
    echo $! > server.pid
    echo "server started successfully!"

    echo "Starting caddy server....."
    nohup ./caddy --conf="$CADDY_CONF_FILE" > ../logs/nohup.log 2>&1&
    echo "Main server started successfully...!!"
}

function stop_(){
    if [ -f redis.pid ]; then
        kill -9 `cat redis.pid`
        echo "redis stopped successfully!"
        rm -rf redis.pid
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

