package models

type SearchResult struct{
	Tokens []map[string]int `json:"tokens"`
	ExecutionTimes ExecutionTimes `json:"execution_times"`
	Data []Data `json:"data"`
}