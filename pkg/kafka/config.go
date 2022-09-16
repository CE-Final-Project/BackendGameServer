package kafka

// Config kafka config
type Config struct {
	Brokers    []string
	GroupID    string
	InitTopics bool
}

// TopicConfig kafka topic config
type TopicConfig struct {
	TopicName         string
	Partitions        int
	ReplicationFactor int
}
