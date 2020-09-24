module github.com/byuoitav/atlona-switcher-microservice

go 1.15

require (
	github.com/byuoitav/atlona-driver v1.5.7
	github.com/byuoitav/common v0.0.0-20200521193927-1fdf4e0a4271
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	gopkg.in/cas.v2 v2.2.0 // indirect
)

replace github.com/byuoitav/atlona-driver => ../atlona-driver/
