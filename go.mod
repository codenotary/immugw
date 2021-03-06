module github.com/codenotary/immugw

go 1.13

require (
	github.com/codenotary/immudb v0.9.3-0.20210422101551-cdf372b3ecf6
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/prometheus/client_golang v1.5.1
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/takama/daemon v0.12.0
	google.golang.org/grpc v1.34.0
)

replace github.com/takama/daemon v0.12.0 => github.com/codenotary/daemon v0.0.0-20200507161650-3d4bcb5230f4
