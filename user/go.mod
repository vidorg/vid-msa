module user

go 1.16

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/kataras/jwt v0.0.8
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/tal-tech/go-zero v1.1.1
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.8
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
