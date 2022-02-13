package plugin

import (
	"context"
	"encoding/json"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	sdcontext "github.com/samwho/streamdeck/context"
)

// PluginGlobalSettings is stored as a Singleton

type PluginGlobalSettings struct {
	Hostname string `json:"hostname"`
	CLIPort  int    `json:"cli_port"`
	HTTPPort int    `json:"http_port"`
}

var instance *PluginGlobalSettings

func init() {
	instance = &PluginGlobalSettings{

		// Default values
		Hostname: "hostname",
		CLIPort:  9090,
		HTTPPort: 9000,
	}
}

func GetPluginGlobalSettings() *PluginGlobalSettings {
	return instance
}

func (pgs PluginGlobalSettings) connectionProps() squeezebox.ConnectionProperties {
	cp := squeezebox.ConnectionProperties{
		Hostname: pgs.Hostname,
		HTTPPort: pgs.HTTPPort,
		CLIPort:  pgs.CLIPort,
	}

	return cp
}

func DidReceiveGlobalSettingsHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	// LogEvent(client, event)

	payload := streamdeck.DidReceiveGlobalSettingsPayload{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	settingsFromPayload := PluginGlobalSettings{}
	if err := json.Unmarshal(payload.Settings, &settingsFromPayload); err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	// Store Global Settings in Server Settings
	//
	serverSettings := GetPluginGlobalSettings()
	serverSettings.Hostname = settingsFromPayload.Hostname
	serverSettings.CLIPort = settingsFromPayload.CLIPort
	serverSettings.HTTPPort = settingsFromPayload.HTTPPort

	return nil
}

func WillAppearRequestGlobalSettingsHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	// LogEvent(client, event)
	var err error

	global := sdcontext.WithContext(context.Background(), pluginUUID)
	if err = client.GetGlobalSettings(global); err != nil {
		general.LogErrorWithEvent(client, event, err)
	}

	return err
}
