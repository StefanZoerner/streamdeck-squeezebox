NAME = de.szoerner.streamdeck.squeezebox
INSTALLDIR = "/Users/stefanz/Library/Application Support/com.elgato.StreamDeck/Plugins"
BUILDDIR = build
RELEASEDIR = release
SDPLUGINDIR = "$(BUILDDIR)/$(NAME).sdPlugin"

build: macexecutable

windowsexecutable:
	GOOS=windows GOARCH=amd64 go build -o streamdeck-squeezebox.exe

macexecutable:
	go build -o streamdeck-squeezebox

sdplugin: build
	rm -rf $(SDPLUGINDIR)
	mkdir -p $(SDPLUGINDIR)
	cp manifest.json $(SDPLUGINDIR)
	cp -r images $(SDPLUGINDIR)
	cp -r html $(SDPLUGINDIR)
	cp *.html $(SDPLUGINDIR)
	cp streamdeck-squeezebox $(SDPLUGINDIR)

uninstall:
	rm -rf $(INSTALLDIR)/$(NAME).sdPlugin

install: uninstall sdplugin
	mv $(SDPLUGINDIR) $(INSTALLDIR)/$(NAME).sdPlugin

