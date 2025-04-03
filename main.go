package main

import (
	"discord-go/api/types"
	"discord-go/bot"
	"fmt"
	"log"
)

func main() {
	b := bot.NewBot("MTM1NTYxMzc2NzcyNzA1OTIzNQ.Gz2EZC.WVlozEeuSi89dOOGyCgIhQnPqZLkgDfxwC4j2o")

	// Message command
	b.AddMessageHandler("!ping", func(ctx *bot.MessageContext) {
		fmt.Printf("Received ping command from %s\n", ctx.Message.Author.Username)
	})

	// Slash commands
	b.AddSlashCommand("ping", "Check if the bot is alive", func(ctx *bot.CommandContext) {
		b.RespondToInteraction(ctx, "Pong! -100ms using discord-go")
	})

	b.AddSlashCommand("advanced", "Test advanced interactions", func(ctx *bot.CommandContext) {
		components := []types.MessageComponent{
			b.NewActionRow(
				b.NewButton("Click me!", "test_button", types.ButtonStylePrimary),
				b.NewButton("Open Modal", "open_modal", types.ButtonStyleSecondary),
			),
			b.NewActionRow(
				b.NewSelectMenu("test_select", "Choose an option", []types.SelectOption{
					{Label: "Option 1", Value: "option_1"},
					{Label: "Option 2", Value: "option_2"},
					{Label: "Option 3", Value: "option_3"},
				}),
			),
		}

		b.RespondWithComponents(ctx, "Here are some interactive components:", components)
	})

	// Button handlers
	b.AddComponentHandler("test_button", func(ctx *bot.ComponentContext) {
		response := types.InteractionResponse{
			Type: types.InteractionResponseTypeUpdateMessage,
			Data: &types.InteractionCallbackData{
				Content: "You clicked the button!",
			},
		}
		err := b.SendInteractionResponse(ctx.Interaction.ID, ctx.Interaction.Token, response)
		if err != nil {
			log.Printf("Failed to handle button click: %v", err)
		}
	})

	b.AddComponentHandler("open_modal", func(ctx *bot.ComponentContext) {
		modal := b.NewModal("test_modal", "Test Modal",
			b.NewTextInputRow("modal_input_1", "Enter some text", types.TextInputStyleShort, true, 100, "Type something here..."),
			b.NewTextInputRow("modal_input_2", "Longer text", types.TextInputStyleParagraph, false, 500, "Type a longer message here..."),
		)
		b.RespondWithModal(ctx, modal)
	})

	// Select menu handler
	b.AddComponentHandler("test_select", func(ctx *bot.ComponentContext) {
		selected := ctx.Values[0]
		response := types.InteractionResponse{
			Type: types.InteractionResponseTypeUpdateMessage,
			Data: &types.InteractionCallbackData{
				Content: fmt.Sprintf("You selected: %s - brought to you by discord-go! by nyrilol!", selected),
			},
		}
		err := b.SendInteractionResponse(ctx.Interaction.ID, ctx.Interaction.Token, response)
		if err != nil {
			log.Printf("Failed to handle select menu: %v", err)
		}
	})

	b.SetModalHandler(func(ctx *bot.ModalContext) {
		log.Println(ctx)

		input1 := ctx.Inputs["modal_input_1"]
		input2 := ctx.Inputs["modal_input_2"]
		log.Printf("Modal submitted - Input 1: %s, Input 2: %s", input1, input2)
	})

	b.Start()
}
