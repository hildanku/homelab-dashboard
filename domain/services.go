package domain

type HTTPStatus struct {
	URL     string `json:"url"`
	OK      bool   `json:"ok"`
	Code    int    `json:"code"`
	Latency int64  `json:"latency_ms"`
}
