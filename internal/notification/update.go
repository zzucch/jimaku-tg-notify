package notification

import (
	"sort"

	"github.com/charmbracelet/log"
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

	var name string
	if (entry.JapaneseName != "" && subscription.Name == entry.JapaneseName) ||
		(entry.Name != "" && subscription.Name == entry.Name) ||
		(entry.EnglishName != "" && subscription.Name == entry.EnglishName) {
		name = ""
	} else {
		if entry.JapaneseName != "" {
			name = entry.JapaneseName
		} else if entry.Name != "" {
			name = entry.Name
		} else if entry.EnglishName != "" {
			name = entry.EnglishName
		}

		if name == "" {
			log.Error("failed to assign name value", "entry", entry)
		}
	}

	return Update{
		TitleID:           subscription.TitleID,
		LatestTimestamp:   lastModified,
		Name:              name,
		NewFileEntryNames: newFileEntryNames,
	}, nil
}
