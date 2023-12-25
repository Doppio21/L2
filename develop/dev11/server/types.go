package server

type GetRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Response struct {
	Events []Event
}


