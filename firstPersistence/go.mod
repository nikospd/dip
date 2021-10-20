module dev.com/dataEntryAPI

go 1.16

replace (
	dev.com/config => ./config
	dev.com/utils => ./utils
)

require (
	dev.com/config v0.0.0-00010101000000-000000000000
	dev.com/utils v0.0.0-00010101000000-000000000000
	github.com/streadway/amqp v1.0.0
	go.mongodb.org/mongo-driver v1.7.2
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/text v0.3.7 // indirect
)
