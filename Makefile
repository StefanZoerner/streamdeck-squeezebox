NAME = de.szoerner.streamdeck.squeezebox
INSTALLDIR = "/Users/stefanz/Library/Application Support/com.elgato.StreamDeck/Plugins"
BUILDDIR = build
RELEASEDIR = release
SDPLUGINDIR = "$(BUILDDIR)/$(NAME).sdPlugin"

fmt:
	go fmt ./...
.PHONY:fmt

build: fmt macexecutable

windowsexecutable:
	GOOS=windows GOARCH=amd64 go build -o streamdeck-squeezebox.exe

macexecutable:
	go build -mod=vendor -o streamdeck-squeezebox

test:
	go test -v -mod=vendor github.com/StefanZoerner/streamdeck-squeezebox/squeezebox

sdplugin: build
	rm -rf $(SDPLUGINDIR)
	mkdir -p $(SDPLUGINDIR)
	cp manifest.json $(SDPLUGINDIR)
	cp -r assets $(SDPLUGINDIR)
	cp -r html $(SDPLUGINDIR)
	cp streamdeck-squeezebox $(SDPLUGINDIR)

uninstall:
	rm -rf $(INSTALLDIR)/$(NAME).sdPlugin

install: uninstall sdplugin
	mv $(SDPLUGINDIR) $(INSTALLDIR)/$(NAME).sdPlugin

