package main

import (
	"context"
	"encoding/json"
	"github.com/samwho/streamdeck"
	"log"
	"os"
	"strconv"
	"strings"
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

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func run(ctx context.Context) error {
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
		SetPlayerStatus(player, "play")
		return nil
	})

	pauseaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.pause")
	pauseaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		SetPlayerStatus(player, "pause")
		return nil
	})

	playtoggleaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.playtoggle")
	playtoggleaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		TogglePlayerStatus(player)
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
		LogEvent(client, event)

		payload := streamdeck.KeyDownPayload{}
		err := json.Unmarshal(event.Payload, &payload)
		if (err != nil) {
			LogError(client, "configureaction", err)
			return err
		}

		settings := Settings{}
		err = json.Unmarshal(payload.Settings, &settings)
		if err != nil {
			LogError(client, "configureaction", err)
			return err
		}

		err = CheckConnectionToPlayer(settings.Hostname, settings.Port, settings.PlayerId)
		if (err != nil) {
			LogError(client, "configure", err)
			client.ShowAlert(ctx)
		} else {
			client.ShowOk(ctx)
		}

		return nil
	})

	configureaction.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		LogEvent(client, event)

		payload := streamdeck.WillAppearPayload{}
		err := json.Unmarshal(event.Payload, &payload)
		if (err != nil) {
			LogError(client, "configureaction", err)
			return err
		}

		sbox_settings := Settings{}
		err = json.Unmarshal(payload.Settings, &sbox_settings)
		if (err != nil) {
			LogError(client, "configureaction", err)
			return err
		}

		PreloadSettings(&sbox_settings)
		err = client.SetSettings(ctx, sbox_settings)
		if (err != nil) {
			LogError(client, "configureaction", err)
			return err
		}

		return nil
	})

	configureaction.RegisterHandler(streamdeck.SendToPlugin, getSettingsFromPIHandler)
}

func getSettingsFromPIHandler (ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	LogEvent(client, event)

	piSettings := PropertyInspectorSettings{}
	err := json.Unmarshal(event.Payload, &piSettings)
	if err != nil {
		LogError(client, "configureaction", err)
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

func LogError(client *streamdeck.Client, action string, err error) {
	client.LogMessage("Error in "+action + ": "+err.Error())
}

func LogEvent(client *streamdeck.Client, event streamdeck.Event) {

	// Determine last part of dot divided action name
	action_name := "???"
	actionParts := strings.Split(event.Action, ".");
	if len(actionParts) > 0 {
		action_name = actionParts[len(actionParts)-1]
	}

	msg := action_name + " " + event.Event + " "
	client.LogMessage("Event : " +msg)
	pl, _ := event.Payload.MarshalJSON()
	client.LogMessage("Payload: "+string(pl)+"\n")
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