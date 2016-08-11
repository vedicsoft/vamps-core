#!/bin/bash
set -e

PROJECT_NAME='vamps-core'
PROJECT_ROOT=`pwd`
export GOBIN="$PROJECT_ROOT"

# setting SERVER_HOME for test cases
HOME=`cd server;pwd`
echo $HOME
export SERVER_HOME=$HOME

echo 'Exporting GO variables.'

if [ -z "$GOPATH" ]; then
 echo "Build failed due to GOPATH has not been set."
 exit 1
fi

command -v godep >/dev/null 2>&1 || { echo >&2 "Godep required. Installing godep.";  go get github.com/tools/godep;}

rm -rf $PROJECT_ROOT/target
mkdir -p $PROJECT_ROOT/target

echo 'Gom install dependencies. This might take some time...'
godep go install ./...

echo "Executing test"
godep go test -v ./...

mv $PROJECT_NAME $PROJECT_ROOT/server/bin/server.bin

echo 'GO build complete.'

cd $PROJECT_ROOT/target

if [ "$1" = "--release" ];then
 echo "Writing version information to versioninfo.md"
 DATE_COMMAND=$(which date)
 TIME_STAMP=`${DATE_COMMAND} '+%Y-%m-%d.%H:%M:%S'`

 echo "Time Stamp : ${TIME_STAMP}" > ../server/versioninfo.md
 LAST_COMMIT_ID=$(git log | head -1 | sed s/'commit '//)
 echo "Last Commit ID : ${LAST_COMMIT_ID}" >> ../server/versioninfo.md
 GIT_BRANCH=$(git branch)
 echo "Branch : ${GIT_BRANCH}" >> ../server/versioninfo.md
fi

echo "Start creating new distribution"
mkdir $PROJECT_NAME
cp -r ../server/* $PROJECT_NAME/
zip -rq $PROJECT_NAME.zip ./$PROJECT_NAME/* -x *.log -x *.out
rm -rf $PROJECT_NAME
echo "Distribution creation complete."