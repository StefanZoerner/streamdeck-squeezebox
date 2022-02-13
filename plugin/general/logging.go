package general

import (
	"github.com/samwho/streamdeck"
	"strings"
)

func LogEvent(client *streamdeck.Client, event streamdeck.Event) {

	// Determine last part of dot divided action name
	actionName := getActionNameFromEvent(event)

	msg := actionName + " " + event.Event + " "
	client.LogMessage("Event : " + msg)
	pl, _ := event.Payload.MarshalJSON()

	client.LogMessage("Payload: " + string(pl) + "\n")
}

func LogErrorWithEvent(client *streamdeck.Client, event streamdeck.Event, err error) {
	actionName := getActionNameFromEvent(event)
	client.LogMessage("Error in " + actionName + " " + event.Event + ": " + err.Error())
}

func LogErrorNoEvent(client *streamdeck.Client, err error) {
	client.LogMessage("Error: " + err.Error())
}

func getActionNameFromEvent(event streamdeck.Event) string {
	actionName := "???"
	actionParts := strings.Split(event.Action, ".")
	if len(actionParts) > 0 {
		actionName = actionParts[len(actionParts)-1]
	}
	return actionName
}
