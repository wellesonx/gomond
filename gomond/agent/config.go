package agent

type NotificationConfig struct {
	// Enable notification for internal sys events
	Enable bool `json:"enable"`
	// Token for Telegram channel
	Token string `json:"token"`
}
