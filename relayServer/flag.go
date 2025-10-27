package relayserver_main

var (
	configFile string
)

type Config struct {
	Key                   string   `json:"key" yaml:"key"`
	Id                    string   `json:"id" yaml:"id"`
	PipeServerAddr        string   `json:"pipe_server_addr" yaml:"pipe_server_addr"`
	ListenRelayServerAddr string   `json:"listen_relay_server_addr" yaml:"listen_relay_server_addr"`
	WhiteIpList           []string `json:"white_ip_list" yaml:"white_ip_list"`
	DebugLog              bool     `json:"debug_log" yaml:"debug_log"`
	RegisterID            string   `json:"register_id" yaml:"register_id"`
	DirectConn            string   `json:"direct_conn" yaml:"direct_conn"`
}
