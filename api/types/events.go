package types

import "time"

// Message Events
type MessageCreateEvent struct {
	ID              string       `json:"id"`
	ChannelID       string       `json:"channel_id"`
	Content         string       `json:"content"`
	Timestamp       time.Time    `json:"timestamp"`
	EditedTimestamp *time.Time   `json:"edited_timestamp"`
	Author          User         `json:"author"`
	Mentions        []User       `json:"mentions"`
	Attachments     []Attachment `json:"attachments"`
	Embeds          []Embed      `json:"embeds"`
	Pinned          bool         `json:"pinned"`
	Type            int          `json:"type"`
}

type MessageUpdateEvent struct {
	ID              string       `json:"id"`
	ChannelID       string       `json:"channel_id"`
	Content         string       `json:"content"`
	EditedTimestamp *time.Time   `json:"edited_timestamp"`
	Mentions        []User       `json:"mentions"`
	Attachments     []Attachment `json:"attachments"`
	Embeds          []Embed      `json:"embeds"`
	Pinned          bool         `json:"pinned"`
	Type            int          `json:"type"`
}

type MessageDeleteEvent struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

type MessageDeleteBulkEvent struct {
	IDs       []string `json:"ids"`
	ChannelID string   `json:"channel_id"`
	GuildID   string   `json:"guild_id,omitempty"`
}

// Reaction Events
type MessageReactionAddEvent struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Emoji     Emoji  `json:"emoji"`
	Member    Member `json:"member,omitempty"`
}

type MessageReactionRemoveEvent struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Emoji     Emoji  `json:"emoji"`
}

type MessageReactionRemoveAllEvent struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

type MessageReactionRemoveEmojiEvent struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Emoji     Emoji  `json:"emoji"`
}

// Guild Events
type GuildCreateEvent struct {
	Guild
}

type GuildUpdateEvent struct {
	Guild
}

type GuildDeleteEvent struct {
	ID          string `json:"id"`
	Unavailable bool   `json:"unavailable"`
}

// Guild Member Events
type GuildMemberAddEvent struct {
	GuildID string      `json:"guild_id"`
	Member  GuildMember `json:"member"`
}

type GuildMemberUpdateEvent struct {
	GuildID                    string     `json:"guild_id"`
	Roles                      []string   `json:"roles"`
	User                       User       `json:"user"`
	Nick                       string     `json:"nick,omitempty"`
	PremiumSince               string     `json:"premium_since,omitempty"`
	Pending                    bool       `json:"pending,omitempty"`
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until,omitempty"`
}

type GuildMemberRemoveEvent struct {
	GuildID string `json:"guild_id"`
	User    User   `json:"user"`
}

// Role Events
type GuildRoleCreateEvent struct {
	GuildID string `json:"guild_id"`
	Role    Role   `json:"role"`
}

type GuildRoleUpdateEvent struct {
	GuildID string `json:"guild_id"`
	Role    Role   `json:"role"`
}

type GuildRoleDeleteEvent struct {
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
}

// Channel Events
type ChannelCreateEvent struct {
	Channel
}

type ChannelUpdateEvent struct {
	Channel
}

type ChannelDeleteEvent struct {
	Channel
}

type ChannelPinsUpdateEvent struct {
	ChannelID        string     `json:"channel_id"`
	GuildID          string     `json:"guild_id,omitempty"`
	LastPinTimestamp *time.Time `json:"last_pin_timestamp,omitempty"`
}

// Thread Events
type ThreadCreateEvent struct {
	Channel
	NewlyCreated bool `json:"newly_created,omitempty"`
}

type ThreadUpdateEvent struct {
	Channel
}

type ThreadDeleteEvent struct {
	ID       string `json:"id"`
	GuildID  string `json:"guild_id"`
	ParentID string `json:"parent_id"`
	Type     int    `json:"type"`
}

type ThreadListSyncEvent struct {
	GuildID    string         `json:"guild_id"`
	ChannelIDs []string       `json:"channel_ids,omitempty"`
	Threads    []Channel      `json:"threads"`
	Members    []ThreadMember `json:"members"`
}

type ThreadMemberUpdateEvent struct {
	ThreadMember
	GuildID string `json:"guild_id"`
}

// Voice Events
type VoiceStateUpdateEvent struct {
	VoiceState
}

type VoiceServerUpdateEvent struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

// Presence Events
type PresenceUpdateEvent struct {
	User         User         `json:"user"`
	Roles        []string     `json:"roles"`
	GuildID      string       `json:"guild_id"`
	Status       string       `json:"status"`
	Activities   []Activity   `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
}

type ClientStatus struct {
	Desktop string `json:"desktop,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
	Web     string `json:"web,omitempty"`
}

// User Events
type UserUpdateEvent struct {
	User
}

// Ready Event
type ReadyEvent struct {
	V                 int         `json:"v"`
	User              User        `json:"user"`
	Guilds            []Guild     `json:"guilds"`
	SessionID         string      `json:"session_id"`
	Shard             []int       `json:"shard,omitempty"`
	Application       Application `json:"application"`
	ReadySupplemental struct {
		MergedPresences struct {
			Friends []PresenceUpdateEvent `json:"friends"`
			Guilds  []struct {
				ID        string                `json:"id"`
				Presences []PresenceUpdateEvent `json:"presences"`
			} `json:"guilds"`
		} `json:"merged_presences"`
	} `json:"ready_supplemental,omitempty"`
}

// Application Command Events
type ApplicationCommandCreateEvent struct {
	ApplicationCommand
	GuildID string `json:"guild_id,omitempty"`
}

type ApplicationCommandUpdateEvent struct {
	ApplicationCommand
	GuildID string `json:"guild_id,omitempty"`
}

type ApplicationCommandDeleteEvent struct {
	ID            string `json:"id"`
	GuildID       string `json:"guild_id,omitempty"`
	Type          int    `json:"type,omitempty"`
	ApplicationID string `json:"application_id"`
}

// Interaction Events
type InteractionCreateEvent struct {
	Interaction
}

// Invite Events
type InviteCreateEvent struct {
	ChannelID         string       `json:"channel_id"`
	Code              string       `json:"code"`
	CreatedAt         time.Time    `json:"created_at"`
	GuildID           string       `json:"guild_id,omitempty"`
	Inviter           User         `json:"inviter,omitempty"`
	MaxAge            int          `json:"max_age"`
	MaxUses           int          `json:"max_uses"`
	TargetType        int          `json:"target_type,omitempty"`
	TargetUser        User         `json:"target_user,omitempty"`
	TargetApplication *Application `json:"target_application,omitempty"`
	Temporary         bool         `json:"temporary"`
	Uses              int          `json:"uses"`
}

type InviteDeleteEvent struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Code      string `json:"code"`
}

// Stage Instance Events
type StageInstanceCreateEvent struct {
	StageInstance
}

type StageInstanceUpdateEvent struct {
	StageInstance
}

type StageInstanceDeleteEvent struct {
	StageInstance
}

// Guild Scheduled Events
type GuildScheduledEventCreateEvent struct {
	GuildScheduledEvent
}

type GuildScheduledEventUpdateEvent struct {
	GuildScheduledEvent
}

type GuildScheduledEventDeleteEvent struct {
	GuildScheduledEvent
}

type GuildScheduledEventUserAddEvent struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

type GuildScheduledEventUserRemoveEvent struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// Auto Moderation Events
type AutoModerationRuleCreateEvent struct {
	AutoModerationRule
}

type AutoModerationRuleUpdateEvent struct {
	AutoModerationRule
}

type AutoModerationRuleDeleteEvent struct {
	AutoModerationRule
}

type AutoModerationActionExecutionEvent struct {
	GuildID              string               `json:"guild_id"`
	Action               AutoModerationAction `json:"action"`
	RuleID               string               `json:"rule_id"`
	RuleTriggerType      int                  `json:"rule_trigger_type"`
	UserID               string               `json:"user_id"`
	ChannelID            string               `json:"channel_id,omitempty"`
	MessageID            string               `json:"message_id,omitempty"`
	AlertSystemMessageID string               `json:"alert_system_message_id,omitempty"`
	Content              string               `json:"content"`
	MatchedKeyword       string               `json:"matched_keyword"`
	MatchedContent       string               `json:"matched_content"`
}

// Typing Start Event
type TypingStartEvent struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	UserID    string `json:"user_id"`
	Timestamp int    `json:"timestamp"`
	Member    Member `json:"member,omitempty"`
}

// Webhooks Update Event
type WebhooksUpdateEvent struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

// Integration Events
type IntegrationCreateEvent struct {
	Integration
	GuildID string `json:"guild_id"`
}

type IntegrationUpdateEvent struct {
	Integration
	GuildID string `json:"guild_id"`
}

type IntegrationDeleteEvent struct {
	ID            string `json:"id"`
	GuildID       string `json:"guild_id"`
	ApplicationID string `json:"application_id,omitempty"`
}

// Entitlement Events
type EntitlementCreateEvent struct {
	Entitlement
}

type EntitlementUpdateEvent struct {
	Entitlement
}

type EntitlementDeleteEvent struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	UserID        string `json:"user_id,omitempty"`
	GuildID       string `json:"guild_id,omitempty"`
	SkuID         string `json:"sku_id"`
}

// Guild Join Request Events
type GuildJoinRequestCreateEvent struct {
	UserID      string           `json:"user_id"`
	GuildID     string           `json:"guild_id"`
	JoinRequest GuildJoinRequest `json:"join_request"`
}

type GuildJoinRequestUpdateEvent struct {
	UserID      string           `json:"user_id"`
	GuildID     string           `json:"guild_id"`
	JoinRequest GuildJoinRequest `json:"join_request"`
}

type GuildJoinRequestDeleteEvent struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
}
