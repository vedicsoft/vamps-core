##Development environment  setup

* Install go-lang 1.6+ (https://golang.org/)
* Make sure you have set GOPATH  variable.
* Add GOPATH/bin to PATH variable
* Clone the source code (https://github.com/vedicsoft/vamps-{component_name}.git) to {GOPATH}/src/github.com/vedicsoft/

## Configuring the DataBase
* Create a database with name ‘vamps’ - change the `{PROJECT_ROOT}/server/configs/config.yaml` if you want to use a
different database name
* Default configuration file is `{PROJECT_ROOT}/server/configs/config.default.yaml` to change and override the
defaults, simply create a new file called `config.yaml` at the same location and have your preferred configs. Do
not change the default configuration file as it's being tracked by the git version control

## Configuring Redis
* [Redis](http://redis.io/) is an open source (BSD licensed), in-memory data structure store, used as database, cache
and message broker. This app use redis primarily as a JWT token storage.
* Configure the `address` and `password` for redis in `config.yaml` file
* This distribution contains an embedded redis instance(`{PROJECT_ROOT}/server/resources/.test/redis-server`) compiled
for Ubuntu 64 bit 16.04 LTS. This is used for integration tests. You have to replace the `redis-server` binary with
the matching redis server for the OS inorder to build the project with test cases.
* To run in development mode you have to configure an external redis instance or use vamps-lb which has an embedded
version of caddy

## Configuring sqlite3 (needed only for the test cases)
* This project use sqlite3 in-memory database to support integration tests
* If you are running a different OS, you have to replace `{PROJECT_ROOT}/server/resources/.test/sqlite3` binary with
the OS compatible version.

## Configuring Caddy web server
* Caddy is a web server like Apache, nginx, or lighttpd, but with different goals, features, and advantages.
* [Caddy](https://caddyserver.com/) is being used to serve the static files of this server.
* We have embedded caddy with this distribution at `{PROJECT_ROOT}/server/webapps/caddy`
* You must replace the binary with the compatible OS version.

## Building the project
* Make sure build.sh file have execute permission.
    * If not provide it the execute permission `chmod +x build.sh`
* Execute build.sh located at the project root directory.
* Build project artifact can be found under {PROJECT_ROOT}/target folder

## Launching the product
* After a successful build a static binary will be placed under {PROJECT_ROOT}/server/bin
* To run the server execute `{PROJECT_ROOT}/server/bin/server.sh start`
* Point your browser to https://localhost:{caddyPort}/{webappname}/
* Default admin credentials
 `Username : admin@super.com Password: admin`

## IDE support for go-lang
* [IDEs and Text editor plugins](https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins)

## Code review comments
* build.sh script does the following for you
1. `go fmt` - automatically fix the majority of mechanical style issues
2. `goimport` -  a superset of gofmt which additionally adds (and removes) import lines as necessary

The rest of [this document](https://github.com/golang/go/wiki/CodeReviewComments) addresses non-mechanical style points.
This document has been prepared by the golang engineering team and it's a must to follow the comments much as possible.

## Platform wiki and other resources
* Please refer [platform wiki](https://github.com/vedicsoft/wiki)


## API development
* Every api should have a swagger definition - developer has to update {SERVER_HOME}/webapps/apieditor/api.yaml file
with the new api.
* Every api should have integration test covering api scenarios.