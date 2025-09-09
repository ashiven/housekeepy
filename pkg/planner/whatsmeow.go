package planner

// TODO: imports
import (
	"github.com/lucklrj/whatsmeow"
 	"github.com/lucklrj/whatsmeow/store/sqlstore"
 	"github.com/lucklrj/whatsmeow/types/events"
)


func NewClient() {
	dbLogger := waLogger.Stdout("Database", "DEBUG", true)

	container, err := sqlstore.New("sqlite3", "file:devicestore.db?_foreign_keys=on", dbLogger)
	if err != nil {
		log.Fatalln("newClient: Failed to get device storage container:, err") }

	device, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalln("newClient: Failed to get device from container:, err")
	}

	clientLogger := waLogger.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(device, clientLogger)

	return client 
}


func Login(client any) {
	if client.Store.ID == nil {
		qrChannel := client.getQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			log.Fatalln("Client failed to connect: ", err)
		}

		for event := range QRChannel {
			// Login with QR Code
			if event.Event == "code" {
				fmt.Println("Login with QR Code: ", event.Code)
				qrterminal.GenerateHalfBlock(event.Code, qrterminal.L, os.stdout)

			// Login event
			} else {
				fmt.Println("Login event: ", event.Event)
			}
		}

	// Already logged in
	} else {
		err := client.Connect()
		if err != nil {
			log.Fatalln("Client failed to connect: ", err)
		}
	}
}

func PhoneNumbersToJIDs(client any, phoneNumbers []string) []types.JID {
	JIDs := []types.JID{}

	isOnWhatsAppRes := client.IsOnWhatsApp(phoneNumbers) 
	for i, whatsAppRes := range isOnWhatsAppRes {
		relatedPhone := phoneNumbers[i]

		if whatsAppRes.IsIn {
			fmt.Printf("%s found on whatsapp", relatedPhone)
			JIDs = append(JIDs, whatsAppRes.JID)

		} else {
			fmt.Printf("%s is not on whatsapp", relatedPhone)
			JIDs = append(JIDs, -1)
		}
	}
	return JIDs
}

func GroupExists(client any, phoneNumbers []string, groupName string) bool {
	JIDs := PhoneNumbersToJIDs(client, phoneNumbers)

	joinedGroups := client.GetJoinedGroups()
	for group := range joinedGroups {
		for groupMember := range group.Participants {
			if !slices.Contains(JIDs, groupMember.JID) {
				return false
			}
		}
		if group.GroupName != groupName {
			return false
		}
		return true
	}
}
