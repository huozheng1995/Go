package common

type ResData struct {
	Success bool        `json:"Success"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

type ConvertReq struct {
	ConvertType string `json:"ConvertType"`
	InputData   string `json:"InputData"`
}

type ConvertRes struct {
	OutputData string `json:"OutputData"`
}
