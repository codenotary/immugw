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
##### Run immugw as a service (using immuadmin)

Please make sure to build or download the immugw and immuadmin component and save them in the same work directory when installing the service.

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
| all configuration files | /etc/immudb                |
| pid file                | /var/lib/immudb/immugw.pid |
| log files               | /var/log/immudb            |

The FreeBSD service is using the following defaults:

| File or configuration   | location            |
| ----------------------- | ------------------- |
| all configuration files | /etc/immudb         |
| pid file                | /var/run/immugw.pid |
| log files               | /var/log/immudb     |



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
  IMMUGW_PKEY=./tools/mtls/4_client/private/localhost.key.pem
  IMMUGW_CERTIFICATE=./tools/mtls/4_client/certs/localhost.cert.pem
  IMMUGW_CLIENTCAS=./tools/mtls/2_intermediate/certs/ca-chain.cert.pem
  IMMUGW_AUDIT="false"
  IMMUGW_AUDIT_PASSWORD=""
  IMMUGW_AUDIT_USERNAME=""

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


### immugw RESTful API reference

You can find the swagger schema here:

https://github.com/codenotary/immudb/blob/master/pkg/api/schema/gw.schema.swagger.json

If you want to run the Swagger UI, simply run the following docker command after you cloned this repo:

```
docker run -d -it -p 8081:8080 --name swagger-immugw -v ${PWD}/pkg/api/gw.schema.swagger.json/openapi.json -e SWAGGER_JSON=/openapi.json  swaggerapi/swagger-ui
```

## License

immugw is [Apache v2.0 License](LICENSE).

immudb re-distributes other open-source tools and libraries - [Acknowledgements](https://github.com/codenotary/immudb/blob/master/ACKNOWLEDGEMENTS.md).
