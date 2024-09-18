package notification

import (
	"github.com/zzucch/jimaku-tg-notify/internal/client"
	"github.com/zzucch/jimaku-tg-notify/internal/storage"
)

func getUpdate(
	subscription storage.Subscription,
	client *client.Client,
) (Update, error) {
	entry, err := client.GetEntryDetails(subscription.TitleID)
	if err != nil {
		return Update{}, err
	}

	lastModified, err := entry.GetLastModified()
	if err != nil {
		return Update{}, err
	}

	var newFileEntryNames []string

	if subscription.LastModified == lastModified {
		lastModified = -1
	} else {
		fileEntries, err := client.GetFileEntries(subscription.TitleID)
		if err != nil {
			return Update{}, err
		}

		for _, fileEntry := range fileEntries {
			newFileEntryNames = append(newFileEntryNames, fileEntry.Name)
		}
	}

	var japaneseName string
	if subscription.JapaneseName == entry.JapaneseName {
		japaneseName = ""
	} else {
		japaneseName = entry.JapaneseName
	}

	return Update{
		TitleID:           subscription.TitleID,
		LatestTimestamp:   lastModified,
		JapaneseName:      japaneseName,
		NewFileEntryNames: newFileEntryNames,
	}, nil
}
