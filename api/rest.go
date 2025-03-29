// api/rest.go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const DiscordAPIURL = "https://discord.com/api/v10"

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{Token: token}
}

func (c *Client) sendRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", DiscordAPIURL, endpoint)

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

	client := &http.Client{
		Timeout: time.Second * 10, // set a timeout to prevent hanging
	}

	return client.Do(req)
}

func (c *Client) GetUser(userID string) (*User, error) {
	endpoint := fmt.Sprintf("/users/%s", userID)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) GetChannel(channelID string) (*Channel, error) {
	endpoint := fmt.Sprintf("/channels/%s", channelID)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var channel Channel
	if err := json.NewDecoder(resp.Body).Decode(&channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (c *Client) GetGuild(guildID string) (*Guild, error) {
	endpoint := fmt.Sprintf("/guilds/%s", guildID)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var guild Guild
	if err := json.NewDecoder(resp.Body).Decode(&guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (c *Client) GetMessages(channelID string, limit int) ([]Message, error) {
	endpoint := fmt.Sprintf("/channels/%s/messages?limit=%d", channelID, limit)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Client) CreateMessage(channelID string, content string) (*Message, error) {
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

	var message Message
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

func (c *Client) CreateGuild(name, region string) (*Guild, error) {
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

	var guild Guild
	if err := json.NewDecoder(resp.Body).Decode(&guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (c *Client) AddGuildMember(guildID, userID, nickname string, roles []string) (*GuildMember, error) {
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

	var member GuildMember
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return nil, err
	}

	return &member, nil
}

func (c *Client) GetGuildMembers(guildID string, limit int) ([]GuildMember, error) {
	endpoint := fmt.Sprintf("/guilds/%s/members?limit=%d", guildID, limit)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// try 2 decode as a slice of members
	var members []GuildMember
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		// got cooked ,single object 😒
		var member GuildMember
		if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
			return nil, fmt.Errorf("failed to decode guild members: %v", err)
		}
		members = append(members, member)
	}

	return members, nil
}

func (c *Client) GetInvite(inviteCode string) (*Invite, error) {
	endpoint := fmt.Sprintf("/invites/%s", inviteCode)
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var invite Invite
	if err := json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}

	return &invite, nil
}

func (c *Client) CreateWebhook(channelID, name, avatar string) (*Webhook, error) {
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

	var webhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}
