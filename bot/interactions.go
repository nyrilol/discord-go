package bot

import (
	"discord-go/api/types"
	"fmt"
)

func NewInteractionHandler(bot *Bot) *InteractionHandler {
	return &InteractionHandler{
		bot:            bot,
		commands:       make(map[string]CommandHandler),
		buttons:        make(map[string]ButtonHandler),
		selectMenus:    make(map[string]SelectMenuHandler),
		modals:         make(map[string]ModalHandler),
		globalCommands: make(map[string]types.ApplicationCommand),
		guildCommands:  make(map[types.Snowflake]map[string]types.ApplicationCommand),
	}
}

func (ih *InteractionHandler) Initialize() {
	ih.bot.On("INTERACTION_CREATE", ih.handleInteraction, types.Interaction{})
}

func (ih *InteractionHandler) handleInteraction(interaction types.Interaction) {
	switch interaction.Type {
	case types.InteractionTypeApplicationCommand:
		ih.handleCommand(&interaction)
	case types.InteractionTypeMessageComponent:
		ih.handleComponent(&interaction)
	case types.InteractionTypeModalSubmit:
		ih.handleModal(&interaction)
	default:
		ih.bot.logger.Warnf("Unhandled interaction type: %d", interaction.Type)
	}
}

func (ih *InteractionHandler) handleCommand(interaction *types.Interaction) {
	commandName := interaction.Data.Name
	handler, exists := ih.commands[commandName]
	if !exists {
		ih.bot.logger.Warnf("No handler found for command: %s", commandName)
		return
	}

	ctx := &CommandContext{
		Interaction: interaction,
		Bot:         ih.bot,
		Options:     make(map[string]interface{}),
	}

	if interaction.Data.Options != nil {
		for _, option := range interaction.Data.Options {
			ctx.Options[option.Name] = option.Value
		}
	}

	handler(ctx)
}

func (ih *InteractionHandler) handleComponent(interaction *types.Interaction) {
	customID := interaction.Data.CustomID
	var handler interface{}
	var exists bool

	switch interaction.Data.ComponentType {
	case types.ComponentTypeButton:
		handler, exists = ih.buttons[customID]
	case types.ComponentTypeSelectMenu:
		handler, exists = ih.selectMenus[customID]
	default:
		ih.bot.logger.Warnf("Unhandled component type: %d", interaction.Data.ComponentType)
		return
	}

	if !exists {
		ih.bot.logger.Warnf("No handler found for component with custom ID: %s", customID)
		return
	}

	ctx := &ComponentContext{
		Interaction: interaction,
		Bot:         ih.bot,
		CustomID:    customID,
		Values:      interaction.Data.Values,
	}

	switch h := handler.(type) {
	case ButtonHandler:
		h(ctx)
	case SelectMenuHandler:
		h(ctx)
	}
}

func (ih *InteractionHandler) handleModal(interaction *types.Interaction) {
	customID := interaction.Data.CustomID
	handler, exists := ih.modals[customID]
	if !exists {
		ih.bot.logger.Warnf("No handler found for modal with custom ID: %s", customID)
		return
	}

	inputs := make(map[string]string)
	for _, row := range interaction.Data.Components {
		actionRow, ok := row.(types.ActionRowComponent)
		if !ok {
			continue
		}

		for _, component := range actionRow.Components {
			if input, ok := component.(types.TextInputComponent); ok {
				inputs[input.CustomID] = input.Value
			}
		}
	}

	ctx := &ModalContext{
		Interaction: interaction,
		Bot:         ih.bot,
		CustomID:    customID,
		Inputs:      inputs,
	}

	handler(ctx)
}

func (ih *InteractionHandler) Command(name string, handler CommandHandler) {
	ih.commandMutex.Lock()
	defer ih.commandMutex.Unlock()
	ih.commands[name] = handler
}

func (ih *InteractionHandler) Button(customID string, handler ButtonHandler) {
	ih.componentMutex.Lock()
	defer ih.componentMutex.Unlock()
	ih.buttons[customID] = handler
}

func (ih *InteractionHandler) SelectMenu(customID string, handler SelectMenuHandler) {
	ih.componentMutex.Lock()
	defer ih.componentMutex.Unlock()
	ih.selectMenus[customID] = handler
}

func (ih *InteractionHandler) Modal(customID string, handler ModalHandler) {
	ih.componentMutex.Lock()
	defer ih.componentMutex.Unlock()
	ih.modals[customID] = handler
}

func (ctx *CommandContext) Respond(content string, ephemeral ...bool) error {
	return ctx.respond(types.InteractionResponseTypeChannelMessageWithSource, content, ephemeral...)
}

func (ctx *CommandContext) Defer(ephemeral ...bool) error {
	return ctx.respond(types.InteractionResponseTypeDeferredChannelMessageWithSource, "", ephemeral...)
}

func (ctx *CommandContext) respond(responseType int, content string, ephemeral ...bool) error {
	flags := 0
	if len(ephemeral) > 0 && ephemeral[0] {
		flags = types.MessageFlagEphemeral
	}

	response := types.InteractionResponse{
		Type: responseType,
		Data: &types.InteractionCallbackData{
			Content: content,
			Flags:   flags,
		},
	}

	return ctx.Bot.gateway.SendInteractionResponse(ctx.Interaction.ID, ctx.Interaction.Token, response)
}

func (ctx *CommandContext) Followup(content string, ephemeral ...bool) error {
	flags := 0
	if len(ephemeral) > 0 && ephemeral[0] {
		flags = types.MessageFlagEphemeral
	}

	message := types.WebhookMessage{
		Content: content,
		Flags:   flags,
	}

	return ctx.Bot.gateway.SendFollowupMessage(ctx.Interaction.Token, message)
}

func (ctx *CommandContext) EditResponse(content string) error {
	return ctx.Bot.gateway.EditOriginalInteractionResponse(ctx.Interaction.Token, content)
}

func (ctx *CommandContext) CreateModal(modal *Modal) error {
	response := types.InteractionResponse{
		Type: types.InteractionResponseTypeModal,
		Data: &types.InteractionCallbackData{
			CustomID: modal.CustomID,
			Title:    modal.Title,
			Components: []types.MessageComponent{
				types.ActionRowComponent{
					Type:       types.ComponentTypeActionRow,
					Components: modal.Components,
				},
			},
		},
	}

	return ctx.Bot.gateway.SendInteractionResponse(ctx.Interaction.ID, ctx.Interaction.Token, response)
}

type Modal struct {
	CustomID   string
	Title      string
	Components []types.MessageComponent
}

func (m *Modal) AddTextInput(customID, label string, style int, options ...TextInputOption) {
	input := types.TextInputComponent{
		Type:     types.ComponentTypeTextInput,
		CustomID: customID,
		Style:    style,
		Label:    label,
	}

	for _, opt := range options {
		opt(&input)
	}

	m.Components = append(m.Components, input)
}

type TextInputOption func(*types.TextInputComponent)

func WithPlaceholder(placeholder string) TextInputOption {
	return func(input *types.TextInputComponent) {
		input.Placeholder = placeholder
	}
}

func WithMinLength(min int) TextInputOption {
	return func(input *types.TextInputComponent) {
		input.MinLength = min
	}
}

func WithMaxLength(max int) TextInputOption {
	return func(input *types.TextInputComponent) {
		input.MaxLength = max
	}
}

func WithRequired(required bool) TextInputOption {
	return func(input *types.TextInputComponent) {
		input.Required = required
	}
}

func WithDefaultValue(value string) TextInputOption {
	return func(input *types.TextInputComponent) {
		input.Value = value
	}
}

func (ih *InteractionHandler) RegisterCommand(command types.ApplicationCommand, guildID ...types.Snowflake) error {
	if len(guildID) > 0 {
		return ih.registerGuildCommand(command, guildID[0])
	}
	return ih.registerGlobalCommand(command)
}

func (ih *InteractionHandler) registerGlobalCommand(command types.ApplicationCommand) error {
	ih.commandMutex.Lock()
	defer ih.commandMutex.Unlock()

	if _, exists := ih.globalCommands[command.Name]; exists {
		return fmt.Errorf("command %s already registered globally", command.Name)
	}

	err := ih.bot.gateway.CreateGlobalApplicationCommand(command)
	if err != nil {
		return err
	}

	ih.globalCommands[command.Name] = command
	return nil
}

func (ih *InteractionHandler) registerGuildCommand(command types.ApplicationCommand, guildID types.Snowflake) error {
	ih.commandMutex.Lock()
	defer ih.commandMutex.Unlock()

	if _, exists := ih.guildCommands[guildID]; !exists {
		ih.guildCommands[guildID] = make(map[string]types.ApplicationCommand)
	}

	if _, exists := ih.guildCommands[guildID][command.Name]; exists {
		return fmt.Errorf("command %s already registered for guild %s", command.Name, guildID)
	}

	err := ih.bot.gateway.CreateGuildApplicationCommand(guildID, command)
	if err != nil {
		return err
	}

	ih.guildCommands[guildID][command.Name] = command
	return nil
}

func (ih *InteractionHandler) SyncCommands() error {
	for _, cmd := range ih.globalCommands {
		err := ih.bot.gateway.CreateGlobalApplicationCommand(cmd)
		if err != nil {
			return err
		}
	}

	for guildID, commands := range ih.guildCommands {
		for _, cmd := range commands {
			err := ih.bot.gateway.CreateGuildApplicationCommand(guildID, cmd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateButton(style int, label, customID string, options ...ButtonOption) types.ButtonComponent {
	button := types.ButtonComponent{
		Type:     types.ComponentTypeButton,
		Style:    style,
		Label:    label,
		CustomID: customID,
	}

	for _, opt := range options {
		opt(&button)
	}

	return button
}

type ButtonOption func(*types.ButtonComponent)

func WithEmoji(emoji types.Emoji) ButtonOption {
	return func(button *types.ButtonComponent) {
		button.Emoji = &emoji
	}
}

func WithDisabled(disabled bool) ButtonOption {
	return func(button *types.ButtonComponent) {
		button.Disabled = disabled
	}
}

func CreateSelectMenu(customID, placeholder string, options []types.SelectOption, minValues, maxValues int) types.SelectMenuComponent {
	return types.SelectMenuComponent{
		Type:        types.ComponentTypeSelectMenu,
		CustomID:    customID,
		Options:     options,
		Placeholder: placeholder,
		MinValues:   &minValues,
		MaxValues:   maxValues,
	}
}

func CreateSelectOption(label, value, description string, defaultOption bool) types.SelectOption {
	return types.SelectOption{
		Label:       label,
		Value:       value,
		Description: description,
		Default:     defaultOption,
	}
}
