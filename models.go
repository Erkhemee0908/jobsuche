package main

type Job struct {
	Titel       string `json:"titel"`
	Refnum      string `json:"refnr"`
	Arbeitgeber string `json:"arbeitgeber"`
	Description string `json:"description"`
}
