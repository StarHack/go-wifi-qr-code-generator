package main

import (
	"bufio"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

var reader = bufio.NewReader(os.Stdin)

func readUserInput(instructions string, validTypes []string) (string, error) {
	fmt.Print(instructions)
	var userInput string
	var err error
	for len(userInput) == 0 {
		userInput, err = reader.ReadString('\n')

		if err != nil {
			return "", err
		}

		userInput = strings.TrimSpace(userInput)

		if len(userInput) > 0 && len(validTypes) > 0 {
			for _, validType := range validTypes {
				if strings.EqualFold(validType, userInput) {
					userInput = validType
					return userInput, nil
				}
			}
			userInput = ""
		}

		if len(userInput) == 0 {
			fmt.Print(instructions)
		}
	}
	return userInput, nil
}

func addLogoToQR(qrImagePath, logoPath, outputPath string) error {
	qrFile, err := os.Open(qrImagePath)
	if err != nil {
		return err
	}
	defer qrFile.Close()

	qrImage, err := png.Decode(qrFile)
	if err != nil {
		return err
	}

	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		err := os.Rename(qrImagePath, outputPath)
		if err != nil {
			return err
		}
		fmt.Println("No logo found. QR code generated without branding.")
		return nil
	}

	logoFile, err := os.Open(logoPath)
	if err != nil {
		return err
	}
	defer logoFile.Close()

	logoImage, err := png.Decode(logoFile)
	if err != nil {
		return err
	}

	logoBounds := logoImage.Bounds()
	logoWidth, logoHeight := logoBounds.Dx(), logoBounds.Dy()

	const logoScale = 0.20
	const marginScale = 1.4
	qrSize := qrImage.Bounds().Dx()
	maxLogoSize := float64(qrSize) * logoScale

	var newLogoWidth, newLogoHeight uint
	aspectRatio := float64(logoWidth) / float64(logoHeight)

	if aspectRatio > 1 {
		newLogoWidth = uint(maxLogoSize)
		newLogoHeight = uint(maxLogoSize / aspectRatio)
	} else {
		newLogoHeight = uint(maxLogoSize)
		newLogoWidth = uint(maxLogoSize * aspectRatio)
	}

	logoResized := resize.Resize(newLogoWidth, newLogoHeight, logoImage, resize.Lanczos3)

	dc := gg.NewContext(qrSize, qrSize)
	dc.DrawImage(qrImage, 0, 0)

	circleRadius := maxLogoSize * marginScale
	centerX, centerY := float64(qrSize)/2, float64(qrSize)/2

	dc.SetColor(color.White)
	dc.DrawCircle(centerX, centerY, circleRadius/2)
	dc.Fill()

	x := int(centerX - float64(newLogoWidth)/2)
	y := int(centerY - float64(newLogoHeight)/2)
	dc.DrawImage(logoResized, x, y)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return png.Encode(outputFile, dc.Image())
}

func main() {
	ssid, err := readUserInput("Enter WiFi name (SSID): ", []string{})
	if err != nil {
		panic(err)
	}

	password, err := readUserInput("Enter WiFi password: ", []string{})
	if err != nil {
		panic(err)
	}

	networkType, err := readUserInput("Enter WiFi type (WPA, WEP, nopass): ", []string{"WPA", "WEP", "nopass"})
	if err != nil {
		panic(err)
	}

	wifiString := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", networkType, ssid, password)

	qrFilePath := "wifi-qr.png"
	err = qrcode.WriteFile(wifiString, qrcode.High, 512, qrFilePath)
	if err != nil {
		fmt.Println("Failed to generate QR code:", err)
		return
	}

	logoFilePath := "logo.png"
	outputFilePath := "wifi-qr-branded.png"

	err = addLogoToQR(qrFilePath, logoFilePath, outputFilePath)
	if err != nil {
		fmt.Println("Failed to add logo:", err)
		return
	}

	fmt.Println("QR code generated successfully:", outputFilePath)
}

