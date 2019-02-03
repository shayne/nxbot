package nxbot

import (
	"fmt"
	"io"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Settings are used to configure access and security to your telegram bot
type Settings struct {
	Token          string
	UserWhitelist  []int
	GroupWhitelist []int64
	HTTPIPPort     string
}

// Recipient is simple struct conforming to the tb.Recipient interface
type Recipient struct {
	ID string
}

// Recipient simply returns the id string of the struct
func (r *Recipient) Recipient() string {
	return r.ID
}

type buttonPair struct {
	id  string
	btn tb.ReplyButton
}

// NxBot wraps the telebot library, making it an Nx specific bot
type NxBot struct {
	settings           *Settings
	bot                *tb.Bot
	cameraBtns         []buttonPair
	cameraBtnHandler   func(string, string, *tb.Message)
	motionEventHandler func(string)
}

// NewNxBot creates a new Nx bot given a bot token
func NewNxBot(settings *Settings) (*NxBot, error) {
	if settings.Token == "" {
		return nil, fmt.Errorf("NewNxBot failed: You must provide a Token")
	}
	if settings.HTTPIPPort == "" {
		return nil, fmt.Errorf("NewNxBot failed: You must provide an HTTP IP:PORT for motion events")
	}
	if settings.UserWhitelist == nil {
		log.Println("WARNING: UserWhitelist not set in settings, your bot is accessible to anyone!")
		settings.UserWhitelist = make([]int, 0)
	}
	if settings.GroupWhitelist == nil {
		settings.GroupWhitelist = make([]int64, 0)
	}
	bot, err := tb.NewBot(tb.Settings{
		Token:  settings.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("NewNxBot failed: %v", err)
	}
	return &NxBot{
		settings:   settings,
		bot:        bot,
		cameraBtns: make([]buttonPair, 0),
	}, nil
}

// AddCamera adds the given camera id, name to the list of cameras
// presented to retrieve snapshots
func (b *NxBot) AddCamera(id, name string) {
	cameraBtn := tb.ReplyButton{Text: name}
	b.cameraBtns = append(b.cameraBtns, buttonPair{id: id, btn: cameraBtn})
}

// OnCameraButtonPressed sets the camera button handler
func (b *NxBot) OnCameraButtonPressed(handler func(string, string, *tb.Message)) {
	b.cameraBtnHandler = handler
}

// SendPhoto sends a photo, from buf, to the given recipient
func (b *NxBot) SendPhoto(r tb.Recipient, buf io.Reader) {
	photo := &tb.Photo{File: tb.FromReader(buf)}
	b.bot.Send(r, photo)
}

// ReplyWithPhoto replies to the given message with the photo from buf
func (b *NxBot) ReplyWithPhoto(m *tb.Message, buf io.Reader) {
	b.SendPhoto(recipient(m), buf)
}

func (b *NxBot) isWhitelisted(m *tb.Message) bool {
	if m.Private() {
		u := m.Sender

		if len(b.settings.UserWhitelist) == 0 {
			return true
		}

		for _, id := range b.settings.UserWhitelist {
			if id == u.ID {
				return true
			}
		}
	}

	if m.FromGroup() {
		c := m.Chat
		for _, id := range b.settings.GroupWhitelist {
			if id == c.ID {
				return true
			}
		}
	}

	return false
}

// Start finishes setting up the Telegram bot then starts and waits indefinitely
func (b *NxBot) Start() {
	for _, btn := range b.cameraBtns {
		id := btn.id
		name := btn.btn.Text
		b.bot.Handle(&btn.btn, func(m *tb.Message) {
			if !b.isWhitelisted(m) {
				return
			}
			if b.cameraBtnHandler != nil {
				b.cameraBtnHandler(id, name, m)
			}
		})
	}

	nbtns := len(b.cameraBtns)
	replyKeys := make([][]tb.ReplyButton, nbtns)
	for i := 0; i < nbtns; i++ {
		replyKeys[i] = []tb.ReplyButton{b.cameraBtns[i].btn}
	}

	// Command: /start <PAYLOAD>
	b.bot.Handle("/start", func(m *tb.Message) {
		if !b.isWhitelisted(m) {
			return
		}

		b.bot.Send(recipient(m), "Hello!", &tb.ReplyMarkup{
			ReplyKeyboard: replyKeys,
		})
	})

	b.startMotionInBackground()
	b.bot.Start()
}

func recipient(m *tb.Message) tb.Recipient {
	var r tb.Recipient = m.Sender
	if m.FromGroup() {
		r = m.Chat
	}
	return r
}
