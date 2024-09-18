package notification

import (
	"strconv"
	"strings"

	"github.com/zzucch/jimaku-tg-notify/internal/storage"
	"github.com/zzucch/jimaku-tg-notify/internal/util"
)

func getUpdateMessage(
	subscription storage.Subscription,
	update Update,
	err error,
) string {
	var sb strings.Builder

	if err != nil {
		sb.WriteString("Failed to get latest subtitle date for jimaku.cc/entry/")
		sb.WriteString(strconv.FormatInt(subscription.TitleID, 10))
		sb.WriteString(":\n")

		sb.WriteString(err.Error())
	} else if subscription.LastModified != update.LatestTimestamp {
		sb.WriteString(update.JapaneseName)
		sb.WriteString("\n")

		sb.WriteString("jimaku.cc/entry/")
		sb.WriteString(strconv.FormatInt(subscription.TitleID, 10))
		sb.WriteString(" at ")
		sb.WriteString(util.TimestampToString(update.LatestTimestamp))

		if len(update.NewFileEntryNames) > 0 {
			sb.WriteString("\nAdded:")
		}

		for _, name := range update.NewFileEntryNames {
			sb.WriteString("\n • ")
			sb.WriteString(name)
		}
	} else {
		return ""
	}

	sb.WriteString("\n\n")

	return sb.String()
}
