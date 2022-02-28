package general

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

type PlayerObserver interface {
	GetID() string

	PlaymodeChanged(newMode string)
	AlbumArtChanged(newURL string)
}

func AddOberserverForPlayer(playerID string, o PlayerObserver) int {

	playerSubjectsLock.Lock()
	defer playerSubjectsLock.Unlock()

	ps, ok := playerSubjects[playerID]
	if !ok {
		ps = newPlayerSubject(playerID)
		playerSubjects[playerID] = ps
	}

	return ps.register(o)
}

func RemoveOberserverForPlayer(playerID string, o PlayerObserver) int {
	var count int

	playerSubjectsLock.Lock()
	defer playerSubjectsLock.Unlock()

	ps, ok := playerSubjects[playerID]
	if ok {
		count = ps.deregister(o)
	}

	return count
}

func RemoveOberserverForAllPlayers(o PlayerObserver) {
	playerSubjectsLock.Lock()
	defer playerSubjectsLock.Unlock()

	for _, subject := range playerSubjects {
		subject.deregister(o)
	}
}

type playerSubject struct {
	observerList []PlayerObserver
	name         string
	playMode     string
	currentURL   string
}

func newPlayerSubject(playerID string) *playerSubject {
	return &playerSubject{
		name:       playerID,
		playMode:   "-",
		currentURL: "-",
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

func (ps *playerSubject) register(o PlayerObserver) int {
	ps.observerList = append(ps.observerList, o)
	return len(ps.observerList)
}

func (ps *playerSubject) deregister(o PlayerObserver) int {
	ps.observerList = removeFromslice(ps.observerList, o)
	return len(ps.observerList)
}

func (ps *playerSubject) notifyPlaymodeChanged(newPlayMode string) {
	for _, observer := range ps.observerList {
		observer.PlaymodeChanged(newPlayMode)
	}
}

func (ps *playerSubject) notifyAlbumArtChanged(newURL string) {
	for _, observer := range ps.observerList {
		observer.AlbumArtChanged(newURL)
	}
}

func removeFromslice(observerList []PlayerObserver, observerToRemove PlayerObserver) []PlayerObserver {
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.GetID() == observer.GetID() {
			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}

func StartTicker() {

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

						cp := GetPluginGlobalSettings().ConnectionProps()
						if playerID == "" {
							playerID = GetPluginGlobalSettings().DefaultPlayerID
						}

						status, _ := squeezebox.GetPlayerMode(cp, playerID)
						subject.updatePlayMode(status)

						url, _ := squeezebox.GetCurrentArtworkURL(cp, playerID)
						subject.updateAlbumArt(url)
					}
				}
			}
		}
	}()

}
