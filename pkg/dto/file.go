package dto

import "time"

type FileEntry struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"last_modified"`
	Size         int64     `json:"size"`
	URL          string    `json:"url"`
}
