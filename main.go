package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/genai"
)

// Your Google API key
const apiKey = "<GEMINI API KEY>" // Replace with your actual Google API key

func main() {
	tgToken := "<Telegram BOT KEY>" // Replace with your actual Telegram bot token

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Telegram bot
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Fatalf("fatal error creating bot: %v", err)
	}
	bot.Debug = false // Set to true for more detailed logs from Telegram API
	log.Printf("Authorized on account: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil { // Ignore any non-message updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "clr":
				userInput := strings.TrimPrefix(update.Message.Text, "/clr ")
				if userInput == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide content after /clr, e.g., /clr Today I learned about machine learning concepts.")
					_, err := bot.Send(msg)
					if err != nil {
						log.Printf("Error sending message: %v", err)
					}
					continue
				}

				// Prepare the standard classroom learning form template
				template := `Fill the form: Classroom Learning

1) Concept Learned
2) New Techniques Learned
3) Related Project/Practice work learned
4) New software/tool/experiment/equipment/machine learned
5) Application of concepts
6) Case studies/examples understood

Each topic each form. The form should be in points simple and straight
Subject name and topics covered: `

				// Combine template with user input
				fullPrompt := template + userInput

				// Call the GenerateContent method with the combined prompt
				result, err := client.Models.GenerateContent(ctx, "gemini-1.5-flash", genai.Text(fullPrompt), nil)
				if err != nil {
					log.Printf("Error generating content: %v", err)
					errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I couldn't generate a response. Please try again later.")
					_, sendErr := bot.Send(errMsg)
					if sendErr != nil {
						log.Printf("Error sending error message: %v", sendErr)
					}
					continue
				}

				// Extract the text from the response and send it back to Telegram
				if len(result.Candidates) > 0 && result.Candidates[0].Content != nil && len(result.Candidates[0].Content.Parts) > 0 {
					var responseText strings.Builder
					for _, part := range result.Candidates[0].Content.Parts {
						if part != nil {
							// The genai.Part struct contains the text at the end
							// We need to extract it from the struct representation
							partStr := fmt.Sprintf("%v", part)
							
							// The text appears to be at the end of the struct after all the <nil> values
							// Look for the pattern: } [text content here]\n}
							if strings.Contains(partStr, "&{") {
								// Find the last occurrence of a meaningful text after the struct fields
								// The structure appears to be: &{<nil> false <nil> <nil> [] <nil> <nil> <nil> <nil> ACTUAL_TEXT}
								
								// Split by the struct opening and find the text content
								parts := strings.Split(partStr, " ")
								var textFound bool
								var textContent strings.Builder
								
								for i, p := range parts {
									// Skip the struct notation parts and look for actual content
									if !textFound && (p == "<nil>" || p == "false" || p == "[]" || p == "&{" || strings.HasPrefix(p, "&{")) {
										continue
									}
									// Once we find non-struct content, start collecting it
									if !strings.Contains(p, "<nil>") && !strings.Contains(p, "false") && p != "[]" && p != "" {
										textFound = true
										// Remove trailing } if it's the last part
										if i == len(parts)-1 {
											p = strings.TrimRight(p, "}")
										}
										if textContent.Len() > 0 {
											textContent.WriteString(" ")
										}
										textContent.WriteString(p)
									}
								}
								
								if textContent.Len() > 0 {
									responseText.WriteString(strings.TrimSpace(textContent.String()))
								}
							} else {
								// Fallback: if the structure is different, just use the string as-is
								responseText.WriteString(partStr)
							}
							
							log.Printf("Part type: %T, extracted text: %s", part, responseText.String())
						}
					}

					if responseText.Len() > 0 {
						response := strings.TrimSpace(responseText.String())
						
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Error sending message: %v", err)
						}
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Could not extract text from the AI response.")
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Error sending message: %v", err)
						}
					}

				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No response or content generated by the AI model.")
					_, err := bot.Send(msg)
					if err != nil {
						log.Printf("Error sending message: %v", err)
					}
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Use /clr followed by your prompt.")
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
			}
		} else {
			// Respond to non-command messages if needed, or just ignore them
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please use the /clr command followed by your learning content. For example: /clr Today I learned about Python programming and data structures.")
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}

// debugPrint remains the same, useful for local debugging
func debugPrint[T any](r *T) {
	response, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response))
}