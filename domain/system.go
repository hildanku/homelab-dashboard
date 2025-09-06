package domain

type Info struct {
	OSName    string `json:"os_name"`
	OSVersion string `json:"os_version"`
	Kernel    string `json:"kernel"`
	Arch      string `json:"arch"`
	GoVersion string `json:"go_version"`
}
