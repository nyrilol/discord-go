package types

import "time"

type MessageUpdateEvent struct {
	ID              string       `json:"id"`
	ChannelID       string       `json:"channel_id"`
	Content         string       `json:"content"`
	EditedTimestamp *time.Time   `json:"edited_timestamp"`
	Mentions        []User       `json:"mentions"`
	Attachments     []Attachment `json:"attachments"`
	Embeds          []Embed      `json:"embeds"`
}

type MessageDeleteEvent struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

// Reaction Events
type MessageReactionAddEvent struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Emoji     Emoji  `json:"emoji"`
}

type MessageReactionRemoveEvent struct {
	UserID    string `json:"user_id"`
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

type VoiceStateUpdateEvent struct {
	VoiceState
}

// Presence Events
type PresenceUpdateEvent struct {
	User       User       `json:"user"`
	Roles      []string   `json:"roles"`
	GuildID    string     `json:"guild_id"`
	Status     string     `json:"status"`
	Activities []Activity `json:"activities"`
}

// User Events
type UserUpdateEvent struct {
	User
}

// too lzfjkzgksdhgiosdhgoisd
type ReadyEvent struct {
	User      User   `json:"user"`
	Shard     []int  `json:"shard"`
	SessionID string `json:"session_id"`
}

// too lzfjkzgksdhgiosdhgoisd
type MessageCreateEvent struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Author    User   `json:"author"`
}
