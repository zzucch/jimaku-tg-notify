package notification

import (
	"sort"

	"github.com/zzucch/jimaku-tg-notify/internal/storage"
	"github.com/zzucch/jimaku-tg-notify/pkg/client"
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

		sort.Slice(fileEntries, func(i, j int) bool {
			return fileEntries[i].LastModified.After(fileEntries[j].LastModified)
		})

		for _, fileEntry := range fileEntries {
			if fileEntry.LastModified.Unix() > subscription.LastModified {
				newFileEntryNames = append(newFileEntryNames, fileEntry.Name)
			}
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
