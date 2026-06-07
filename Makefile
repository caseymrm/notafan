APP=NotAFan
EXE=notafan
VERSION=v1.0.1
IDENTITY=Developer ID Application: Rational Creation LLC (AP2AEA9WAW)
IDENTIFIER=notafan.caseymrm.github.com

SDK=$(shell xcrun --sdk macosx --show-sdk-path)
BUILD=build
APPDIR=$(APP).app
MACOS=$(APPDIR)/Contents/MacOS

# menuet's UNUserNotificationCenter requires macOS 10.14+; arm64 requires 11+
AMD64_MIN=10.14
ARM64_MIN=11.0

GO=go
GOFLAGS=-trimpath -ldflags="-s -w"

.PHONY: all amd64 arm64 universal app sign adhoc notarize zip release clean test

all: app

$(BUILD):
	mkdir -p $(BUILD)

test:
	$(GO) test ./...

amd64: $(BUILD)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
	  CGO_CFLAGS="-mmacosx-version-min=$(AMD64_MIN) -isysroot $(SDK)" \
	  CGO_LDFLAGS="-mmacosx-version-min=$(AMD64_MIN)" \
	  $(GO) build $(GOFLAGS) -o $(BUILD)/$(EXE)-amd64 .

arm64: $(BUILD)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 \
	  CGO_CFLAGS="-mmacosx-version-min=$(ARM64_MIN) -isysroot $(SDK)" \
	  CGO_LDFLAGS="-mmacosx-version-min=$(ARM64_MIN)" \
	  $(GO) build $(GOFLAGS) -o $(BUILD)/$(EXE)-arm64 .

universal: amd64 arm64
	lipo -create -output $(BUILD)/$(EXE) $(BUILD)/$(EXE)-amd64 $(BUILD)/$(EXE)-arm64
	lipo -info $(BUILD)/$(EXE)

app: universal
	cp $(BUILD)/$(EXE) $(MACOS)/$(EXE)

# Ad-hoc sign (no Developer ID required). Lets the bundled binary run locally
# after the security prompt; cannot be notarized.
adhoc: app
	codesign --force --deep --sign - $(APPDIR)

sign: app
	codesign --force --options runtime --timestamp \
	  --sign "$(IDENTITY)" $(APPDIR)

zip:
	rm -f $(BUILD)/$(APP).app.zip
	ditto -c -k --keepParent $(APPDIR) $(BUILD)/$(APP).app.zip

# Requires a keychain profile named AC_NOTARY: xcrun notarytool store-credentials AC_NOTARY
notarize: sign zip
	xcrun notarytool submit $(BUILD)/$(APP).app.zip --keychain-profile AC_NOTARY --wait
	xcrun stapler staple $(APPDIR)
	$(MAKE) zip

release: zip
	gh release create $(VERSION) $(BUILD)/$(APP).app.zip \
	  --title "$(VERSION) — Apple Silicon universal build" \
	  --notes-file RELEASE_NOTES.md

clean:
	rm -rf $(BUILD)
