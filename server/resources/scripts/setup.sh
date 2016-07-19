#!/bin/bash
set -e

source setup_configs.sh

function default_(){
    echo "Installing VAMPS databse..."
    mysql -u $DASHBOARD_DB_USERNAME -p$DASHBOARD_DB_PASSWORD -h $DASHBOARD_DB_HOST < ../sql/dashboard-db.sql
    echo "Dashboard DB installed successfully."
}

function clean_(){
    echo "Cleaning dashboard databses..."
    mysql -u $DASHBOARD_DB_USERNAME -p$DASHBOARD_DB_PASSWORD -h $DASHBOARD_DB_HOST < ../sql/cleanall.sql
    echo "Cleaning dashboard database complete"
}


function summarize_(){
    echo "Sumarizing databses..."
    mysql -u $SUMMARY_DB_USERNAME -p$SUMMARY_DB_PASSWORD -h $SUMMARY_DB_HOST < ../sql/summarize.sql
    echo "Summarizing complete"
}

case "$1" in
        "")
           default_
           ;;

        clean)
            clean_
            ;;

        *)
            echo $"Usage: $0 {start|stop|restart|status}"
            exit 1
esac

