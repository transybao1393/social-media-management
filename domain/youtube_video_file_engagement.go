package domain

import (
	"time"

	"google.golang.org/api/youtube/v3"
)

type YoutubeFileUploadInfo struct {
	FileName        string                   `json:"file_name" bson:"file_name"`
	FileSize        int64                    `json:"file_size" bson:"file_size"`
	FileContentType string                   `json:"file_content_type" bson:"file_content_type"`
	VideoEngagement *youtube.VideoStatistics `json:"video_engagement" bson:"video_engagement"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
