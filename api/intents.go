package api

const (
	IntentGuilds                      = 1 << 0
	IntentGuildMembers                = 1 << 1 // privleged
	IntentGuildModeration             = 1 << 2
	IntentGuildEmojisAndStickers      = 1 << 3
	IntentGuildIntegrations           = 1 << 4
	IntentGuildWebhooks               = 1 << 5
	IntentGuildInvites                = 1 << 6
	IntentGuildVoiceStates            = 1 << 7
	IntentGuildPresences              = 1 << 8 // privleged
	IntentGuildMessages               = 1 << 9
	IntentGuildMessageReactions       = 1 << 10
	IntentGuildMessageTyping          = 1 << 11
	IntentDirectMessages              = 1 << 12
	IntentDirectMessageReactions      = 1 << 13
	IntentDirectMessageTyping         = 1 << 14
	IntentMessageContent              = 1 << 15 // privleged
	IntentGuildScheduledEvents        = 1 << 16
	IntentAutoModerationConfiguration = 1 << 20
	IntentAutoModerationExecution     = 1 << 21
	IntentGuildMessagePolls           = 1 << 22
	IntentDirectMessagePolls          = 1 << 23

	IntentAllNonPrivileged = IntentGuilds |
		IntentGuildModeration |
		IntentGuildEmojisAndStickers |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentGuildScheduledEvents |
		IntentAutoModerationConfiguration |
		IntentAutoModerationExecution |
		IntentGuildMessagePolls |
		IntentDirectMessagePolls

	IntentAll = IntentAllNonPrivileged | IntentGuildMembers | IntentGuildPresences | IntentMessageContent
)
