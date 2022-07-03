module github.com/kujilabo/cocotola-translator-api

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v63.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.25
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.4.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/mattn/go-sqlite3 v1.14.12
	github.com/onrik/gorm-logrus v0.3.0
	github.com/onrik/logrus v0.9.0
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.4.1
	github.com/swaggo/swag v1.8.1
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.31.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.32.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.6.3
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.6.3
	go.opentelemetry.io/otel/sdk v1.6.3
	go.opentelemetry.io/otel/trace v1.7.0
	golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.3.3
	gorm.io/driver/sqlite v1.3.1
	gorm.io/gorm v1.23.4
)
