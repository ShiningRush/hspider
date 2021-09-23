module github.com/shiningrush/hspider

go 1.12

require (
	github.com/emicklei/go-restful/v3 v3.6.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/shiningrush/droplet v0.2.5
	github.com/shiningrush/droplet/wrapper/gorestful v0.1.0
	github.com/shiningrush/goreq v0.1.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace (
	github.com/shiningrush/goreq => ../../github/goreq
	github.com/shiningrush/droplet => ../../github/droplet
	github.com/shiningrush/droplet/wrapper/gorestful => ../../github/droplet/wrapper/gorestful
)
