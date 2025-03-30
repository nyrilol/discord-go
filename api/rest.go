// api/rest.go
package api

import (
	"bytes"
	"discord-go/api/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const DiscordAPIURL = "https://discord.com/api/v10"

type Client struct {
	Token      string
	HTTPClient *http.Client
	RateLimits map[string]time.Time
	Mutex      sync.Mutex
}

func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		RateLimits: make(map[string]time.Time),
	}
}

func (c *Client) sendRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", DiscordAPIURL, endpoint)
	c.Mutex.Lock()
	if reset, exists := c.RateLimits[endpoint]; exists && time.Now().Before(reset) {
		delay := time.Until(reset)
		c.Mutex.Unlock()
		time.Sleep(delay)
	} else {
		c.Mutex.Unlock()
	}

	var requestBody []byte
	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bot "+c.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// add check for rate limit
	if resp.StatusCode == 429 {
		var rateLimit struct {
			RetryAfter float64 `json:"retry_after"`
		}
		json.NewDecoder(resp.Body).Decode(&rateLimit)
		c.Mutex.Lock()
		c.RateLimits[endpoint] = time.Now().Add(time.Duration(rateLimit.RetryAfter) * time.Second)
		c.Mutex.Unlock()
		resp.Body.Close()
		return c.sendRequest(method, endpoint, body)
	}

	if resetAfter := resp.Header.Get("X-RateLimit-Reset-After"); resetAfter != "" {
		if seconds, err := strconv.ParseFloat(resetAfter, 64); err == nil {
			c.Mutex.Lock()
			c.RateLimits[endpoint] = time.Now().Add(time.Duration(seconds) * time.Second)
			c.Mutex.Unlock()
		}
	}

	return resp, nil
}

func (c *Client) GetUser(userID string) (*types.User, error) {
	endpoint := fmt.Sprintf("/users/%s", userID)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) GetChannel(channelID string) (*types.Channel, error) {
	endpoint := fmt.Sprintf("/channels/%s", channelID)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var channel types.Channel
	if err := json.NewDecoder(resp.Body).Decode(&channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (c *Client) GetGuild(guildID string) (*types.Guild, error) {
	endpoint := fmt.Sprintf("/guilds/%s", guildID)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var guild types.Guild
	if err := json.NewDecoder(resp.Body).Decode(&guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (c *Client) GetMessages(channelID string, limit int) ([]types.Message, error) {
	endpoint := fmt.Sprintf("/channels/%s/messages?limit=%d", channelID, limit)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var messages []types.Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Client) CreateMessage(channelID string, content string) (*types.Message, error) {
	endpoint := fmt.Sprintf("/channels/%s/messages", channelID)
	messageData := struct {
		Content string `json:"content"`
	}{
		Content: content,
	}

	resp, err := c.sendRequest("POST", endpoint, messageData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var message types.Message
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (c *Client) DeleteMessage(channelID, messageID string) error {
	endpoint := fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID)
	resp, err := c.sendRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete message, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) CreateGuild(name, region string) (*types.Guild, error) {
	endpoint := "/guilds"
	guildData := struct {
		Name   string `json:"name"`
		Region string `json:"region"`
	}{
		Name:   name,
		Region: region,
	}

	resp, err := c.sendRequest("POST", endpoint, guildData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var guild types.Guild
	if err := json.NewDecoder(resp.Body).Decode(&guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (c *Client) AddGuildMember(guildID, userID, nickname string, roles []string) (*types.GuildMember, error) {
	endpoint := fmt.Sprintf("/guilds/%s/members/%s", guildID, userID)
	memberData := struct {
		Nickname string   `json:"nickname"`
		Roles    []string `json:"roles"`
	}{
		Nickname: nickname,
		Roles:    roles,
	}

	resp, err := c.sendRequest("PUT", endpoint, memberData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var member types.GuildMember
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return nil, err
	}

	return &member, nil
}

func (c *Client) GetGuildMembers(guildID string, limit int) ([]types.GuildMember, error) {
	endpoint := fmt.Sprintf("/guilds/%s/members?limit=%d", guildID, limit)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// try 2 decode as a slice of members
	var members []types.GuildMember
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		// got cooked ,single object 😒
		var member types.GuildMember
		if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
			return nil, fmt.Errorf("failed to decode guild members: %v", err)
		}
		members = append(members, member)
	}

	return members, nil
}

func (c *Client) GetInvite(inviteCode string) (*types.Invite, error) {
	endpoint := fmt.Sprintf("/invites/%s", inviteCode)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var invite types.Invite
	if err := json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}

	return &invite, nil
}

func (c *Client) CreateWebhook(channelID, name, avatar string) (*types.Webhook, error) {
	endpoint := fmt.Sprintf("/channels/%s/webhooks", channelID)
	webhookData := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}{
		Name:   name,
		Avatar: avatar,
	}

	resp, err := c.sendRequest("POST", endpoint, webhookData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var webhook types.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

func (c *Client) ModifyGuild(guildID, name string) (*types.Guild, error) {
	endpoint := fmt.Sprintf("/guilds/%s", guildID)
	guildData := struct {
		Name string `json:"name"`
	}{Name: name}

	resp, err := c.sendRequest("PATCH", endpoint, guildData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var guild types.Guild
	if err := json.NewDecoder(resp.Body).Decode(&guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (c *Client) AddReaction(channelID, messageID, emoji string) error {
	endpoint := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/@me", channelID, messageID, emoji)
	resp, err := c.sendRequest("PUT", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to add reaction, status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) DeleteReaction(channelID, messageID, emoji string) error {
	endpoint := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/@me", channelID, messageID, emoji)
	resp, err := c.sendRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete reaction, status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateDM(userID string) (string, error) {
	endpoint := "/users/@me/channels"
	dmData := struct {
		RecipientID string `json:"recipient_id"`
	}{RecipientID: userID}

	resp, err := c.sendRequest("POST", endpoint, dmData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var dmChannel struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&dmChannel); err != nil {
		return "", err
	}

	return dmChannel.ID, nil
}

func (c *Client) SendDM(userID, message string) (*types.Message, error) {
	dmID, err := c.CreateDM(userID)
	if err != nil {
		return nil, err
	}
	return c.CreateMessage(dmID, message)
}
