// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

//Config is the struct for all config variables
type Config struct {
	Brokers   []string `config:"brokers"`
	Topics    []string `config:"topics"`
	Group     string   `config:"group"`
	Partition int      `config:"partition"`
	//might want to add offset type
}

//DefaultConfig is the default values for configuration variables
var DefaultConfig = Config{
	Brokers:   []string{"localhost:9092"},
	Topics:    []string{"events"},
	Group:     "kafkabeat",
	Partition: 0,
}
