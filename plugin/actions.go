package plugin

import (
	"context"
	"encoding/json"
	"strconv"
	"os"

	"github.com/samwho/streamdeck"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
)

const player = "00:04:20:22:c2:54"

type PropertyInspectorSettings struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	PlayerId string `json:"player_id"`
}

type Settings struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	PlayerId string `json:"player_id"`
}

func Run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := streamdeck.NewClient(ctx, params)
	setup(client)

	return client.Run()
}

func setup(client *streamdeck.Client) {

	playaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.play")
	playaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		squeezebox.SetPlayerStatus(player, "play")
		return nil
	})

	pauseaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.pause")
	pauseaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		squeezebox.SetPlayerStatus(player, "pause")
		return nil
	})


	playtoggleaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.playtoggle")
	playtoggleaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		status, err := squeezebox.TogglePlayerStatus(player)
		if err != nil {
			client.ShowAlert(ctx)
			logError(client, "playtoggle", err)
			return err
		}

		client.LogMessage("New status: "+status)
		// TODO: Change Key Image

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		p := streamdeck.WillAppearPayload{}
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return err
		}

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.SetSettings, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.SendToPlugin, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {

		return nil

	})

	configureaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.configure")

	// KeyDown
	//
	configureaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		payload := streamdeck.KeyDownPayload{}
		err := json.Unmarshal(event.Payload, &payload)
		if (err != nil) {
			logError(client, "configureaction", err)
			return err
		}

		settings := Settings{}
		err = json.Unmarshal(payload.Settings, &settings)
		if err != nil {
			logError(client, "configureaction", err)
			return err
		}

		err = squeezebox.CheckConnectionToPlayer(settings.Hostname, settings.Port, settings.PlayerId)

		if (err != nil) {
			logError(client, "configure", err)
			client.ShowAlert(ctx)
		} else {
			client.ShowOk(ctx)
		}

		return nil
	})

	configureaction.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		payload := streamdeck.WillAppearPayload{}
		err := json.Unmarshal(event.Payload, &payload)
		if (err != nil) {
			logError(client, "configureaction", err)
			return err
		}

		sbox_settings := Settings{}
		err = json.Unmarshal(payload.Settings, &sbox_settings)
		if (err != nil) {
			logError(client, "configureaction", err)
			return err
		}

		PreloadSettings(&sbox_settings)
		err = client.SetSettings(ctx, sbox_settings)
		if (err != nil) {
			logError(client, "configureaction", err)
			return err
		}

		return nil
	})

	configureaction.RegisterHandler(streamdeck.SendToPlugin, GetSettingsFromPIHandler)
}

func GetSettingsFromPIHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	piSettings := PropertyInspectorSettings{}
	err := json.Unmarshal(event.Payload, &piSettings)
	if err != nil {
		logError(client, "configureaction", err)
		return err
	} else {
		newSettings := Settings{}
		newSettings.Hostname = piSettings.Hostname
		newSettings.Port, _ = strconv.Atoi(piSettings.Port)
		newSettings.PlayerId = piSettings.PlayerId

		client.SetSettings(ctx, newSettings);
	}

	return nil
}


func PreloadSettings(settings *Settings) {

	if settings.Hostname == "" {
		settings.Hostname = "hostname"
	}

	if settings.Port == 0 {
		settings.Port = 9090
	}

	if settings.PlayerId == "" {
		settings.PlayerId = "none"
	}

}
