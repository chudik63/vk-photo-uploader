package entity

type Photo struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Size         uint64 `json:"size"`
	Type         string `json:"type"`
	LastModified uint64 `json:"lastmodified"`
}
