package scripts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func getLatestTag(githubUsername, imageName string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", githubUsername, imageName)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch tags: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Response: %s\n", body)

	var tags []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &tags); err != nil {
		return "", fmt.Errorf("failed to parse tags: %v", err)
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}

	return tags[0].Name, nil
}

func updateDockerCompose(imageName, tag string) error {
	fmt.Printf("Pulling image with tag: %s\n", tag)
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", imageName, tag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull Docker image: %v", err)
	}

	fmt.Println("Restarting Docker Compose...")
	cmd = exec.Command("docker", "compose", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Docker Compose: %v", err)
	}

	cmd = exec.Command("docker", "compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Docker Compose: %v", err)
	}

	return nil
}

func CheckTag() {
	githubUsername := os.Getenv("GITHUB_USERNAME")
	imageName := os.Getenv("IMAGE_NAME")

	checkIntervalStr := os.Getenv("CHECK_INTERVAL")
	checkInterval := 10 * time.Second
	if checkIntervalStr != "" {
		if interval, err := time.ParseDuration(checkIntervalStr); err == nil {
			checkInterval = interval
		}
	}

	currentTag := ""

	for {
		fmt.Println("Checking for the latest tag...")
		latestTag, err := getLatestTag(githubUsername, imageName)
		fmt.Printf("Latest tag: %s\n", latestTag)
		fmt.Printf("Current tag: %s\n", currentTag)
		if err != nil {
			fmt.Printf("Error fetching latest tag: %v\n", err)
		} else if latestTag != currentTag {
			fmt.Printf("New tag found: %s (current: %s)\n", latestTag, currentTag)
			// if err := updateDockerCompose(imageName, latestTag); err != nil {
			// 	fmt.Printf("Error updating Docker Compose: %v\n", err)
			// } else {
			// 	currentTag = latestTag
			// 	fmt.Println("Update successful!")
			// }
		} else {
			fmt.Println("No new tag found.")
		}

		time.Sleep(checkInterval)
	}
}
