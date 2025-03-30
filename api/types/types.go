// api/types.go
package types

import "encoding/json"

// User struct
type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	Email         string `json:"email,omitempty"`
	Verified      bool   `json:"verified,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Flags         int    `json:"flags"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
}

// Message struct
type Message struct {
	ID                string       `json:"id"`
	ChannelID         string       `json:"channel_id"`
	Content           string       `json:"content"`
	Timestamp         string       `json:"timestamp"`
	EditedTimestamp   string       `json:"edited_timestamp,omitempty"`
	Author            User         `json:"author"`
	Attachments       []Attachment `json:"attachments"`
	Embeds            []Embed      `json:"embeds"`
	Reactions         []Reaction   `json:"reactions"`
	MentionedUsers    []string     `json:"mention_user_ids"`
	MentionedRoles    []string     `json:"mention_role_ids"`
	MentionedChannels []string     `json:"mention_channel_ids"`
	MentionEveryone   bool         `json:"mention_everyone"`
	Pinned            bool         `json:"pinned"`
	TTS               bool         `json:"tts"`
}

// Attachment struct
type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// Embed struct
type Embed struct {
	Title       string          `json:"title,omitempty"`
	Type        string          `json:"type,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Timestamp   string          `json:"timestamp,omitempty"`
	Color       int             `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField    `json:"fields,omitempty"`
}

// EmbedFooter struct
type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// EmbedImage struct
type EmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// EmbedThumbnail struct
type EmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// EmbedVideo struct
type EmbedVideo struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// EmbedProvider struct
type EmbedProvider struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// EmbedAuthor struct
type EmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// EmbedField struct
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// Reaction struct
type Reaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

// Emoji struct
type Emoji struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Animated bool   `json:"animated"`
}

// Channel struct
type Channel struct {
	ID                   string                `json:"id"`
	Type                 int                   `json:"type"`
	GuildID              string                `json:"guild_id,omitempty"`
	Position             int                   `json:"position,omitempty"`
	Name                 string                `json:"name"`
	Topic                string                `json:"topic,omitempty"`
	NSFW                 bool                  `json:"nsfw"`
	LastMessageID        string                `json:"last_message_id,omitempty"`
	RateLimitPerUser     int                   `json:"rate_limit_per_user,omitempty"`
	UserLimit            int                   `json:"user_limit,omitempty"`
	Bitrate              int                   `json:"bitrate,omitempty"`
	ParentID             string                `json:"parent_id,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
}

// PermissionOverwrite struct
type PermissionOverwrite struct {
	ID    string `json:"id"`
	Type  int    `json:"type"`
	Allow int    `json:"allow"`
	Deny  int    `json:"deny"`
}

// Guild struct
type Guild struct {
	ID                          string        `json:"id"`
	Name                        string        `json:"name"`
	Icon                        string        `json:"icon,omitempty"`
	Splash                      string        `json:"splash,omitempty"`
	OwnerID                     string        `json:"owner_id"`
	Region                      string        `json:"region"`
	MemberCount                 int           `json:"member_count"`
	VerificationLevel           int           `json:"verification_level"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	Features                    []string      `json:"features,omitempty"`
	Emojis                      []Emoji       `json:"emojis,omitempty"`
	Roles                       []Role        `json:"roles"`
	Channels                    []Channel     `json:"channels"`
	Members                     []GuildMember `json:"members,omitempty"`
	VanityURLCode               string        `json:"vanity_url_code,omitempty"`
}

// Role struct
type Role struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Color       int             `json:"color"`
	Hoist       bool            `json:"hoist"`
	Position    int             `json:"position"`
	Permissions json.RawMessage `json:"permissions"`
	Managed     bool            `json:"managed"`
	Mentionable bool            `json:"mentionable"`
}

// Permission struct
type Permission struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Webhook struct
type Webhook struct {
	ID        string `json:"id"`
	Type      int    `json:"type"`
	GuildID   string `json:"guild_id,omitempty"`
	ChannelID string `json:"channel_id"`
	User      User   `json:"user"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Token     string `json:"token"`
	URL       string `json:"url"`
}

// GuildMember struct
type GuildMember struct {
	User         User     `json:"user"`
	Nickname     string   `json:"nickname,omitempty"`
	Roles        []string `json:"roles"`
	JoinedAt     string   `json:"joined_at"`
	PremiumSince string   `json:"premium_since,omitempty"`
	Deaf         bool     `json:"deaf"`
	Mute         bool     `json:"mute"`
	Pending      bool     `json:"pending"`
	Permissions  int      `json:"permissions"`
}

// GatewayEvent struct
type GatewayEvent struct {
	T  string          `json:"t"`
	D  json.RawMessage `json:"d"`
	S  int64           `json:"s"`
	OP int             `json:"op"`
}

// VoiceState struct
type VoiceState struct {
	GuildID    string      `json:"guild_id"`
	ChannelID  string      `json:"channel_id,omitempty"`
	UserID     string      `json:"user_id"`
	Member     GuildMember `json:"member,omitempty"`
	Deaf       bool        `json:"deaf"`
	Mute       bool        `json:"mute"`
	SelfDeaf   bool        `json:"self_deaf"`
	SelfMute   bool        `json:"self_mute"`
	Suppressed bool        `json:"suppress"`
	SessionID  string      `json:"session_id"`
}

// PresenceUpdate struct
type PresenceUpdate struct {
	User         User         `json:"user"`
	Status       string       `json:"status"`
	Activities   []Activity   `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
}

// Activity struct
type Activity struct {
	Name          string `json:"name"`
	Type          int    `json:"type"`
	URL           string `json:"url,omitempty"`
	Start         string `json:"start,omitempty"`
	End           string `json:"end,omitempty"`
	ApplicationID string `json:"application_id,omitempty"`
	State         string `json:"state,omitempty"`
	Details       string `json:"details,omitempty"`
}

// ClientStatus struct
type ClientStatus struct {
	Desktop string `json:"desktop,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
	Web     string `json:"web,omitempty"`
}

// Ban struct
type Ban struct {
	User   User   `json:"user"`
	Reason string `json:"reason,omitempty"`
}

// Invite struct
type Invite struct {
	Code           interface{} `json:"code"` //cannot unmarshal number into Go struct field Invite.code of type string
	Guild          Guild       `json:"guild"`
	Inviter        User        `json:"inviter"`
	TargetUser     User        `json:"target_user,omitempty"`
	TargetUserType int         `json:"target_user_type,omitempty"`
	MaxAge         int         `json:"max_age"`
	Uses           int         `json:"uses"`
	Temporary      bool        `json:"temporary"`
}

// GuildScheduledEvent struct
type GuildScheduledEvent struct {
	ID             string `json:"id"`
	GuildID        string `json:"guild_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ScheduledStart string `json:"scheduled_start"`
	ScheduledEnd   string `json:"scheduled_end"`
	EntityType     int    `json:"entity_type"`
	ChannelID      string `json:"channel_id"`
	UserCount      int    `json:"user_count"`
	PrivacyLevel   int    `json:"privacy_level"`
}

// Interaction struct
type Interaction struct {
	ID        string          `json:"id"`
	Type      int             `json:"type"`
	Data      json.RawMessage `json:"data"`
	GuildID   string          `json:"guild_id,omitempty"`
	ChannelID string          `json:"channel_id"`
	Member    GuildMember     `json:"member,omitempty"`
}
