module dev.com/pullMechanism

go 1.16

replace dev.com/utils => ./utils

replace dev.com/config => ./config

require (
	dev.com/config v0.0.0-00010101000000-000000000000
	dev.com/utils v0.0.0-00010101000000-000000000000
	github.com/streadway/amqp v1.0.0
	go.mongodb.org/mongo-driver v1.7.4
)
