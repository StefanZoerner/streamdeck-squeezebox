package plugin

import (
	"github.com/samwho/streamdeck"
	"strings"
)

func logError(client *streamdeck.Client, action string, err error) {
	client.LogMessage("Error in "+action + ": "+err.Error())
}

func logEvent(client *streamdeck.Client, event streamdeck.Event) {

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