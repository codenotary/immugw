<!--
---

title: "immugw"

custom_edit_url: https://github.com/codenotary/immudb/edit/master/README.md
---

-->

# immugw [![License](https://img.shields.io/github/license/codenotary/immudb)](LICENSE) <img align="right" src="img/Black%20logo%20-%20no%20background.png" height="47px" />

[![Build Status](https://travis-ci.com/codenotary/immudb.svg?branch=master)](https://travis-ci.com/codenotary/immudb)
[![Go Report Card](https://goreportcard.com/badge/github.com/codenotary/immugw)](https://goreportcard.com/report/github.com/codenotary/immugw)
[![Slack](https://img.shields.io/badge/join%20slack-%23immutability-brightgreen.svg)](https://slack.vchain.us/)
[![Discuss at immudb@googlegroups.com](https://img.shields.io/badge/discuss-immudb%40googlegroups.com-blue.svg)](https://groups.google.com/group/immudb)
[![Immudb Careers](https://img.shields.io/badge/careers-We%20are%20hiring!-blue?style=flat)](https://immudb.io/careers/)


**immugw** is the intelligent REST proxy that connects to immudb and provides a RESTful interface for applications. We recommend to run immudb and immugw on separate machines to enhance security


#### Build the binaries yourself

To build the binaries yourself, simply clone this repo and run

```
make all
```

#### immugw first start

##### Run immugw binary

```bash
# run immugw in the foreground
./immugw
```
##### Run immugw as a service

Service installation and management are supported on Linux, Windows, OSX and FreeBSD operating systems.

```
# install immugw service
./immuadmin service immugw install

# check current immugw service status
./immuadmin service immugw status

# stop immugw service
./immuadmin service immugw stop

# start immugw service
./immuadmin service immugw start
```

The linux service is using the following defaults:

| File or configuration   | location                   |
| ----------------------- | -------------------------- |
| executable              | /usr/sbin/immugw           |
| all configuration files | /etc/immugw                |
| pid file                | /var/lib/immugw/immugw.pid |
| log files               | /var/log/immugw            |

The FreeBSD service is using the following defaults:

| File or configuration   | location            |
| ----------------------- | ------------------- |
| executable              | /usr/sbin/immugw    |
| all configuration files | /etc/immugw         |
| pid file                | /var/run/immugw.pid |
| log files               | /var/log/immugw     |

The Windows service is using the following defaults:

| File or configuration   | location                             |
| ----------------------- | ------------------------------------
| executable              | Program Files\Immugw\immugw.exe      |
| configuration file      | ProgramData\Immugw\config\immugw.toml|
| all data files          | ProgramData\Immugw\                  |
| pid file                | ProgramData\Immugw\config\immugw.pid |
| log file                | ProgramData\Immugw\config\immugw.log |



Simply run `./immugw -d` to start immugw on the same machine as immudb (test or dev environment) or pointing to the remote immudb system ```./immugw --immudb-address "immudb-server"```.

If you want to stop immugw în that case you need to find the process `ps -ax | grep immugw` and then `kill -15 <pid>`. Windows PowerShell would be `Get-Process immugw* | Stop-Process`.

```bash
immu gateway: a smart REST proxy for immudb - the lightweight, high-speed immutable database for systems and applications.
It exposes all gRPC methods with a REST interface while wrapping all SAFE endpoints with a verification service.

Environment variables:
  IMMUGW_ADDRESS=0.0.0.0
  IMMUGW_PORT=3323
  IMMUGW_IMMUDB_ADDRESS=127.0.0.1
  IMMUGW_IMMUDB_PORT=3322
  IMMUGW_DIR=.
  IMMUGW_PIDFILE=
  IMMUGW_LOGFILE=
  IMMUGW_DETACHED=false
  IMMUGW_MTLS=false
  IMMUGW_SERVERNAME=localhost
  IMMUGW_AUDIT=false
  IMMUGW_AUDIT_INTERVAL=5m
  IMMUGW_AUDIT_USERNAME=immugwauditor
  IMMUGW_AUDIT_PASSWORD=
  IMMUGW_AUDIT_SIGNATURE=ignore
  IMMUGW_PKEY=
  IMMUGW_CERTIFICATE=
  IMMUGW_CLIENTCAS=

Usage:
  immugw [flags]
  immugw [command]

Available Commands:
  help        Help about any command
  version     Show the immugw version

Flags:
  -a, --address string            immugw host address (default "0.0.0.0")
      --audit                     enable audit mode (continuously fetches latest root from server, checks consistency against a local root and saves the latest root locally)
      --audit-interval duration   interval at which audit should run (default 5m0s)
      --audit-password string     immudb password used to login during audit; can be plain-text or base64 encoded (must be prefixed with 'enc:' if it is encoded)
      --audit-username string     immudb username used to login during audit (default "immugwauditor")
      --certificate string        server certificate file path (default "./tools/mtls/4_client/certs/localhost.cert.pem")
      --clientcas string          clients certificates list. Aka certificate authority (default "./tools/mtls/2_intermediate/certs/ca-chain.cert.pem")
      --config string             config file (default path are configs or $HOME. Default filename is immugw.toml)
  -d, --detached                  run immudb in background
      --dir string                program files folder (default ".")
  -h, --help                      help for immugw
  -k, --immudb-address string     immudb host address (default "127.0.0.1")
  -j, --immudb-port int           immudb port number (default 3322)
      --logfile string            log path with filename. E.g. /tmp/immugw/immugw.log
  -m, --mtls                      enable mutual tls
      --pidfile string            pid path with filename. E.g. /var/run/immugw.pid
      --pkey string               server private key path (default "./tools/mtls/4_client/private/localhost.key.pem")
  -p, --port int                  immugw port number (default 3323)
      --servername string         used to verify the hostname on the returned certificates (default "localhost")

Use "immugw [command] --help" for more information about a command.

```

### Docker

**immugw**  is also available as docker images on dockerhub.com.

| Component  | Container image                                |
| ---------- | ---------------------------------------------- |
| immugw     | https://hub.docker.com/r/codenotary/immugw     |

#### Run immugw

```
docker run -it -d -p 3323:3323 --name immugw --env IMMUGW_IMMUDB_ADDRESS=immudb codenotary/immugw:latest
```


#### Build the container images yourself

If you want to build the container images yourself, simply clone this repo and run

```
docker build -t myown/immugw:latest -f Dockerfile .
```

## Why immugw

**immugw** provides a simple solution to interact with immudb with REST protocol, without taking in charge the merkle tree root hash file management and concurrency related complexity.

#### immugw communication

**immugw**  proxies REST client communication and gRPC server interface. For security purposes immugw should not run on the same server as immudb. The following diagram shows how the communication works:

![immugw communication explained](img/immugw-diagram.png)

### CURL examples

#### Login
```shell script
curl --location --request POST '127.0.0.1:3323/login' \
--header 'Authorization;' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user": "aW1tdWRi",
    "password": "aW1tdWRi"
}'
```
#### Use Database
```shell script
curl --location --request GET '127.0.0.1:3323/db/use/defaultdb' \
--header 'Content-Type: application/json' \
--header 'Authorization: {{token}}'
```
#### Login
```shell script
curl --location --request POST '127.0.0.1:3323/login' \
--header 'Authorization;' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user": "aW1tdWRi",
    "password": "aW1tdWRi"
}'
```
#### Verified Set
```shell script
curl --location --request POST '127.0.0.1:3323/db/verified/set' \
--header 'Content-Type: application/json' \
--header 'Authorization: {{token}}' \
--data-raw '{
  "setRequest": {
    "KVs": [
      {
        "key": "a2V5MQ==",
	   "value": "dmFsMQ=="
      }
    ]
  }
}'
```
#### Verified Get
```shell script
curl --location --request POST '127.0.0.1:3323/db/verified/get' \
--header 'Content-Type: application/json' \
--header 'Authorization: {{token}}' \
--data-raw '{
  "keyRequest": {
    "key": "a2V5MQ=="
  }
}'
```
#### Verified Reference
```shell script
curl --location --request POST '127.0.0.1:3323/db/verified/setreference' \
--header 'Content-Type: application/json' \
--header 'Authorization: {{token}}' \
--data-raw '{
  "referenceRequest": {
    "key": "dGFnMQ==",
    "referencedKey": "a2V5MQ==",
    "atTx": "0"
  }
}'
```
#### Verified ZAdd
```shell script
curl --location --request POST '127.0.0.1:3323/db/verified/zadd' \
--header 'Content-Type: application/json' \
--header 'Authorization: {{token}}' \
--data-raw '{
  "zAddRequest": {
    "set": "c2V0MQ==",
    "score": 15.5,
    "key": "a2V5MQ==",
    "atTx": "0"
  }
}'
```
#### ZScan
```shell script
curl --location --request POST '127.0.0.1:3323/db/zscan' \
--header 'Content-Type: application/json' \
--header 'Authorization: {{token}}' \
--data-raw '{
  "set": "c2V0MQ=="
}'
```
#### History
```shell script
curl --location --request POST '127.0.0.1:3323/db/history' \
--header 'Authorization: {{token}}' \
--header 'Content-Type: application/json' \
--data-raw '{
  "key": "a2V5NQ=="
}'
```
#### Verified Transaction
```shell script
curl --location --request GET '127.0.0.1:3323/db/verified/tx/1' \
--header 'Content-Type: application/json' \
--header 'Authorization: {{token}}'
```
#### SQL Exec
```shell script
curl --location --request POST '127.0.0.1:3323/db/sqlexec' \
--header 'Authorization: {{token}}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "sql":"CREATE TABLE mytable23 (id INTEGER, amount INTEGER, total INTEGER, title VARCHAR, content BLOB, isPresent BOOLEAN, PRIMARY KEY id)"
}'
```

#### SQL Exec insert
```shell script
curl --location --request POST '127.0.0.1:3323/db/sqlexec' \
--header 'Authorization: v2.public.eyJkYXRhYmFzZSI6IjUiLCJleHAiOiIyMDIxLTEwLTI4VDE4OjU1OjAyKzAyOjAwIiwic3ViIjoiaW1tdWRiIn3-aNUXqydajYFR9Aa7-q40JepLuA0tsPXeR1nRo75jA1H45RZZU9Twt6EVi-4bS4gpzeQcRNEdJs8U5oM5urcM.aW1tdWRi' \
--header 'Content-Type: application/json' \
--data-raw '{
    "sql":"INSERT INTO myTable23 (id, amount, title, content, isPresent) VALUES (2, 1000, '\''title 1'\'', x'\''626C6F6220636F6E74656E74'\'', true)"
}'
```
> byte arrays need to be hex encoded
#### SQL Query
```shell script
curl --location --request POST '127.0.0.1:3323/db/sqlquery' \
--header 'Authorization: {{token}}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "sql":"SELECT * from mytable23;"
}'
```

#### SQL Verifiable sql get
Its possible also to tamperproof verify a SQL row.
```shell script
curl --location --request POST '127.0.0.1:3323/db/verified/row' \
--header 'Authorization: v2.public.eyJkYXRhYmFzZSI6IjUiLCJleHAiOiIyMDIxLTEwLTI4VDE4OjU1OjAyKzAyOjAwIiwic3ViIjoiaW1tdWRiIn3-aNUXqydajYFR9Aa7-q40JepLuA0tsPXeR1nRo75jA1H45RZZU9Twt6EVi-4bS4gpzeQcRNEdJs8U5oM5urcM.aW1tdWRi' \
--header 'Content-Type: application/json' \
--data-raw '{
    "row": {
            "columns": [
                "(testdb1.mytable23.id)",
                "(testdb1.mytable23.amount)",
                "(testdb1.mytable23.total)",
                "(testdb1.mytable23.title)",
                "(testdb1.mytable23.content)",
                "(testdb1.mytable23.ispresent)"
            ]
            },
            "values": [
                {
                    "n": "2"
                },
                {
                    "n": "1000"
                },
                {
                    "null": null
                },
                {
                     "s": "title 1"
                },
                {
                    "bs": "YmxvYiBjb250ZW50"
                },
                {
                    "b": true
                }
            ],
        "table": "mytable23",
        "pkValues": [
            {
                "n": "2"
            }
        ]
      }'
```
> byte arrays need to be b64 encoded
#### Logout
```shell script
curl --location --request POST '127.0.0.1:3323/logout' \
--header 'Authorization: {{token}}' \
--header 'Content-Type: application/json'
```


## License

immugw is [Apache v2.0 License](LICENSE).

immudb re-distributes other open-source tools and libraries - [Acknowledgements](https://github.com/codenotary/immudb/blob/master/ACKNOWLEDGEMENTS.md).
