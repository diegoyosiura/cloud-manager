package messaging

import (
	"context"
	"fmt"
	"github.com/diegoyosiura/cloud-manager/pkg/authentication"
	"net/smtp"
	"sync"
	"time"
)

const MaxOCIMessages = 10

type OciManager struct {
	Auth        *authentication.OCIAuth // OCI authentication details.
	sendContext context.Context
	Client      smtp.Auth

	Messages   []Message
	MessagesMT *sync.RWMutex
}

func (o *OciManager) setup() (bool, error) {
	o.Client = smtp.PlainAuth("", string(o.Auth.EmailUser), string(o.Auth.EmailPassword), o.Auth.EmailHost)
	return true, nil
}

func (o *OciManager) AddMessage(m Message) {
	o.MessagesMT.Lock()
	defer o.MessagesMT.Unlock()
	o.Messages = append(o.Messages, m)
}

func (o *OciManager) AddMessages(m []Message) {
	o.MessagesMT.Lock()
	defer o.MessagesMT.Unlock()
	o.Messages = append(o.Messages, m...)
}
func (o *OciManager) CancelSend() (bool, error) {
	ready, err := o.setup()

	if !ready {
		return false, err
	}
	panic("implement me")
}

func (o *OciManager) Send() (chan Message, bool, error) {
	ready, err := o.setup()

	if !ready {
		return nil, false, err
	}

	ch := o.sendMessage()

	return ch, true, nil
}

func (o *OciManager) SendStatus() (float64, error) {
	ready, err := o.setup()

	if !ready {
		return 0.0, err
	}

	sent := 0.0
	o.MessagesMT.Lock()
	defer o.MessagesMT.Unlock()
	for _, msg := range o.Messages {
		if msg.Status == Sent {
			sent++
		}
	}

	return sent / float64(len(o.Messages)), nil
}

func (o *OciManager) sendMessage() chan Message {
	ch := make(chan Message, MaxOCIMessages)

	go func() {
		defer close(ch)
		wg := &sync.WaitGroup{}

		tm := len(o.Messages)
		for i := 0; i < tm; i++ {
			o.MessagesMT.Lock()
			m := o.Messages[i]
			o.MessagesMT.Unlock()
			m.Status = Queued
			ch <- m
			wg.Add(1)
			go o.send(ch, m, wg)
		}

		wg.Wait()
	}()
	return ch
}

func (o *OciManager) send(ch chan Message, m Message, wg *sync.WaitGroup) {
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

	err = smtp.SendMail(fmt.Sprintf(`%s:%s`, o.Auth.EmailHost, o.Auth.EmailPort), o.Client, m.From.Address, list, data)

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
