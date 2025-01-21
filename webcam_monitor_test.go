package main

import "testing"

type MockCommandRunner struct {
	Error   error
	Command string
	Args    []string
}

func (m *MockCommandRunner) Run(command string, args ...string) error {
	m.Command = command
	m.Args = args
	return m.Error
}

func TestSendMessage(t *testing.T) {
	mockRunner := &MockCommandRunner{}
	recipient := "test@example.com"
	message := "test message"

	err := sendMessage(mockRunner, recipient, message)
	if err != nil {
		t.Fatalf("Got an unexpected error: %v\n", err)
	}

	expectedCommand := "osascript"
	expectedArgs := []string{
		"-e",
		`
        tell application "Messages"
			set targetService to 1st service whose service type = iMessage
			set targetBuddy to buddy "test@example.com" of targetService
			send "test message" to targetBuddy
		end tell
        `,
	}
	if mockRunner.Command != expectedCommand || len(mockRunner.Args) != len(expectedArgs) {
		t.Errorf(
			"expected command %q with args %v, but got  command %q with args %v\n",
			expectedCommand, expectedArgs, mockRunner.Command, mockRunner.Args,
		)
	}
}

func TestParseLine(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected WebcamEvent
	}{
		{
			Input:    "random log",
			Expected: WebcamEvent{"", false},
		},
		{
			Input:    "",
			Expected: WebcamEvent{"", false},
		},
		{
			Input:    "AVCaptureSessionDidStartRunningNotification",
			Expected: WebcamEvent{"started", true},
		},
		{
			Input:    "AVCaptureSessionDidStopRunningNotification",
			Expected: WebcamEvent{"stopped", true},
		},
	}
	for _, testCase := range testCases {
		result := parseLine(testCase.Input)
		if result != testCase.Expected {
			t.Errorf("Expected  %v, got %v\n", testCase.Expected, result)
		}
	}
}
