module github.com/codenotary/immugw

go 1.13

require (
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/codenotary/immudb v1.0.5
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.1 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mwitkow/go-proto-validators v0.3.2 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/pseudomuto/protoc-gen-doc v1.5.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/takama/daemon v0.12.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/ugorji/go v1.1.4 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210824181836-a4879c3d0e89 // indirect
	google.golang.org/grpc v1.40.0
)

replace github.com/takama/daemon v0.12.0 => github.com/codenotary/daemon v0.0.0-20200507161650-3d4bcb5230f4
