package messaging

import (
	"bytes"
	"net/mail"
	"testing"
)

// Helper function to generate a sample message
// This function creates a base Message object that can be reused across tests.
// It ensures consistency and reduces duplication in test setup.
func generateSampleMessage() Message {
	return NewMessage(
		mail.Address{Name: "Test", Address: "from@email.com"},
		"Test Subject",                // Subject of the email
		"This is a test body.",        // Body content of the email
		"text/plain",                  // Content type for the body
		[]string{"to@example.com"},    // Primary recipient
		[]string{"cc@example.com"},    // CC recipient
		[]string{"bcc@example.com"},   // BCC recipient
		[]string{"reply@example.com"}, // Reply-To address
	)
}

// Test creating a new message with default values
// Verifies that the NewMessage function initializes a Message object correctly.
func TestNewMessage(t *testing.T) {
	msg := generateSampleMessage()

	// Check if the subject is correctly set
	if msg.Subject != "Test Subject" {
		t.Errorf("expected subject 'Test Subject', got '%s'", msg.Subject)
	}

	// Ensure the body content type defaults to "text/plain"
	if msg.BodyContentType != "text/plain" {
		t.Errorf("expected content type 'text/plain', got '%s'", msg.BodyContentType)
	}

	// Validate that the "To" list contains the expected email address
	if len(msg.MailTo) != 1 || msg.MailTo[0] != "to@example.com" {
		t.Errorf("expected 'to@example.com' in MailTo, got '%v'", msg.MailTo)
	}
}

// Test adding attachments
// Verifies that files can be added as attachments to the message.
func TestAttach(t *testing.T) {
	msg := generateSampleMessage()

	// Attach a sample file (ensure the file exists at the specified path)
	err := msg.Attach("testdata/sample.txt")
	if err != nil {
		t.Errorf("unexpected error attaching file: %v", err)
	}

	// Validate that the attachment was added
	if len(msg.Attachments) != 1 {
		t.Errorf("expected 1 attachment, got %d", len(msg.Attachments))
	}

	// Check that the attachment has the correct filename
	attachment, ok := msg.Attachments["sample.txt"]
	if !ok {
		t.Error("attachment 'sample.txt' not found")
	}

	// Ensure the filename matches the expected value
	if attachment.Filename != "sample.txt" {
		t.Errorf("expected filename 'sample.txt', got '%s'", attachment.Filename)
	}
}

// Test attaching inline content
// Verifies that files can be added as inline attachments.
func TestInline(t *testing.T) {
	msg := generateSampleMessage()

	// Attach a sample file as inline (ensure the file exists at the specified path)
	err := msg.Inline("testdata/image.png")
	if err != nil {
		t.Errorf("unexpected error attaching inline file: %v", err)
	}

	// Validate that the inline attachment was added
	if len(msg.Attachments) != 1 {
		t.Errorf("expected 1 inline attachment, got %d", len(msg.Attachments))
	}

	// Check that the inline attachment exists and is marked as inline
	attachment, ok := msg.Attachments["image.png"]
	if !ok {
		t.Error("inline attachment 'image.png' not found")
	}

	// Ensure the inline flag is correctly set
	if !attachment.Inline {
		t.Error("attachment should be marked as inline")
	}
}

// Test adding custom headers
// Verifies that custom headers can be added to the message.
func TestAddHeader(t *testing.T) {
	msg := generateSampleMessage()

	// Add a custom header
	header := msg.AddHeader("X-Custom-Header", "CustomValue")

	// Validate that the header has the correct key and value
	if header.Key != "X-Custom-Header" || header.Value != "CustomValue" {
		t.Errorf("expected header 'X-Custom-Header: CustomValue', got '%s: %s'", header.Key, header.Value)
	}

	// Ensure the header is added to the message
	if len(msg.Headers) != 1 {
		t.Errorf("expected 1 header, got %d", len(msg.Headers))
	}
}

// Test generating recipients list
// Verifies that the Tolist method compiles all recipients (To, CC, BCC).
func TestTolist(t *testing.T) {
	msg := generateSampleMessage()

	// Generate the list of recipients
	recipients, err := msg.Tolist()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Define the expected list of recipients
	expectedRecipients := []string{"to@example.com", "cc@example.com", "bcc@example.com"}

	// Check that the number of recipients matches the expected count
	if len(recipients) != len(expectedRecipients) {
		t.Errorf("expected %d recipients, got %d", len(expectedRecipients), len(recipients))
	}

	// Ensure all expected recipients are present in the list
	for _, recipient := range expectedRecipients {
		found := false
		for _, r := range recipients {
			if r == recipient {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected recipient '%s' not found", recipient)
		}
	}
}

// Test encoding the message to bytes
// Verifies that the Bytes method constructs the message with headers, body, and attachments.
func TestBytes(t *testing.T) {
	msg := generateSampleMessage()

	// Add a custom header
	msg.AddHeader("X-Test-Header", "TestValue")

	// Encode the message into bytes
	data, err := msg.Bytes()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check for the presence of essential headers and their values
	if !bytes.Contains(data, []byte("Subject: =?UTF-8?B?VGVzdCBTdWJqZWN0?=")) {
		t.Error("missing or invalid subject header")
	}
	if !bytes.Contains(data, []byte("To: to@example.com")) {
		t.Error("missing 'To' header")
	}
	if !bytes.Contains(data, []byte("Cc: cc@example.com")) {
		t.Error("missing 'Cc' header")
	}
	if !bytes.Contains(data, []byte("Reply-To: reply@example.com")) {
		t.Error("missing 'Reply-To' header")
	}
	if !bytes.Contains(data, []byte("X-Test-Header: TestValue")) {
		t.Error("missing custom header")
	}
}
