package message

type HeartBeatMessage struct {
	ServiceName string `json:"service_name"`
	Ip          string `json:"ip"`
}
