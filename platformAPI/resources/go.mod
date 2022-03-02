module dev.com/resources

go 1.16

replace dev.com/utils => ../utils

require (
	dev.com/utils v0.0.0-00010101000000-000000000000
	github.com/labstack/echo/v4 v4.6.3
	go.mongodb.org/mongo-driver v1.7.2
)
