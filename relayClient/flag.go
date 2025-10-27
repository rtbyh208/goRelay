package relayclient_main

var (
	configFile string
)

type Config struct {
	Key            string       `json:"key" yaml:"key"`
	RegisterID     string       `json:"register_id" yaml:"register_id"`
	PipeServerAddr string       `json:"pipe_server_addr" yaml:"pipe_server_addr"`
	DebugLog       bool         `json:"debug_log" yaml:"debug_log"`
	RealServerInfo []RealServer `json:"realServer_info" yaml:"realServer_info"`
	HttpServer     string       `json:"http_server" yaml:"http_server"`
}

type RealServer struct {
	ID             string `json:"id" yaml:"id"`
	RealServerAddr string `json:"real_Server_Addr" yaml:"real_Server_Addr"`
}
