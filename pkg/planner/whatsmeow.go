package planner

import (
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLogger "go.mau.fi/whatsmeow/util/log"
)

func NewWhatsmeowClient() *whatsmeow.Client {
	dbLogger := waLogger.Stdout("Database", "DEBUG", true)
	ctx := context.Background()

	container, err := sqlstore.New(ctx, "sqlite3", "file:devicestore.db?_foreign_keys=on", dbLogger)
	if err != nil {
		fmt.Println("[ERROR] NewWhatsmeowClient: ", err)
	}

	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		fmt.Println("[ERROR] NewWhatsmeowClient: ", err)
	}

	clientLogger := waLogger.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(device, clientLogger)

	return client
}

func Login(client *whatsmeow.Client) {
	if client.Store.ID == nil {
		qrChannel, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			fmt.Println("[ERROR] Login: ", err)
		}

		for event := range qrChannel {
			// Login with QR Code
			if event.Event == "code" {
				fmt.Println("[INFO] Login with QR Code: ", event.Code)
				qrterminal.GenerateHalfBlock(event.Code, qrterminal.L, os.Stdout)

				// Login event
			} else {
				fmt.Println("[INFO] Login event: ", event.Event)
			}
		}

		// Already logged in
	} else {
		err := client.Connect()
		if err != nil {
			fmt.Println("[ERROR] Login: ", err)
		}
	}
}

func PhoneNumbersToJIDs(client *whatsmeow.Client, phoneNumbers []string) map[string]types.JID {
	JIDs := map[string]types.JID{}

	isOnWhatsAppRes, err := client.IsOnWhatsApp(context.Background(), phoneNumbers)
	if err != nil {
		fmt.Println("[ERROR] IsOnWhatsApp: ", err)
	}

	for _, whatsAppRes := range isOnWhatsAppRes {
		phoneNumber := "+" + whatsAppRes.JID.User

		if whatsAppRes.IsIn {
			fmt.Printf("[INFO] %s found on whatsapp\n", phoneNumber)
			JIDs[phoneNumber] = whatsAppRes.JID

		} else {
			fmt.Printf("[INFO] %s is not on whatsapp\n", phoneNumber)
			JIDs[phoneNumber] = types.JID{}
		}
	}
	return JIDs
}
