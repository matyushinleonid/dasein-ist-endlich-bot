package model

type Session struct {
	Stage   int      `json:"stage"`
	Answers []string `json:"answers"`
}
