package qrgenerator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

// Generates a QR code for WiFi credentials and saves it to a file.
func GenerateWiFiQRCodeToFile(ssid, password, networkType, outputFilePath string, size int) error {
	qrImage, err := GenerateWiFiQRCode(ssid, password, networkType, size)
	if err != nil {
		return err
	}

	// Save QR code to file
	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, qrImage)
}

// Generates a QR code and returns it as an in-memory image.
func GenerateWiFiQRCode(ssid, password, networkType string, size int) (image.Image, error) {
	wifiString := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", networkType, ssid, password)

	qrCode, err := qrcode.New(wifiString, qrcode.High)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	qrCode.DisableBorder = true

	qrImage := qrCode.Image(size)
	return qrImage, nil
}

// Adds a logo to a QR code file and saves the result.
func AddLogoToQRFile(qrImagePath, logoPath, outputFilePath string) error {
	qrImage, err := loadImageFromFile(qrImagePath)
	if err != nil {
		return err
	}

	logoImage, err := loadImageFromFile(logoPath)
	if err != nil {
		return err
	}

	brandedQR := AddLogoToQR(qrImage, logoImage)

	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, brandedQR)
}

// Adds a logo to the QR code and returns an in-memory image.
func AddLogoToQR(qrImage image.Image, logoImage image.Image) image.Image {
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

	return dc.Image()
}

// Helper function to load an image from a file
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

// Converts an image.Image to a PNG byte slice.
func EncodeImageToPNGBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
