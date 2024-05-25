package telegram

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) aiHandleVoice(message *tgbotapi.Message) status.CommandResponse {
	log.Printf("Generic vocal command - asking to PaolaGPT")

	if !t.chatgptClient.RateLimiter.Allow(message.Chat.ID) {
		log.Println("Rate limit exceeded")

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	fileID := message.Voice.FileID
	file, err := t.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		log.Printf("Error getting file: %s", err)

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	// Download file
	fileURL, err := t.bot.GetFileDirectURL(file.FileID)
	if err != nil {
		log.Printf("Error getting file URL: %s", err)

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	err = downloadFile(fileURL, file.FilePath)
	if err != nil {
		log.Printf("Error downloading file: %s", err)

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	log.Printf("Received voice message file: %+v\n", file)

	// By default, the Whisper API only supports files that are less than 25 MB.
	// If you have an audio file that is longer than that, you will need to break it up into
	// chunks of 25 MB's or less or used a compressed audio format. To get the best performance,
	// we suggest that you avoid breaking the audio up mid-sentence as this may cause some context to be lost.

	transcript, err := t.chatgptClient.Transcript(file.FilePath)
	if err != nil {
		log.Printf("Error transcribing voice: %s", err)

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	err = removeFile(file.FilePath)
	if err != nil {
		log.Printf("Error removing file: %s", err)

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, transcript)
	reply.ReplyToMessageID = message.MessageID
	t.SendMessage(&reply, nil)

	var threadID string
	if userStatus, exists := status.Get(message.Chat.ID); exists {
		threadID = userStatus.ThreadID
	}

	answer, err := t.useAssistant(transcript, message.Chat.ID, &threadID)
	if err != nil {
		log.Printf("Error using the AI assistant: %s", err)

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if threadID != "" {
		status.SetThread(message.Chat.ID, threadID)
	} else {
		status.ResetThread(message.Chat.ID)
	}

	status.ResetNext(message.Chat.ID)
	reply = tgbotapi.NewMessage(message.Chat.ID, answer)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}

func downloadFile(url string, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %s", err)
	}
	defer response.Body.Close()

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %s", err)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("error copying file: %s", err)
	}

	err = convertOggToMp3(path, strings.Replace(path, ".oga", ".mp3", 1))
	if err != nil {
		return fmt.Errorf("error converting file: %s", err)
	}

	return nil
}

func convertOggToMp3(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-acodec", "libmp3lame", outputPath)
	err := cmd.Run()
	return err
}

func removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to remove file: %s", err)
	}
	return nil
}
