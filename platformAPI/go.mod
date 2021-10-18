module dev.com/platformAPI

go 1.16

replace dev.com/utils => ./utils

replace dev.com/resources => ./resources

require (
	dev.com/resources v0.0.0-00010101000000-000000000000
	dev.com/utils v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/labstack/echo/v4 v4.6.0
	go.mongodb.org/mongo-driver v1.7.2
)
