package models

type Resource struct {
	Name  string   `json:"name"`
	Host  string   `json:"host"`
	Ports []string `json:"ports"`
}
