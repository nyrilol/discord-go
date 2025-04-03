// api/types.go
package types

import (
	"encoding/json"
	"io"
	"strconv"
)

type Snowflake string

func (s Snowflake) String() string {
	return string(s)
}

func (s Snowflake) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *Snowflake) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = Snowflake(str)
		return nil
	}

	var num int64
	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}
	*s = Snowflake(strconv.FormatInt(num, 10))
	return nil
}

func (s Snowflake) Equal(other Snowflake) bool {
	return s == other
}

func (s Snowflake) IsEmpty() bool {
	return s == ""
}

// or else guildmember is stupid
type Member = GuildMember

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

// WebhookMessage represents a message that can be sent via a webhook (used for interaction followups)
type WebhookMessage struct {
	Content         string             `json:"content,omitempty"`
	Username        string             `json:"username,omitempty"`
	AvatarURL       string             `json:"avatar_url,omitempty"`
	TTS             bool               `json:"tts,omitempty"`
	Embeds          []*Embed           `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions   `json:"allowed_mentions,omitempty"`
	Components      []MessageComponent `json:"components,omitempty"`
	Files           []*File            `json:"-"`
	PayloadJSON     string             `json:"payload_json,omitempty"`
	Attachments     []*Attachment      `json:"attachments,omitempty"`
	Flags           int                `json:"flags,omitempty"`
	ThreadName      string             `json:"thread_name,omitempty"`
}

// booboo
type File struct {
	Name        string
	ContentType string
	Reader      io.Reader
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
	Permissions  string   `json:"permissions"`
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

// ThreadMember struct
type ThreadMember struct {
	ID            string `json:"id,omitempty"`
	UserID        string `json:"user_id,omitempty"`
	JoinTimestamp string `json:"join_timestamp"`
	Flags         int    `json:"flags"`
	Member        Member `json:"member,omitempty"`
}

// Application struct
type Application struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Icon                string `json:"icon,omitempty"`
	Description         string `json:"description,omitempty"`
	BotPublic           bool   `json:"bot_public"`
	BotRequireCodeGrant bool   `json:"bot_require_code_grant"`
	Owner               User   `json:"owner,omitempty"`
	Flags               int    `json:"flags,omitempty"`
}

// ApplicationCommand struct
type ApplicationCommand struct {
	ID                string                     `json:"id"`
	Type              int                        `json:"type,omitempty"`
	ApplicationID     string                     `json:"application_id"`
	GuildID           string                     `json:"guild_id,omitempty"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                       `json:"default_permission,omitempty"`
	Version           string                     `json:"version"`
}

// ApplicationCommandOption struct
type ApplicationCommandOption struct {
	Type        int                              `json:"type"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Required    bool                             `json:"required,omitempty"`
	Choices     []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options     []ApplicationCommandOption       `json:"options,omitempty"`
}

// ApplicationCommandOptionChoice struct
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// StageInstance struct
type StageInstance struct {
	ID                   string `json:"id"`
	GuildID              string `json:"guild_id"`
	ChannelID            string `json:"channel_id"`
	Topic                string `json:"topic"`
	PrivacyLevel         int    `json:"privacy_level"`
	DiscoverableDisabled bool   `json:"discoverable_disabled"`
}

// AutoModerationRule struct
type AutoModerationRule struct {
	ID              string                        `json:"id"`
	GuildID         string                        `json:"guild_id"`
	Name            string                        `json:"name"`
	CreatorID       string                        `json:"creator_id"`
	EventType       int                           `json:"event_type"`
	TriggerType     int                           `json:"trigger_type"`
	TriggerMetadata AutoModerationTriggerMetadata `json:"trigger_metadata"`
	Actions         []AutoModerationAction        `json:"actions"`
	Enabled         bool                          `json:"enabled"`
	ExemptRoles     []string                      `json:"exempt_roles"`
	ExemptChannels  []string                      `json:"exempt_channels"`
}

// AutoModerationTriggerMetadata struct
type AutoModerationTriggerMetadata struct {
	KeywordFilter     []string `json:"keyword_filter,omitempty"`
	RegexPatterns     []string `json:"regex_patterns,omitempty"`
	Presets           []int    `json:"presets,omitempty"`
	AllowList         []string `json:"allow_list,omitempty"`
	MentionTotalLimit int      `json:"mention_total_limit,omitempty"`
}

// AutoModerationAction struct
type AutoModerationAction struct {
	Type     int                          `json:"type"`
	Metadata AutoModerationActionMetadata `json:"metadata,omitempty"`
}

// AutoModerationActionMetadata struct
type AutoModerationActionMetadata struct {
	ChannelID       string `json:"channel_id,omitempty"`
	DurationSeconds int    `json:"duration_seconds,omitempty"`
	CustomMessage   string `json:"custom_message,omitempty"`
}

// Integration struct
type Integration struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Type              string                 `json:"type"`
	Enabled           bool                   `json:"enabled"`
	Syncing           bool                   `json:"syncing,omitempty"`
	RoleID            string                 `json:"role_id,omitempty"`
	EnableEmoticons   bool                   `json:"enable_emoticons,omitempty"`
	ExpireBehavior    int                    `json:"expire_behavior,omitempty"`
	ExpireGracePeriod int                    `json:"expire_grace_period,omitempty"`
	User              User                   `json:"user,omitempty"`
	Account           IntegrationAccount     `json:"account"`
	SyncedAt          string                 `json:"synced_at,omitempty"`
	SubscriberCount   int                    `json:"subscriber_count,omitempty"`
	Revoked           bool                   `json:"revoked,omitempty"`
	Application       IntegrationApplication `json:"application,omitempty"`
}

// IntegrationAccount struct
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IntegrationApplication struct
type IntegrationApplication struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description"`
	Bot         User   `json:"bot,omitempty"`
}

// Entitlement struct
type Entitlement struct {
	ID            string `json:"id"`
	SkuID         string `json:"sku_id"`
	ApplicationID string `json:"application_id"`
	UserID        string `json:"user_id,omitempty"`
	GuildID       string `json:"guild_id,omitempty"`
	Type          int    `json:"type"`
	Consumed      bool   `json:"consumed,omitempty"`
	StartsAt      string `json:"starts_at,omitempty"`
	EndsAt        string `json:"ends_at,omitempty"`
}

// GuildJoinRequest struct
type GuildJoinRequest struct {
	UserID          string `json:"user_id"`
	GuildID         string `json:"guild_id"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

// InteractionType represents the type of interaction
const (
	InteractionTypePing                           = 1
	InteractionTypeApplicationCommand             = 2
	InteractionTypeMessageComponent               = 3
	InteractionTypeApplicationCommandAutocomplete = 4
	InteractionTypeModalSubmit                    = 5
)

// InteractionResponseType represents the type of response to an interaction
const (
	InteractionResponseTypePong                                 = 1
	InteractionResponseTypeChannelMessageWithSource             = 4
	InteractionResponseTypeDeferredChannelMessageWithSource     = 5
	InteractionResponseTypeDeferredUpdateMessage                = 6
	InteractionResponseTypeUpdateMessage                        = 7
	InteractionResponseTypeApplicationCommandAutocompleteResult = 8
	InteractionResponseTypeModal                                = 9
)

// ComponentType represents the type of component
const (
	ComponentTypeActionRow  = 1
	ComponentTypeButton     = 2
	ComponentTypeSelectMenu = 3
	ComponentTypeTextInput  = 4
)

// ButtonStyle represents the style of a button
const (
	ButtonStylePrimary   = 1
	ButtonStyleSecondary = 2
	ButtonStyleSuccess   = 3
	ButtonStyleDanger    = 4
	ButtonStyleLink      = 5
)

// TextInputStyle represents the style of a text input
const (
	TextInputStyleShort     = 1
	TextInputStyleParagraph = 2
)

// Interaction struct - enhanced version
type Interaction struct {
	ID            Snowflake       `json:"id"`
	ApplicationID Snowflake       `json:"application_id"`
	Type          int             `json:"type"`
	Data          InteractionData `json:"data,omitempty"`
	GuildID       string          `json:"guild_id,omitempty"`
	ChannelID     string          `json:"channel_id,omitempty"`
	Member        *GuildMember    `json:"member,omitempty"`
	User          *User           `json:"user,omitempty"`
	Token         string          `json:"token"`
	Version       int             `json:"version"`
	Message       *Message        `json:"message,omitempty"`
	Locale        string          `json:"locale,omitempty"`
	GuildLocale   string          `json:"guild_locale,omitempty"`
}

// InteractionData represents the data payload of an interaction
type InteractionData struct {
	ID            Snowflake                             `json:"id,omitempty"`
	Name          string                                `json:"name,omitempty"`
	Type          int                                   `json:"type,omitempty"`
	Resolved      *ResolvedData                         `json:"resolved,omitempty"`
	Options       []ApplicationCommandInteractionOption `json:"options,omitempty"`
	CustomID      string                                `json:"custom_id,omitempty"`
	ComponentType int                                   `json:"component_type,omitempty"`
	Values        []string                              `json:"values,omitempty"`
	TargetID      string                                `json:"target_id,omitempty"`
	Components    []MessageComponent                    `json:"components,omitempty"`
}

// ResolvedData contains resolved data for command options
type ResolvedData struct {
	Users       map[string]User        `json:"users,omitempty"`
	Members     map[string]GuildMember `json:"members,omitempty"`
	Roles       map[string]Role        `json:"roles,omitempty"`
	Channels    map[string]Channel     `json:"channels,omitempty"`
	Messages    map[string]Message     `json:"messages,omitempty"`
	Attachments map[string]Attachment  `json:"attachments,omitempty"`
}

// ApplicationCommandInteractionOption represents an option in an application command interaction
type ApplicationCommandInteractionOption struct {
	Name    string                                `json:"name"`
	Type    int                                   `json:"type"`
	Value   interface{}                           `json:"value,omitempty"`
	Options []ApplicationCommandInteractionOption `json:"options,omitempty"`
	Focused bool                                  `json:"focused,omitempty"`
}

// InteractionResponse represents a response to an interaction
type InteractionResponse struct {
	Type int                      `json:"type"`
	Data *InteractionCallbackData `json:"data,omitempty"`
}

// InteractionCallbackData represents the data in an interaction response
type InteractionCallbackData struct {
	TTS             bool               `json:"tts,omitempty"`
	Content         string             `json:"content,omitempty"`
	Embeds          []*Embed           `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions   `json:"allowed_mentions,omitempty"`
	Flags           int                `json:"flags,omitempty"`
	Components      []MessageComponent `json:"components,omitempty"`
	Attachments     []*Attachment      `json:"attachments,omitempty"`
	CustomID        string             `json:"custom_id,omitempty"`
	Title           string             `json:"title,omitempty"`
}

// MessageComponent represents a message component
type MessageComponent interface{}

// ActionRowComponent represents an action row component
type ActionRowComponent struct {
	Type       int                `json:"type"`
	Components []MessageComponent `json:"components"`
}

// ButtonComponent represents a button component
type ButtonComponent struct {
	Type     int    `json:"type"`
	Style    int    `json:"style"`
	Label    string `json:"label,omitempty"`
	Emoji    *Emoji `json:"emoji,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

// SelectMenuComponent represents a select menu component
type SelectMenuComponent struct {
	Type        int            `json:"type"`
	CustomID    string         `json:"custom_id"`
	Options     []SelectOption `json:"options,omitempty"`
	Placeholder string         `json:"placeholder,omitempty"`
	MinValues   *int           `json:"min_values,omitempty"`
	MaxValues   int            `json:"max_values,omitempty"`
	Disabled    bool           `json:"disabled,omitempty"`
}

// SelectOption represents an option in a select menu
type SelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Emoji       *Emoji `json:"emoji,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

// TextInputComponent represents a text input component (for modals)
type TextInputComponent struct {
	Type        int    `json:"type"`
	CustomID    string `json:"custom_id"`
	Style       int    `json:"style"`
	Label       string `json:"label"`
	MinLength   int    `json:"min_length,omitempty"`
	MaxLength   int    `json:"max_length,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Value       string `json:"value,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
}

// AllowedMentions controls mention parsing in messages
type AllowedMentions struct {
	Parse       []string `json:"parse,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Users       []string `json:"users,omitempty"`
	RepliedUser bool     `json:"replied_user,omitempty"`
}

type ApplicationCommandInteractionData struct {
	ID       int                                   `json:"id"`
	Name     string                                `json:"name"`
	Type     int                                   `json:"type"`
	Resolved *ResolvedData                         `json:"resolved,omitempty"`
	Options  []ApplicationCommandInteractionOption `json:"options,omitempty"`
	GuildID  string                                `json:"guild_id,omitempty"`
	TargetID string                                `json:"target_id,omitempty"`
}

// MessageFlags represents flags for a message
const (
	MessageFlagCrossposted          = 1 << 0
	MessageFlagIsCrosspost          = 1 << 1
	MessageFlagSuppressEmbeds       = 1 << 2
	MessageFlagSourceMessageDeleted = 1 << 3
	MessageFlagUrgent               = 1 << 4
	MessageFlagEphemeral            = 1 << 6
	MessageFlagLoading              = 1 << 7
)

// ApplicationCommandType represents the type of application command
const (
	ApplicationCommandTypeChatInput = 1
	ApplicationCommandTypeUser      = 2
	ApplicationCommandTypeMessage   = 3
)

// ApplicationCommandOptionType represents the type of application command option
const (
	ApplicationCommandOptionTypeSubCommand      = 1
	ApplicationCommandOptionTypeSubCommandGroup = 2
	ApplicationCommandOptionTypeString          = 3
	ApplicationCommandOptionTypeInteger         = 4
	ApplicationCommandOptionTypeBoolean         = 5
	ApplicationCommandOptionTypeUser            = 6
	ApplicationCommandOptionTypeChannel         = 7
	ApplicationCommandOptionTypeRole            = 8
	ApplicationCommandOptionTypeMentionable     = 9
	ApplicationCommandOptionTypeNumber          = 10
	ApplicationCommandOptionTypeAttachment      = 11
)
