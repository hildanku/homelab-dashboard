package domain

type HTTPStatus struct {
	URL     string `json:"url"`
	OK      bool   `json:"ok"`
	Code    int    `json:"code"`
	Latency int64  `json:"latency_ms"`
}

type ProcStatus struct {
	Name  string `json:"name"`
	Found bool   `json:"found"`
	Count int    `json:"count"`
}
