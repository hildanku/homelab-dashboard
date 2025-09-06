package domain

type Snapshot struct {
	CPUUsage  float64 `json:"cpu_usage"`
	RAMTotal  uint64  `json:"ram_total"`
	RAMFree   uint64  `json:"ram_free"`
	RAMUsed   uint64  `json:"ram_used"`
	RAMUsage  float64 `json:"ram_usage"`
	DiskTotal uint64  `json:"disk_total"`
	DiskFree  uint64  `json:"disk_free"`
	DiskUsed  uint64  `json:"disk_used"`
	DiskUsage float64 `json:"disk_usage"`
}
