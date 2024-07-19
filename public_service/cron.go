package public_service

type CronJob struct {
	Exp string `json:"exp"`
	Cmd string `json:"cmd"`
}
