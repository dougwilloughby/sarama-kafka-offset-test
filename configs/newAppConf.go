package configs

// Configuration properties for the application. Once initialized, this information should
// be treated as immutable since it may be used by various go routines. Only properties that
// should go in here are those related to setting up the application
type AppConf struct {
	SrcTopic               string
	SrcPartition           int32
	CustId                 string
	AgentId                string
	Brokers                []string
	NumPartitions          int32
	ReplicationFactor      int16
}


// NewAppConf returns a new configuration initialized to "sensible" default values. Here, sensible
// means that the application won't panic or go off into the weeds if the values aren't defined.
// This doesn't mean that the app will work...
func NewAppConf() *AppConf {
	return &AppConf{
		SrcTopic:          "history-transactions",
		AgentId:           "f6d6d81c-6104-457b-acae-b8949ae4cecd", // Demo account
		CustId:            "ce8e6040-4e29-405f-bdfd-015ae91b11a3",
		Brokers:           []string{"bigdata-worker2.dc.res0.local:9092", "bigdata-worker1.dc.res0.local:9092", "bigdata-worker0.dc.res0.local:9092"},
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
}
