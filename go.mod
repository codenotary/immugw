module github.com/codenotary/immugw

go 1.13

require (
	github.com/codenotary/immudb v1.4.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/prometheus/client_golang v1.12.2
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.12.0
	github.com/stretchr/testify v1.7.1
	github.com/takama/daemon v0.12.0
	google.golang.org/grpc v1.46.2
)

replace github.com/takama/daemon v0.12.0 => github.com/codenotary/daemon v0.0.0-20200507161650-3d4bcb5230f4

replace github.com/spf13/afero => github.com/spf13/afero v1.5.1
