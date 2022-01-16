package plugin

import (
	"context"
	"os"

	"github.com/samwho/streamdeck"
)

var pluginUUID string

func Run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	pluginUUID = params.PluginUUID
	client := streamdeck.NewClient(ctx, params)

	client.RegisterHandler(streamdeck.DidReceiveGlobalSettings, DidReceiveGlobalSettingsHandler)
	setup(client)

	return client.Run()
}

func setup(client *streamdeck.Client) {
	setupConfigurationAction(client)
	setupVolumeActions(client)
	setupPlaymodeActions(client)
}


