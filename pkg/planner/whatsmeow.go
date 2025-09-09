package planner

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"

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
		log.Fatalln("[ERROR} NewWhatsmeowClient: Failed to get device storage container:, err")
	}

	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatalln("[ERROR] NewWhatsmeowClient: Failed to get device from container:, err")
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
			log.Fatalln("[ERROR] Client failed to connect: ", err)
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
			log.Fatalln("[ERROR] Client failed to connect: ", err)
		}
	}
}

func PhoneNumbersToJIDs(client *whatsmeow.Client, phoneNumbers []string) []types.JID {
	JIDs := []types.JID{}

	isOnWhatsAppRes, err := client.IsOnWhatsApp(phoneNumbers)
	if err != nil {
		log.Fatalln("[ERROR] Client failed to get isOnWhatsApp: ", err)
	}

	for i, whatsAppRes := range isOnWhatsAppRes {
		relatedPhone := phoneNumbers[i]

		if whatsAppRes.IsIn {
			fmt.Printf("[INFO] %s found on whatsapp", relatedPhone)
			JIDs = append(JIDs, whatsAppRes.JID)

		} else {
			fmt.Printf("[INFO] %s is not on whatsapp", relatedPhone)
			JIDs = append(JIDs, types.JID{})
		}
	}
	return JIDs
}

func GroupExists(client *whatsmeow.Client, phoneNumbers []string, groupName string) bool {
	JIDs := PhoneNumbersToJIDs(client, phoneNumbers)

	joinedGroups, err := client.GetJoinedGroups()
	if err != nil {
		log.Fatalln("[ERROR] Client failed to get joinedGroups: ", err)
	}

	for _, group := range joinedGroups {
		for _, groupMember := range group.Participants {
			if !slices.Contains(JIDs, groupMember.JID) {
				return false
			}
		}
		if group.Name != groupName {
			return false
		}
	}

	return true
}
