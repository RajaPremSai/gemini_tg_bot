# Gemini Telegram Classroom Learning Bot

This project is a Telegram bot that leverages Google's Gemini API to help users fill out a classroom learning form in a structured, point-wise manner. Users interact with the bot using the `/clr` command, and the bot responds with a formatted summary based on the provided input.

## Features

- Accepts user input via the `/clr` command in Telegram.
- Uses Gemini (Google AI) to generate structured classroom learning summaries.
- Responds with a point-wise form covering concepts learned, techniques, projects, tools, applications, and examples.
- Handles errors and provides user-friendly feedback.

## Requirements

- Go 1.20 or newer
- A Telegram Bot Token ([How to get one](https://core.telegram.org/bots#how-do-i-create-a-bot))
- A Google Gemini API Key ([Get one here](https://aistudio.google.com/app/apikey))

## Setup

1. **Clone the repository:**

   ```sh
   git clone <your-repo-url>
   cd Gemini-TG-BOT
   ```

2. **Configure your credentials:**

   - Edit `config.yaml` and set your `tgToken` (Telegram Bot Token) and `geminiToken` (Google Gemini API Key):

     ```yaml
     tgToken: "<YOUR_TELEGRAM_BOT_TOKEN>"
     geminiToken: "<YOUR_GEMINI_API_KEY>"
     preamble: "E5: "
     ```

3. **Install dependencies:**

   ```sh
   go mod tidy
   ```

4. **Run the bot:**
   ```sh
   go run main.go
   ```

## Usage

- In your Telegram app, start a chat with your bot.
- Use the `/clr` command followed by your learning content.  
  **Example:**

  ```
  /clr Today I learned about machine learning concepts and implemented a small project using Python.
  ```

- The bot will reply with a structured classroom learning form based on your input.

## Example Output

```
Fill the form: Classroom Learning

1) Concept Learned
2) New Techniques Learned
3) Related Project/Practice work learned
4) New software/tool/experiment/equipment/machine learned
5) Application of concepts
6) Case studies/examples understood

Each topic each form. The form should be in points simple and straight
Subject name and topics covered: Today I learned about machine learning concepts and implemented a small project using Python.
```

## Libraries Used

- [github.com/go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) — Telegram Bot API for Go
- [google.golang.org/genai](https://pkg.go.dev/google.golang.org/genai) — Google Gemini API Go SDK
- [github.com/spf13/viper](https://github.com/spf13/viper) — Go configuration with YAML support

## Things I Learned & Experience

- Integrating Google Gemini API with Go and handling its evolving SDK.
- Building a Telegram bot that interacts with users and external APIs.
- Parsing and structuring AI-generated responses for user-friendly output.
- Managing configuration securely using YAML and Viper.
- Handling Go module dependencies and troubleshooting version mismatches.
- Error handling and logging for robust bot operation.

## Notes

- The bot only responds to the `/clr` command.
- If you do not provide content after `/clr`, the bot will prompt you for more information.
- Make sure your API keys are kept private and not shared publicly.

## License

MIT License

---

**Contributions are welcome!**
