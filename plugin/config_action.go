package plugin

import (
	"context"
	"encoding/json"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	sdcontext "github.com/samwho/streamdeck/context"
	"strconv"
)

type ConfigurationDataFromPI struct {
	Command  string `json:"command"`
	Hostname string `json:"hostname"`
	CliPort  string `json:"cli_port"`
}

type ConfigurationMessage struct {
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}

func setupConfigurationAction(client *streamdeck.Client) {

	configureaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.configure")
	configureaction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)

	configureaction.RegisterHandler(streamdeck.SendToPlugin, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		fromPI := ConfigurationDataFromPI{}
		if err := json.Unmarshal(event.Payload, &fromPI); err != nil {
			logError(client, event, err)
			return err
		}

		if fromPI.Command == "setConnection" {

			newGlobalSettings := PluginGlobalSettings{}
			newGlobalSettings.Hostname = fromPI.Hostname
			newGlobalSettings.CliPort, _ = strconv.Atoi(fromPI.CliPort)

			globalCtx := sdcontext.WithContext(context.Background(), pluginUUID)
			if err := client.SetGlobalSettings(globalCtx, newGlobalSettings); err != nil {
				logError(client, event, err)
				return err
			}

			// Enforce Reload of Global S3ttings via an Event
			if err := client.GetGlobalSettings(globalCtx); err != nil {
				logError(client, event, err)
				return err
			}

		} else if fromPI.Command == "testConnection" {

			hostname := fromPI.Hostname
			cliPort, _ := strconv.Atoi(fromPI.CliPort)

			conProps := squeezebox.NewConnectionProperties(hostname, 9002, cliPort)

			error := squeezebox.CheckConnectionCLI(conProps)
			if error != nil {
				client.ShowAlert(ctx)

				msgPayload := ConfigurationMessage{
					Type:    "caution",
					Summary: "Failed.",
					Content: error.Error(),
				}
				client.SendToPropertyInspector(ctx, msgPayload)

			} else {
				client.ShowOk(ctx)

				msgPayload := ConfigurationMessage{
					Type:    "info",
					Summary: "Success.",
					Content: "Connection to Logitech Media Server successfully establiished.",
				}
				client.SendToPropertyInspector(ctx, msgPayload)
			}
		}

		return nil
	})
}
