package general

import (
	"context"
	"encoding/json"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	sdcontext "github.com/samwho/streamdeck/context"
)

// PluginGlobalSettings is stored as a Singleton

type PluginGlobalSettings struct {
	Hostname          string `json:"hostname"`
	CLIPort           int    `json:"cli_port"`
	HTTPPort          int    `json:"http_port"`
	DefaultPlayerID   string `json:"default_player_id"`
	DefaultPlayerName string `json:"default_player_name"`
}

var instance *PluginGlobalSettings

func init() {
	instance = &PluginGlobalSettings{

		// Default values
		Hostname:          "",
		CLIPort:           9090,
		HTTPPort:          9000,
		DefaultPlayerID:   "",
		DefaultPlayerName: "",
	}
}

var PluginUUID string

func GetPluginGlobalSettings() *PluginGlobalSettings {
	return instance
}

func (pgs PluginGlobalSettings) ConnectionProps() squeezebox.ConnectionProperties {
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
		LogErrorWithEvent(client, event, err)
		return err
	}

	settingsFromPayload := PluginGlobalSettings{}
	if err := json.Unmarshal(payload.Settings, &settingsFromPayload); err != nil {
		LogErrorWithEvent(client, event, err)
		return err
	}

	// Store Global settings in Server settings
	//
	serverSettings := GetPluginGlobalSettings()
	serverSettings.Hostname = settingsFromPayload.Hostname
	serverSettings.CLIPort = settingsFromPayload.CLIPort
	serverSettings.HTTPPort = settingsFromPayload.HTTPPort
	serverSettings.DefaultPlayerID = settingsFromPayload.DefaultPlayerID
	serverSettings.DefaultPlayerName = settingsFromPayload.DefaultPlayerName

	return nil
}

func WillAppearRequestGlobalSettingsHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	// LogEvent(client, event)
	var err error

	global := sdcontext.WithContext(context.Background(), PluginUUID)
	if err = client.GetGlobalSettings(global); err != nil {
		LogErrorWithEvent(client, event, err)
	}

	return err
}
