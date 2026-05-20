package models

type ErrorGeneral struct {
	Error      error
	StatusCode int
}

type ResponseGeneral struct {
	StatusCode int         `json:"statusCode"`
	Error      interface{} `json:"error"`
	Coleccion  interface{} `json:"data,omitempty"`
}
