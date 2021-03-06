NAME = de.szoerner.streamdeck.squeezebox
INSTALLDIR = "/Users/stefanz/Library/Application Support/com.elgato.StreamDeck/Plugins"
BUILDDIR = build
RELEASEDIR = release
SDPLUGINDIR = "$(BUILDDIR)/$(NAME).sdPlugin"

.DEFAULT_GOAL := install

.PHONY:fmt
fmt:
	go fmt ./...

.PHONY:build
build: fmt macexecutable windowsexecutable

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
	cp streamdeck-squeezebox.exe $(SDPLUGINDIR)

uninstall:
	rm -rf $(INSTALLDIR)/$(NAME).sdPlugin

install: uninstall sdplugin
	mv $(SDPLUGINDIR) $(INSTALLDIR)/$(NAME).sdPlugin

distribute: sdplugin
	DistributionTool -b -i $(SDPLUGINDIR) -o ~/Desktop/
