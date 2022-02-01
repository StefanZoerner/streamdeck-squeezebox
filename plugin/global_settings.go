package plugin

import (
	"context"
	"encoding/json"
	"github.com/samwho/streamdeck"
	sdcontext "github.com/samwho/streamdeck/context"
)

// PluginGlobalSettings is stored as a Singleton

type PluginGlobalSettings struct {
	Hostname string `json:"hostname"`
	CliPort  int    `json:"cli_port"`
}

var instance *PluginGlobalSettings

func init() {
	instance = &PluginGlobalSettings{
		// Default values
		Hostname: "hostname",
		CliPort:  9090,
	}
}

func GetPluginGlobalSettings() *PluginGlobalSettings {
	return instance
}

func DidReceiveGlobalSettingsHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	payload := streamdeck.DidReceiveGlobalSettingsPayload{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logError(client, event, err)
		return err
	}

	settingsFromPayload := PluginGlobalSettings{}
	if err := json.Unmarshal(payload.Settings, &settingsFromPayload); err != nil {
		logError(client, event, err)
		return err
	}

	// Store Global Settings in Server Settings
	//
	serverSettings := GetPluginGlobalSettings()
	serverSettings.Hostname = settingsFromPayload.Hostname
	serverSettings.CliPort = settingsFromPayload.CliPort

	return nil
}

func WillAppearRequestGlobalSettingsHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	global := sdcontext.WithContext(context.Background(), pluginUUID)
	if err := client.GetGlobalSettings(global); err != nil {
		logError(client, event, err)
		return err
	}

	return nil
}
