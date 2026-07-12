package models

type ExportRequest struct {
	InputPath string  `json:"inputPath,omitempty"`
	Pitch     int     `json:"pitch"`
	Speed     float64 `json:"speed"`
	Volume    float64 `json:"volume"`
	Bass      float64 `json:"bass"`
	Treble    float64 `json:"treble"`
	Echo      bool    `json:"echo"`
	Reverb    bool    `json:"reverb"`
	FadeIn    bool    `json:"fadeIn"`
	FadeOut   bool    `json:"fadeOut"`
	Normalize bool    `json:"normalize"`
}

type UploadResponse struct {
	ID         string `json:"id"`
	FileName   string `json:"fileName"`
	StoredName string `json:"storedName"`
	Size       int64  `json:"size"`
}

type ExportResponse struct {
	Success     bool   `json:"success"`
	DownloadURL string `json:"downloadUrl"`
	Message     string `json:"message,omitempty"`
}
