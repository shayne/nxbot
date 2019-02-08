package main

import (
	"log"

	"github.com/jacknx/nxbot/internal/app/nxbot"
	"github.com/jacknx/nxbot/pkg/nxapi"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	api, err := nxapi.NewAPI(conf.NxIPPort, conf.NxUser, conf.NxPass)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Nx API")
	cameras, err := api.GetCameras()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Retrieved cameras from Nx API")

	bot, err := nxbot.NewNxBot(&nxbot.Settings{
		Token:          conf.TgToken,
		HTTPIPPort:     conf.HTTPIPPort,
		UserWhitelist:  conf.TgUserWhitelist,
		GroupWhitelist: conf.TgGroupWhitelist,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created Nx Telegram Bot")

	for _, camera := range cameras {
		bot.AddCamera(camera.ID, camera.Name)
	}

	bot.OnCameraButtonPressed(func(id string, name string, m *tb.Message) {
		log.Printf("Camera snapshot requested: %s = %s\n", id, name)
		buf, err := api.GetSnapshot(id)
		if err != nil {
			log.Fatal(err)
		}
		bot.ReplyWithPhoto(m, buf)
	})

	bot.OnMotionEvent(func(cameraID string) {
		log.Printf("Motion event received: %s\n", cameraID)
		buf, err := api.GetSnapshot(cameraID)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range conf.TgMotionRecipients {
			log.Printf("Sending motion event to: %s\n", r)
			buf.Reset()
			bot.SendPhoto(&nxbot.Recipient{ID: r}, buf)
		}
	})

	log.Printf("Starting Nx Telegram Bot and motion-event HTTP server on port %s\n", conf.HTTPIPPort)
	bot.Start()
}
