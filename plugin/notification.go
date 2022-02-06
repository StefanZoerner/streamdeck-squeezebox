package plugin

import (
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"sync"
	"time"
)

// https://golangbyexample.com/observer-design-pattern-golang/

var (
	playerSubjects     map[string]*playerSubject
	playerSubjectsLock sync.Mutex
)

func init() {
	playerSubjects = make(map[string]*playerSubject)
}

type observer interface {
	getID() string

	playmodeChanged(newMode string)
	albumArtChanged(newURL string)
}

type subject interface {
	register(o observer)
	deregister(o observer)

	notifyPlaymodeChanged(newMode string)
	albumArtChanged(newURL string)
}

func addOberserverForPlayer(playerID string, o observer) int {

	playerSubjectsLock.Lock()
	defer playerSubjectsLock.Unlock()

	ps, ok := playerSubjects[playerID]
	if !ok {
		ps = newPlayerSubject(playerID)
		playerSubjects[playerID] = ps
	}

	return ps.register(o)
}

func removeOberserverForPlayer(playerID string, o observer) int {
	var count int

	playerSubjectsLock.Lock()
	defer playerSubjectsLock.Unlock()

	ps, ok := playerSubjects[playerID]
	if ok {
		count = ps.deregister(o)
	}

	return count
}

func removeOberserverForAllPlayers(o observer) {
	playerSubjectsLock.Lock()
	defer playerSubjectsLock.Unlock()

	for _, subject := range playerSubjects {
		subject.deregister(o)
	}
}

type playerSubject struct {
	observerList []observer
	name         string
	playMode     string
	currentURL   string
}

func newPlayerSubject(playerID string) *playerSubject {
	return &playerSubject{
		name:       playerID,
		playMode:   "",
		currentURL: "",
	}
}

func (ps *playerSubject) updatePlayMode(newPlayMode string) {
	if ps.playMode != newPlayMode {
		ps.playMode = newPlayMode
		ps.notifyPlaymodeChanged(newPlayMode)
	}
}

func (ps *playerSubject) updateAlbumArt(newURL string) {
	if ps.currentURL != newURL {
		ps.currentURL = newURL
		ps.notifyAlbumArtChanged(newURL)
	}
}

func (ps *playerSubject) register(o observer) int {
	ps.observerList = append(ps.observerList, o)
	return len(ps.observerList)
}

func (ps *playerSubject) deregister(o observer) int {
	ps.observerList = removeFromslice(ps.observerList, o)
	return len(ps.observerList)
}

func (ps *playerSubject) notifyPlaymodeChanged(newPlayMode string) {
	for _, observer := range ps.observerList {
		observer.playmodeChanged(newPlayMode)
	}
}

func (ps *playerSubject) notifyAlbumArtChanged(newURL string) {
	for _, observer := range ps.observerList {
		observer.albumArtChanged(newURL)
	}
}

func removeFromslice(observerList []observer, observerToRemove observer) []observer {
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.getID() == observer.getID() {
			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}

func startTicker() {

	// https://gobyexample.com/tickers

	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for playerID, subject := range playerSubjects {

					if len(subject.observerList) > 0 {
						status, _ := getPlayMode(playerID)
						subject.updatePlayMode(status)
						url, _ := getAlbumArtURL(playerID)
						subject.updateAlbumArt(url)
					}
				}
			}
		}
	}()

}

func getPlayMode(playerID string) (string, error) {
	gs := GetPluginGlobalSettings()
	status, err := squeezebox.GetPlayerMode(gs.Hostname, gs.CLIPort, playerID)
	return status, err
}

func getAlbumArtURL(playerID string) (string, error) {
	cp := GetPluginGlobalSettings().connectionProps()
	status, err := squeezebox.GetCurrentArtworkURL(cp, playerID)
	return status, err
}
