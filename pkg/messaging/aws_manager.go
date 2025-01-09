package messaging

import (
	"fmt"
	"github.com/diegoyosiura/cloud-manager/pkg/authentication"
	"net/smtp"
	"sync"
	"time"
)

type AWSManager struct {
	Auth   *authentication.AWSAuth // AWS authentication details.
	Client smtp.Auth

	Messages   []Message
	MessagesMT *sync.RWMutex
}

func (a *AWSManager) setup() (bool, error) {
	a.Client = smtp.PlainAuth("", string(a.Auth.EmailUser), string(a.Auth.EmailPassword), a.Auth.EmailHost)
	return true, nil
}

func (a *AWSManager) AddMessage(m Message) {
	a.MessagesMT.Lock()
	defer a.MessagesMT.Unlock()
	a.Messages = append(a.Messages, m)
}

func (a *AWSManager) AddMessages(m []Message) {
	a.MessagesMT.Lock()
	defer a.MessagesMT.Unlock()
	a.Messages = append(a.Messages, m...)
}
func (a *AWSManager) CancelSend() (bool, error) {
	ready, err := a.setup()

	if !ready {
		return false, err
	}
	panic("implement me")
}

func (a *AWSManager) Send() (chan Message, bool, error) {
	ready, err := a.setup()

	if !ready {
		return nil, false, err
	}

	ch := a.sendMessage()

	return ch, true, nil
}

func (a *AWSManager) SendStatus() (float64, error) {
	ready, err := a.setup()

	if !ready {
		return 0.0, err
	}

	sent := 0.0
	a.MessagesMT.Lock()
	defer a.MessagesMT.Unlock()
	for _, msg := range a.Messages {
		if msg.Status == Sent {
			sent++
		}
	}

	return sent / float64(len(a.Messages)), nil
}

func (a *AWSManager) sendMessage() chan Message {
	ch := make(chan Message, MaxOCIMessages)

	go func() {
		defer close(ch)
		wg := &sync.WaitGroup{}

		tm := len(a.Messages)
		for i := 0; i < tm; i++ {
			a.MessagesMT.Lock()
			m := a.Messages[i]
			a.MessagesMT.Unlock()
			m.Status = Queued
			ch <- m
			wg.Add(1)
			go a.send(ch, m, wg)
		}

		wg.Wait()
	}()
	return ch
}

func (a *AWSManager) send(ch chan Message, m Message, wg *sync.WaitGroup) {
	defer wg.Done()
	m.Status = Sending
	ch <- m

	list, err := m.Tolist()
	if err != nil {
		m.Status = SendError
		m.DateStatus = time.Now()
		m.Error = err
		ch <- m
		return
	}

	data, err := m.Bytes()
	if err != nil {
		m.Status = SendError
		m.DateStatus = time.Now()
		m.Error = err
		ch <- m
		return
	}

	err = smtp.SendMail(fmt.Sprintf(`%s:%s`, a.Auth.EmailHost, a.Auth.EmailPort), a.Client, m.From.Address, list, data)

	if err != nil {
		m.Status = SendError
		m.DateStatus = time.Now()
		m.Error = err
		ch <- m
		return
	}

	m.Status = Sent
	m.DateStatus = time.Now()
	ch <- m
}
