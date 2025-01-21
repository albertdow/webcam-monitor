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

type WebcamEvent struct {
	State string
	Valid bool
}

type CommandRunner interface {
	Run(command string, args ...string) error
}

type DefaultCommandRunner struct{}

func (r *DefaultCommandRunner) Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

func sendMessage(runner CommandRunner, recipient string, message string) error {
	script := fmt.Sprintf(MessageTemplate, recipient, message)
	return runner.Run("osascript", "-e", script)
}

func parseLine(line string) WebcamEvent {
	if strings.Contains(line, "AVCaptureSessionDidStartRunningNotification") {
		return WebcamEvent{"started", true}
	} else if strings.Contains(line, "AVCaptureSessionDidStopRunningNotification") {
		return WebcamEvent{"stopped", true}
	}
	return WebcamEvent{"", false}
}

func monitorWebcam(recipient string) {
	cmd := exec.Command("log", "stream", "--predicate", LogPredicate)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error with pipe stdout: %v\n", err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting cmd: %v\n", err)
	}
	runner := &DefaultCommandRunner{}
	scanner := bufio.NewScanner(stdout)
	skipFirstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		if skipFirstLine {
			skipFirstLine = false
			continue
		}
		parsedLine := parseLine(line)
		if parsedLine.Valid {
			switch parsedLine.State {
			case "started":
				fmt.Println("Webcam started...")
				err := sendMessage(runner, recipient, "Webcam is on!!")
				if err != nil {
					fmt.Printf("Could not send message: %v\n", err)
				} else {
					fmt.Println("Message sent.")
				}
			case "stopped":
				fmt.Println("Webcam stopped...")
				err := sendMessage(runner, recipient, "Webcam is now off!!")
				if err != nil {
					fmt.Printf("Could not send message: %v\n", err)
				} else {
					fmt.Println("Message sent.")
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading log stream: %v\n", err)
	}
}

func main() {
	recipient := os.Getenv("MESSAGE_RECIPIENT")
	if recipient == "" {
		fmt.Println("Missing MESSAGE_RECIPIENT in environment variables.")
		return
	}
	fmt.Println("Starting webcam monitor...")
	monitorWebcam(recipient)
}
