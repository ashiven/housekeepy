package main

import (
	"context"
	"fmt"
	"household-planner/pkg/backend"
	"household-planner/pkg/planner"
	"os"
	"slices"
	"time"

	"github.com/kr/pretty"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func main() {
	fmt.Println("[INFO] Starting Household Planner...")
	debug := len(os.Args) > 1 && slices.Contains(os.Args, "-d")
	useWhatsApp := len(os.Args) > 1 && slices.Contains(os.Args, "-w")

	myHousehold, err := planner.NewHousehold()
	if err != nil {
		fmt.Println("[ERROR] Failed to create household:", err)
		return
	}

	backend.SetHousehold(myHousehold)
	go backend.StartServer()

	for {
		if debug {
			fmt.Println("[DEBUG] Starting next day in one minute...: ")
			time.Sleep(1 * time.Minute)
		} else {
			planner.WaitUntilNoon()
		}

		fmt.Println("[INFO] A new day has started, assigning tasks...")

		myHousehold.ClearAssignments()
		myHousehold.UpdateCurrentMember()
		myHousehold.AssignDailyTasks()
		myHousehold.AssignWeeklyTasks()
		myHousehold.AssignMonthlyTasks()

		// Send messages via whatsmeow
		if useWhatsApp {
			client := planner.NewWhatsmeowClient()

			// NOTE: Needs QR Login via terminal on first startup
			planner.Login(client)

			phoneNumbers := myHousehold.PhoneNumbers()
			JIDs := planner.PhoneNumbersToJIDs(client, phoneNumbers)

			for i, member := range myHousehold.Members {
				assignedTasks := myHousehold.GetAssignedTasks(member)
				message := planner.CreateDailyTaskMessage(assignedTasks, member)

				if debug {
					pretty.Println("[DEBUG] member: ", member)
					pretty.Println("[DEBUG] message: ", message)
					pretty.Println("[DEBUG] phoneNumber: ", phoneNumbers[i])
					pretty.Println("[DEBUG] JID: ", JIDs[i])

				} else {
					JID := JIDs[i]
					waMessage := &waE2E.Message{Conversation: proto.String(message)}
					_, err := client.SendMessage(context.Background(), JID, waMessage)
					if err != nil {
						fmt.Printf("[ERROR] Failed to deliver message to %s.\n", member.Name)
					}
				}
			}

			// Send messages via twilio sms
		} else {
			client := planner.InitializeTwilioClient()
			for _, member := range myHousehold.Members {
				assignedTasks := myHousehold.GetAssignedTasks(member)
				planner.SendMessageSms(client, member, assignedTasks, debug)
			}
		}
	}
}
