package plugin

import (
	"context"
	"fmt"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	"time"
)

type PlayerSettings struct {
	PlayerId   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

func setupVolumeActions(client *streamdeck.Client) {

	// Volume Up
	//
	volumeUpAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.volumeup")
	volumeUpAction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		volume, err := squeezebox.ChangePlayerVolume(player, +10)
		if err != nil {
			client.ShowAlert(ctx)
			logError(client, "volumeup", err)
			return err
		}

		go displayTextAsTitleForTwoSeconds(ctx, client, fmt.Sprintf("%d", volume))

		return nil
	})

	volumeUpAction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	volumeUpAction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)

	// Volume Down
	//
	volumeDownAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.volumedown")
	volumeDownAction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		volume, err := squeezebox.ChangePlayerVolume(player, -10)
		if err != nil {
			client.ShowAlert(ctx)
			logError(client, "volumedown", err)
			return err
		}

		go displayTextAsTitleForTwoSeconds(ctx, client, fmt.Sprintf("%d", volume))

		return nil
	})

	volumeDownAction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	volumeDownAction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)
}


func displayTextAsTitleForTwoSeconds (ctx context.Context, client *streamdeck.Client, text string) {
	client.SetTitle(ctx, text, streamdeck.HardwareAndSoftware)
	time.Sleep(2 * time.Second)
	client.SetTitle(ctx, "", streamdeck.HardwareAndSoftware)
}