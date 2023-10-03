package domain

type YoutubeVideoUploadPayload struct {
	VideoPath string `json:"video_path"`
	ClientKey string `json:"client_key"`
}
