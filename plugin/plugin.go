package plugin

import (
	"context"
	"github.com/samwho/streamdeck"
	"os"
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

	go StartTicker()

	return client.Run()
}

func setup(client *streamdeck.Client) {
	setupConfigurationAction(client)
	setupVolumeAction(client)
	setupPlaytoggleAction(client)
	setupTrackActions(client)
	setupAlbumArtAction(client)
}
