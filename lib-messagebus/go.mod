module github.com/ODIM-Project/ODIM/lib-messagebus

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/segmentio/kafka-go v0.3.5
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.1.0 // indirect
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
