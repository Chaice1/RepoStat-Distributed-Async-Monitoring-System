package dto

type ServicesInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponse struct {
	Status   string         `json:"status"`
	Services []ServicesInfo `json:"services"`
}
