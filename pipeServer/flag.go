package pipeserver_main

type Config struct {
	Key                  string   `json:"Key" yaml:"Key"`
	ListenPipeServerAddr string   `json:"addr" yaml:"addr"`
	WhiteIpList          []string `json:"white_ip_list" yaml:"white_ip_list"`
	BlackIpList          []string `json:"black_ip_list" yaml:"black_ip_list"`
	DebugLog             bool     `json:"debug_log" yaml:"debug_log"`
}
