package domain

type Config struct {
	HTTPTargets    []string `json:"http_targets"`
	ProcessTargets []string `json:"process_targets"`
}
