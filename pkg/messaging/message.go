package messaging

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Global regex for sanitizing filenames (compiled once for reuse)
var validFilenameRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// Buffer pool for optimized memory allocation when creating email content
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Message represents an email with metadata, recipients, and content.
type Message struct {
	ID              string                 // Message Identifier
	Subject         string                 // Email subject
	Body            string                 // Email body content
	Error           error                  // Error content
	Status          MessageStatus          // Current status of the email (e.g., NotSent, Sent)
	From            mail.Address           // Sender's email address
	MailTo          []string               // Primary recipients
	CC              []string               // Carbon copy recipients
	BCC             []string               // Blind carbon copy recipients
	Reply           []string               // Reply-To addresses
	BodyContentType string                 // MIME type of the body content (e.g., text/plain, text/html)
	Headers         []Header               // Additional custom headers
	Attachments     map[string]*Attachment // Attachments associated with the email
	DateReceived    time.Time              // Timestamp when the email was created
	DateStatus      time.Time              // Timestamp when the status was last updated
}

// NewMessage initializes a new Message object with default values if not provided.
func NewMessage(from mail.Address, subject, body, bodyContentType string, mailTo, cc, bcc, reply []string) Message {
	// Default to "text/html" if no bodyContentType is provided
	if bodyContentType == "" {
		bodyContentType = "text/html"
	}

	return Message{
		From:            from,
		Subject:         subject,
		Body:            body,
		Status:          NotSent,
		MailTo:          mailTo,
		CC:              cc,
		BCC:             bcc,
		Reply:           reply,
		BodyContentType: bodyContentType,
		Attachments:     make(map[string]*Attachment),
		DateReceived:    time.Now(),
		DateStatus:      time.Now(),
	}
}

// attach adds a file to the message's attachments, optionally setting it as inline content.
func (m *Message) attach(file string, inline bool) error {
	// Read the file contents
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to attach file '%s': %w", file, err)
	}

	// Sanitize the filename to prevent malicious input
	filename := sanitizeFilename(filepath.Base(file))

	// Store the attachment
	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}

	return nil
}

// sanitizeFilename replaces invalid characters in a filename with underscores.
func sanitizeFilename(filename string) string {
	return validFilenameRegex.ReplaceAllString(filename, "_")
}

// AddTo appends a recipient to the "To" list.
func (m *Message) AddTo(address mail.Address) []string {
	m.MailTo = append(m.MailTo, address.String())
	return m.MailTo
}

// AddCc appends a recipient to the "CC" list.
func (m *Message) AddCc(address mail.Address) []string {
	m.CC = append(m.CC, address.String())
	return m.CC
}

// AddBcc appends a recipient to the "BCC" list.
func (m *Message) AddBcc(address mail.Address) []string {
	m.BCC = append(m.BCC, address.String())
	return m.BCC
}

// AttachBuffer adds an attachment to the message directly from a buffer.
// It sanitizes the filename and sets the inline flag as specified.
func (m *Message) AttachBuffer(filename string, buf []byte, inline bool) error {
	// Check for empty buffer
	if len(buf) == 0 {
		return fmt.Errorf("buffer for attachment '%s' is empty", filename)
	}

	// Store the attachment
	m.Attachments[sanitizeFilename(filename)] = &Attachment{
		Filename: sanitizeFilename(filename),
		Data:     buf,
		Inline:   inline,
	}
	return nil
}

// Attach adds a file as a regular attachment (not inline).
func (m *Message) Attach(file string) error {
	return m.attach(file, false)
}

// Inline adds a file as an inline attachment.
func (m *Message) Inline(file string) error {
	return m.attach(file, true)
}

// AddHeader appends a custom header to the message.
func (m *Message) AddHeader(key, value string) Header {
	header := Header{Key: key, Value: value}
	m.Headers = append(m.Headers, header)
	return header
}

// Tolist compiles and validates all recipients from "To", "CC", and "BCC" lists.
func (m *Message) Tolist() ([]string, error) {
	// Combine all recipient lists
	allRecipients := append([]string{}, m.MailTo...)
	allRecipients = append(allRecipients, m.CC...)
	allRecipients = append(allRecipients, m.BCC...)

	// Validate and parse email addresses
	parsedAddresses := []string{}
	for _, recipient := range allRecipients {
		address, err := mail.ParseAddress(recipient)
		if err != nil {
			return nil, fmt.Errorf("invalid address '%s': %w", recipient, err)
		}
		parsedAddresses = append(parsedAddresses, address.Address)
	}

	return parsedAddresses, nil
}

// Bytes constructs the message into a byte slice suitable for sending via SMTP.
func (m *Message) Bytes() ([]byte, error) {
	// Get a buffer from the pool
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()

	// Validate the "From" address
	if _, err := mail.ParseAddress(m.From.Address); err != nil {
		return nil, fmt.Errorf("invalid 'From' address: %w", err)
	}

	// Add "From" and "Date" headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", m.From.String()))
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))

	// Add "To" and "CC" headers
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.MailTo, ", ")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(m.CC, ", ")))
	}

	// Encode and add the "Subject" header
	if !isUTF8(m.Subject) {
		return nil, fmt.Errorf("subject contains non-UTF-8 characters")
	}
	encodedSubject := base64.StdEncoding.EncodeToString([]byte(m.Subject))
	buf.WriteString(fmt.Sprintf("Subject: =?UTF-8?B?%s?=\r\n", encodedSubject))

	// Add "Reply-To" header if applicable
	if len(m.Reply) > 0 {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", strings.Join(m.Reply, ", ")))
	}

	// Add MIME version and custom headers
	buf.WriteString("MIME-Version: 1.0\r\n")
	for _, header := range m.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", header.Key, header.Value))
	}

	// Handle body and attachments
	if len(m.Attachments) > 0 {
		// Add multipart boundary for attachments
		boundary := "f46d043c813270fc6b04c2d223da"
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))

		// Add body content
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.BodyContentType))
		buf.WriteString(m.Body + "\r\n")

		// Add attachments
		for _, att := range m.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			mimeType := mime.TypeByExtension(filepath.Ext(att.Filename))
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
			buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimeType))
			buf.WriteString(fmt.Sprintf("Content-Disposition: %s; filename=\"%s\"\r\n", "attachment", att.Filename))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")

			// Encode and add attachment content
			encoded := make([]byte, base64.StdEncoding.EncodedLen(len(att.Data)))
			base64.StdEncoding.Encode(encoded, att.Data)
			buf.Write(encoded)
			buf.WriteString("\r\n")
		}

		// Close the multipart boundary
		buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// Add plain body content
		buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.BodyContentType))
		buf.WriteString(m.Body + "\r\n")
	}

	return buf.Bytes(), nil
}

// Send transmits the email message using the specified SMTP server.
func Send(addr string, auth smtp.Auth, m *Message) error {
	data, err := m.Bytes()
	if err != nil {
		return err
	}
	recipients, err := m.Tolist()
	if err != nil {
		return err
	}
	return smtp.SendMail(addr, auth, m.From.Address, recipients, data)
}

// isUTF8 checks if the given string contains only valid UTF-8 characters.
func isUTF8(s string) bool {
	for _, r := range s {
		if r == '�' {
			return false
		}
	}
	return true
}
