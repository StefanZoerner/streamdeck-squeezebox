package actions

import (
	"context"
	"encoding/json"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	sdcontext "github.com/samwho/streamdeck/context"
	"strconv"
)

type ConfigurationDataFromPI struct {
	Command  string `json:"command"`
	Hostname string `json:"hostname"`
	CliPort  string `json:"cli_port"`
	HttpPort string `json:"http_port"`
}

type ConfigurationMessage struct {
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}

func SetupConfigurationAction(client *streamdeck.Client) {

	configureaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.configure")
	configureaction.RegisterHandler(streamdeck.WillAppear, general.WillAppearRequestGlobalSettingsHandler)
	configureaction.RegisterHandler(streamdeck.SendToPlugin, configHanderSendToPlugin)
}

func configHanderSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	fromPI := ConfigurationDataFromPI{}
	if err := json.Unmarshal(event.Payload, &fromPI); err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	if fromPI.Command == "setConnection" {

		newGlobalSettings := general.PluginGlobalSettings{}
		newGlobalSettings.Hostname = fromPI.Hostname
		newGlobalSettings.CLIPort, _ = strconv.Atoi(fromPI.CliPort)
		newGlobalSettings.HTTPPort, _ = strconv.Atoi(fromPI.HttpPort)

		globalCtx := sdcontext.WithContext(context.Background(), general.PluginUUID)
		if err := client.SetGlobalSettings(globalCtx, newGlobalSettings); err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		// Enforce Reload of Global S3ttings via an Event
		if err := client.GetGlobalSettings(globalCtx); err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

	} else if fromPI.Command == "testConnection" {

		hostname := fromPI.Hostname
		cliPort, _ := strconv.Atoi(fromPI.CliPort)
		httpPort, _ := strconv.Atoi(fromPI.HttpPort)

		conProps := squeezebox.NewConnectionProperties(hostname, httpPort, cliPort)

		err := squeezebox.CheckConnectionCLI(conProps)
		if err != nil {
			client.ShowAlert(ctx)

			msgPayload := ConfigurationMessage{
				Type:    "caution",
				Summary: "Failed.",
				Content: err.Error(),
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
}
