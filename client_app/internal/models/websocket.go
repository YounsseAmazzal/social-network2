package models 

type WebSocketInput struct{
	WsData map[string]interface{} `json:"wsdata"`
	Datatype string  `json:"wsdatatype"`  
}