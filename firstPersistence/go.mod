module dev.com/firstPersistence

go 1.16

replace (
	dev.com/config => ./config
	dev.com/utils => ./utils
)

require (
	dev.com/config v0.0.0-00010101000000-000000000000
	dev.com/utils v0.0.0-00010101000000-000000000000
	github.com/gobuffalo/genny v0.1.1 // indirect
	github.com/gobuffalo/gogen v0.1.1 // indirect
	github.com/karrick/godirwalk v1.10.3 // indirect
	github.com/pelletier/go-toml v1.7.0 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/streadway/amqp v1.0.0
	go.mongodb.org/mongo-driver v1.8.4
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/text v0.3.7 // indirect
)
