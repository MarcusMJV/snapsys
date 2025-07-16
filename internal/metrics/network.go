package metrics

type NetworkInterfaceStats struct {
	Interface string `json:"interface"`
	RxBytes   int    `json:"rx_bytes"`
	RxPackets int    `json:"rx_packets"`
	TxBytes   int    `json:"tx_bytes"`
	TxPackets int    `json:"tx_packets"`
}

type NetworkInterfaces []NetworkInterfaceStats
