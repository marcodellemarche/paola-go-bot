package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"paola-go-bot/chatgpt"
	"paola-go-bot/telegram/commands"
	"paola-go-bot/telegram/utils"
	"time"
)

func (t *Telegram) useAssistant(message string, chatID int64, threadID *string) (string, error) {
	var runID string

	if *threadID == "" {
		run, err := t.chatgptClient.CreateThreadAndRun(message)
		if err != nil {
			return "", fmt.Errorf("error creating thread and run: %s", err)
		}
		// log.Printf("Created thread and run: %+v\n", run)

		*threadID = run.ThreadID
		runID = run.ID
	} else {
		_, err := t.chatgptClient.CreateMessage(*threadID, message)
		if err != nil {
			return "", fmt.Errorf("error getting message: %s", err)
		}
		// log.Printf("Got message: %+v\n", message)

		run, err := t.chatgptClient.CreateRun(*threadID)
		if err != nil {
			return "", fmt.Errorf("error creating run: %s", err)
		}
		// log.Printf("Create run: %+v\n", run)
		runID = run.ID
	}

	maxTries := 10
	interval := 500 * time.Millisecond

	for i := 0; i < maxTries; i++ {
		time.Sleep(interval)

		run, err := t.chatgptClient.GetRun(*threadID, runID)
		if err != nil {
			return "", fmt.Errorf("error getting run: %s", err)
		}
		// log.Printf("Got run: %+v\n", run)

		if run.Status == "completed" {
			log.Printf("Conversation continued with the assistant")
			msgList, err := t.chatgptClient.ListMessages(*threadID)
			if err != nil {
				return "", fmt.Errorf("error listing messages: %s", err)
			}
			// log.Printf("List of messages: %+v\n", msgList)

			answer := "Sò na sega"
			for _, msg := range msgList.Data {
				if msg.Role == "assistant" {
					log.Printf("Message from assistant: %+v\n", msg)

					answer = msg.Content[0].Text.Value
					break
				}
			}

			return answer, nil
		}

		if run.Status == "requires_action" {
			if run.RequiredAction.SubmitToolOutputs.ToolCalls == nil {
				return "", fmt.Errorf("no tool calls in required action")
			}

			toolCall := run.RequiredAction.SubmitToolOutputs.ToolCalls[0]
			var answer string

			switch toolCall.Function.Name {
			case "set_birthday":
				log.Printf("Set birthday: %s", toolCall.Function.Arguments)
				answer, err = setBirthday(toolCall.Function.Arguments, chatID)
			case "list_birthdays":
				log.Printf("List birthdays")
				answer, err = listBirthdays(chatID)
			case "delete_birthday":
				log.Printf("Delete birthday: %s", toolCall.Function.Arguments)
				answer, err = deleteBirthday(toolCall.Function.Arguments, chatID)
			case "answer_question":
				log.Printf("Answer question: %s", toolCall.Function.Arguments)
				answer, err = askAI(toolCall.Function.Arguments, t.chatgptClient)
			case "random_insult":
				log.Printf("Random insult")
				answer = utils.RandomInsult()
			default:
				return "", fmt.Errorf("unexpected tool call: %s", toolCall.Function.Name)
			}

			if err != nil {
				log.Printf("Error asking AI: %s", err)

				err = t.chatgptClient.CancelRun(*threadID, runID)
				if err != nil {
					return "", fmt.Errorf("error asking AI: error cancelling run: %s", err)
				}

				return "", fmt.Errorf("error asking AI: %s", err)
			}

			err = t.chatgptClient.SubmitToolOutputs(*threadID, runID, toolCall.ID, fmt.Sprintf(`{"message":"%s"}`, answer))
			if err != nil {
				log.Printf("Error submitting tool outputs: %s", err)

				err = t.chatgptClient.CancelRun(*threadID, runID)
				if err != nil {
					return "", fmt.Errorf("error submitting tool outputs: error cancelling run: %s", err)
				}

				return "Impara a stà al mondo", nil
			}

			*threadID = ""

			return answer, nil
		}
	}

	return "", fmt.Errorf("run did not complete")
}

// Input looks like {"name":"Antonio Giaguaro","date":"2023-11-07","contact_id":"1234567890}
func setBirthday(arguments string, chatID int64) (string, error) {
	type BirthdayArguments struct {
		Name      string `json:"name"`
		Date      string `json:"date"`
		ContactID int64  `json:"contact_id"`
	}

	// Parse arguments into name and date
	var birthday BirthdayArguments
	err := json.Unmarshal([]byte(arguments), &birthday)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling arguments: %s", err)
	}

	dateFormats := []string{"2006-01-02", "01-02"}

	var birthdayDate time.Time
	for _, format := range dateFormats {
		birthdayDate, err = time.Parse(format, birthday.Date)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", fmt.Errorf("error parsing date: %s", err)
	}

	parsedDay := birthdayDate.Day()
	parsedMonth := birthdayDate.Month()

	return commands.SetBirthday(birthday.Name, birthday.ContactID, uint8(parsedDay), uint8(parsedMonth), chatID)
}

// Input looks like {"name":"Antonio Giaguaro"}
func deleteBirthday(arguments string, chatID int64) (string, error) {
	type BirthdayArguments struct {
		Name string `json:"name"`
	}

	// Parse arguments into name and date
	var birthday BirthdayArguments
	err := json.Unmarshal([]byte(arguments), &birthday)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling arguments: %s", err)
	}

	return commands.DeleteBirthday(birthday.Name, chatID)
}

func listBirthdays(chatID int64) (string, error) {
	return commands.ListBirthdays(chatID)
}

// Input looks like {"question":"Quanti giri fa una boccia?"}
func askAI(arguments string, chatgptClient *chatgpt.ChatGPT) (string, error) {
	type QuestionArguments struct {
		Question string `json:"question"`
	}

	// Parse arguments into name and date
	var question QuestionArguments
	err := json.Unmarshal([]byte(arguments), &question)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling arguments: %s", err)
	}

	answer, err := chatgptClient.CreateCompletion(question.Question, nil)
	if err != nil {
		return "", fmt.Errorf("error creating completion: %s", err)
	}

	return answer, nil
}
