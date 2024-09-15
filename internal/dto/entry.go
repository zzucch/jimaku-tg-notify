package dto

import "github.com/zzucch/jimaku-tg-notify/internal/util"

type Entry struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	LastModified string `json:"last_modified"`
	AnilistID    int64  `json:"anilist_id"`
	EnglishName  string `json:"english_name"`
	JapaneseName string `json:"japanese_name"`
}

func (e *Entry) GetLatestSubtitleTimestamp() (int64, error) {
	latestSubtitleTime, err := util.RFC3339ToUnixTimestamp(e.LastModified)
	if err != nil {
		return 0, err
	}

	return latestSubtitleTime, nil
}
