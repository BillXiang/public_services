package main

type CronJob struct {
	Exp string `json:"exp"`
	Cmd string `json:"cmd"`
}
