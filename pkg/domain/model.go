package domain

// BGPSummary parsed summary monitoring data
type BGPSummary struct {
	LocalASNumber          int
	RIBEntries             int
	NumberOfPeers          int
	TotalNumberOfNeighbors int
}

type IXPServer struct {
	IXP     string
	City    string
	Country string
}

type QueryResult struct {
	Nonce  string `json:"nonce"`
	Status string `json:"status"`
	Result string `json:"result"`
}
