package queries

//use key value to build query object
var kvLogParser map[string]map[string]func(map[string]string) QueryInterface

//use json to build query object
var jsonLogParser map[string]map[string]func([]byte) QueryInterface

// initializes accepted types of logs
func init() {
	kvLogParser = map[string]map[string]func(map[string]string) QueryInterface{
		"bro": map[string]func(map[string]string) QueryInterface{
			"conn": BroConnQueryFromKV,
		},
	}

	jsonLogParser = map[string]map[string]func([]byte) QueryInterface{
		"bro": map[string]func([]byte) QueryInterface{
			"conn": BroConnQueryFromJSON,
		},
	}
}
