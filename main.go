package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

	// WIFI:T:<NetworkType>;S:<SSID>;P:<Password>;;
	wifiString := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", networkType, ssid, password)

	// QR encode and save to disk
	err = qrcode.WriteFile(wifiString, qrcode.Medium, 256, "wifi-qr.png")
	if err != nil {
		fmt.Println("Failed to generate QR code:", err)
		return
	}

	fmt.Println("QR code generated successfully: wifi-qr.png")
}
