# go-wifi-qr-code-generator

![Image](https://github.com/user-attachments/assets/ca634509-e3d3-40dd-a205-64ed2deafd82)

This projects generates QR codes for provided WiFi credentials. Users may then use their camera app to connect to the network in an easy way. The project may be used both as a library and as a standalone program.

QR codes include error correction to allow partial recovery of lost or damaged code areas. This is particularly useful if a QR code gets dirty, scratched, or has a logo overlay. We overlay a logo with 20% of the QR code size and error correction adds redundancy for up to 30%, thus this works.

## Library usage

### File based example

```go
package main

import (
	"fmt"

	"github.com/StarHack/go-wifi-qr-code-generator/qrgenerator"
)

func main() {
	ssid := "example network"
	password := "test123"
	networkType := "WPA"
	qrFilePath := "wifi-qr.png"
	outputFilePath := "wifi-qr-branded.png"
	logoFilePath := "logo.png"

	err := qrgenerator.GenerateWiFiQRCodeToFile(ssid, password, networkType, qrFilePath, 512)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}

	err = qrgenerator.AddLogoToQRFile(qrFilePath, logoFilePath, outputFilePath)
	if err != nil {
		fmt.Println("Error adding logo to QR code:", err)
		return
	}

	fmt.Println("QR code generated successfully:", outputFilePath)
}
```

### In-memory example

```go
package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/StarHack/go-wifi-qr-code-generator/qrgenerator"
)

func main() {
	ssid := "example network"
	password := "test123"
	networkType := "WPA"
	logoImg, err := loadImageFromFile("logo.png")
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}

	qrImg, err := qrgenerator.GenerateWiFiQRCode(ssid, password, networkType, 512)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}

	newImg := qrgenerator.AddLogoToQR(qrImg, logoImg)

	// use newImg

	fmt.Println("QR code generated successfully:")
}

func loadImageFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}


```

## Standalone usage

```
git clone https://github.com/StarHack/go-wifi-qr-code-generator.git
cd go-wifi-qr-code-generator
go get
go build
./go-wifi-qr-code-generator
```

A user types:

- Their WiFi name (SSID)
- WiFi password
- Type (one of: `WPA`, `WEP`, `nopass`)

The program afterwards generates a file called `wifi-qr.png` (and `wifi-qr-branded.png` if `logo.png` exists).

