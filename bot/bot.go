package bot

type BOT interface {
	SendMessage(messages []string, channelID int64)
}
