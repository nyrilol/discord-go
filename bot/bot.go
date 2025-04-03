package bot

import (
	"discord-go/api"
	"discord-go/api/types"
	"discord-go/gateway"
	"discord-go/utils"
	"strings"
	"sync"
)

type Bot struct {
	token    string
	gateway  *gateway.Gateway
	logger   utils.Logger
	handlers sync.Map
	commands map[string]CommandHandler
	mu       sync.RWMutex
}

type InteractionHandler struct {
	bot            *Bot
	commands       map[string]CommandHandler
	buttons        map[string]ButtonHandler
	selectMenus    map[string]SelectMenuHandler
	modals         map[string]ModalHandler
	componentMutex sync.Mutex
	commandMutex   sync.Mutex
	globalCommands map[string]types.ApplicationCommand
	guildCommands  map[types.Snowflake]map[string]types.ApplicationCommand // guildID -> commandName -> command
}

type CommandHandler func(ctx *CommandContext)
type ButtonHandler func(ctx *ComponentContext)
type SelectMenuHandler func(ctx *ComponentContext)
type ModalHandler func(ctx *ModalContext)
type CommandContext struct {
	Interaction *types.Interaction
	Bot         *Bot
	Options     map[string]interface{}
}
type ComponentContext struct {
	Interaction *types.Interaction
	Bot         *Bot
	CustomID    string
	Values      []string
}
type ModalContext struct {
	Interaction *types.Interaction
	Bot         *Bot
	CustomID    string
	Inputs      map[string]string
}

type MessageContext struct {
	Message *types.Message
	Bot     *Bot
}

func NewBot(token string, intents ...int) *Bot {
	if !isValidToken(token) {
		panic("Token not in right format")
	}

	intent_value := api.IntentAll
	if len(intents) > 0 {
		intent_value = intents[0]
	}

	gw := gateway.NewGateway(token, intent_value)
	bot := &Bot{
		token:    token,
		gateway:  gw,
		logger:   utils.NewLogger(),
		commands: make(map[string]CommandHandler),
	}

	bot.registerDefaultHandlers()
	return bot
}

func isValidToken(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
	}
	return true
}

func (bot *Bot) SendInteractionResponse(interactionID types.Snowflake, token string, response types.InteractionResponse) error {
	return bot.gateway.SendInteractionResponse(interactionID, token, response)
}

func (bot *Bot) CreateGlobalApplicationCommand(command types.ApplicationCommand) error {
	return bot.gateway.CreateGlobalApplicationCommand(command)
}

func (bot *Bot) registerDefaultHandlers() {
	bot.On("READY", func(event types.ReadyEvent) {
		bot.logger.Infof("Bot is ready: %s (Shard %d)", event.User.Username, event.Shard)
		bot.gateway.RemoveHandler("READY")
	}, types.ReadyEvent{})
}

func (b *Bot) RegisterCommand(name string, handler CommandHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.commands[name] = handler

	if _, exists := b.commands["__internal_handler_registered"]; !exists {
		b.gateway.RegisterHandler("INTERACTION_CREATE", b.handleInteraction, types.Interaction{})
		b.commands["__internal_handler_registered"] = nil // mark as registered
	}
}

func (b *Bot) handleInteraction(interaction types.Interaction) {
	if interaction.Type != types.InteractionTypeApplicationCommand {
		return
	}

	b.mu.RLock()
	handler, exists := b.commands[interaction.Data.Name]
	b.mu.RUnlock()

	if !exists {
		return
	}

	options := make(map[string]interface{})
	for _, option := range interaction.Data.Options {
		options[option.Name] = option.Value
	}

	ctx := &CommandContext{
		Interaction: &interaction,
		Bot:         b,
		Options:     options,
	}

	handler(ctx)
}

func (b *Bot) AddMessageHandler(prefix string, handler func(*MessageContext)) {
	b.On("MESSAGE_CREATE", func(event types.Message) {
		if event.Content == prefix {
			handler(&MessageContext{
				Message: &event,
				Bot:     b,
			})
		}
	}, types.Message{})
}

func (b *Bot) AddSlashCommand(name, description string, handler func(*CommandContext)) {
	cmd := types.ApplicationCommand{
		Name:        name,
		Description: description,
		Type:        types.ApplicationCommandTypeChatInput,
	}

	if err := b.CreateGlobalApplicationCommand(cmd); err != nil {
		b.logger.Errorf("Failed to register %s command: %v", name, err)
		return
	}

	b.RegisterCommand(name, handler)
}

// component things
func (b *Bot) NewActionRow(components ...types.MessageComponent) types.ActionRowComponent {
	return types.ActionRowComponent{
		Type:       types.ComponentTypeActionRow,
		Components: components,
	}
}

func (b *Bot) NewButton(label, customID string, style int) types.ButtonComponent {
	return types.ButtonComponent{
		Type:     types.ComponentTypeButton,
		Style:    style,
		Label:    label,
		CustomID: customID,
	}
}

func (b *Bot) NewSelectMenu(customID, placeholder string, options []types.SelectOption) types.SelectMenuComponent {
	return types.SelectMenuComponent{
		Type:        types.ComponentTypeSelectMenu,
		CustomID:    customID,
		Placeholder: placeholder,
		Options:     options,
	}
}

// response helpers
func (b *Bot) RespondToInteraction(ctx *CommandContext, content string) error {
	response := types.InteractionResponse{
		Type: types.InteractionResponseTypeChannelMessageWithSource,
		Data: &types.InteractionCallbackData{
			Content: content,
		},
	}
	return b.SendInteractionResponse(ctx.Interaction.ID, ctx.Interaction.Token, response)
}

func (b *Bot) RespondWithComponents(ctx *CommandContext, content string, components []types.MessageComponent) error {
	response := types.InteractionResponse{
		Type: types.InteractionResponseTypeChannelMessageWithSource,
		Data: &types.InteractionCallbackData{
			Content:    content,
			Components: components,
		},
	}
	return b.SendInteractionResponse(ctx.Interaction.ID, ctx.Interaction.Token, response)
}

// modal helpers
func (b *Bot) NewModal(customID, title string, rows ...types.ActionRowComponent) types.InteractionResponse {
	components := make([]types.MessageComponent, len(rows))
	for i, row := range rows {
		components[i] = row
	}

	return types.InteractionResponse{
		Type: types.InteractionResponseTypeModal,
		Data: &types.InteractionCallbackData{
			CustomID:   customID,
			Title:      title,
			Components: components,
		},
	}
}

func (b *Bot) NewTextInputRow(customID, label string, style int, required bool, maxLength int, placeholder string) types.ActionRowComponent {
	return b.NewActionRow(types.TextInputComponent{
		Type:        types.ComponentTypeTextInput,
		CustomID:    customID,
		Label:       label,
		Style:       style,
		Placeholder: placeholder,
		Required:    required,
		MaxLength:   maxLength,
	})
}

func (b *Bot) RespondWithModal(ctx *ComponentContext, modal types.InteractionResponse) error {
	return b.SendInteractionResponse(ctx.Interaction.ID, ctx.Interaction.Token, modal)
}

// cpompent handlers
func (b *Bot) AddComponentHandler(customID string, handler func(*ComponentContext)) {
	b.On("INTERACTION_CREATE", func(interaction types.Interaction) {
		if interaction.Type == types.InteractionTypeMessageComponent && interaction.Data.CustomID == customID {
			var values []string
			if interaction.Data.Values != nil {
				values = interaction.Data.Values
			}

			handler(&ComponentContext{
				Interaction: &interaction,
				Bot:         b,
				CustomID:    customID,
				Values:      values,
			})
		}
	}, types.Interaction{})
}

// modal handler
func (b *Bot) SetModalHandler(handler func(*ModalContext)) {
	b.On("INTERACTION_CREATE", func(interaction types.Interaction) {
		if interaction.Type == types.InteractionTypeModalSubmit {
			inputs := make(map[string]string)
			for _, row := range interaction.Data.Components {
				if actionRow, ok := row.(types.ActionRowComponent); ok {
					for _, component := range actionRow.Components {
						if textInput, ok := component.(types.TextInputComponent); ok {
							inputs[textInput.CustomID] = textInput.Value
						}
					}
				}
			}

			handler(&ModalContext{
				Interaction: &interaction,
				Bot:         b,
				CustomID:    interaction.Data.CustomID,
				Inputs:      inputs,
			})
		}
	}, types.Interaction{})
}

func (bot *Bot) Start() {
	bot.logger.Info("Starting bot...")
	bot.gateway.Connect("wss://gateway.discord.gg/?v=10&encoding=json")
	select {}
}

func (bot *Bot) On(eventName string, handler interface{}, event_type interface{}) {
	eventName = strings.ToUpper(eventName) // just incase retard user
	bot.handlers.Store(eventName, handler)
	bot.gateway.RegisterHandler(eventName, handler, event_type)
}

func (bot *Bot) RemoveHandler(eventName string) {
	eventName = strings.ToUpper(eventName) // just incase retard user
	bot.handlers.Delete(eventName)
	bot.gateway.RemoveHandler(eventName)
}
