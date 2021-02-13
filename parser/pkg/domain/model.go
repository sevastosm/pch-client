package domain

// BGPSummary parsed summary monitoring data
type BGPSummary struct {
	LocalASNumber          int
	RIBEntries             int
	NumberOfPeers          int
	TotalNumberOfNeighbors int
}

// IXP Server per protocol details
type IXPServer struct {
	IXP      string
	City     string
	Country  string
	Protocol string //IPv4 or IPv6
}

type QueryResult struct {
	Nonce  string `json:"nonce"`
	Status string `json:"status"`
	Result string `json:"result"`
}
