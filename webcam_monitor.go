package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	MessageTemplate = `
		tell application "Messages"
			set targetService to 1st service whose service type = iMessage
			set targetBuddy to buddy "%s" of targetService
			send "%s" to targetBuddy
		end tell
	`
	LogPredicate = `
        (eventMessage CONTAINS "AVCaptureSessionDidStartRunningNotification" || 
        eventMessage CONTAINS "AVCaptureSessionDidStopRunningNotification")
    `
)

func monitorWebcam(recipient string) {
	cmd := exec.Command("log", "stream", "--predicate", LogPredicate)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error with pipe stdout: %v\n", err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting cmd: %v\n", err)
	}
	scanner := bufio.NewScanner(stdout)
	skipFirstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		if skipFirstLine {
			skipFirstLine = false
			continue
		}
		if strings.Contains(line, "AVCaptureSessionDidStartRunningNotification") {
			fmt.Println("Webcam started...")
			err := sendMessage(recipient, "Webcam is on!!")
			if err != nil {
				fmt.Printf("Could not send message: %v\n", err)
			} else {
				fmt.Println("Message sent.")
			}
		} else if strings.Contains(line, "AVCaptureSessionDidStopRunningNotification") {
			fmt.Println("Webcam stopped...")
			err := sendMessage(recipient, "Webcam is off now!!")
			if err != nil {
				fmt.Printf("Could not send message: %v\n", err)
			} else {
				fmt.Println("Message sent.")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading log stream: %v\n", err)
	}
}

func sendMessage(recipient string, message string) error {
	script := fmt.Sprintf(MessageTemplate, recipient, message)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

func main() {
	recipient := os.Getenv("MESSAGE_RECIPIENT")
	if recipient == "" {
		fmt.Println("Missing MESSAGE_RECIPIENT in environment variables.")
		return
	}
	fmt.Println("Starting webcam daemon...")
	monitorWebcam(recipient)
}
