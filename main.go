package main

import (
	"fmt"
	"household-planner/pkg/backend"
	"household-planner/pkg/planner"
	"os"
	"time"
)

const whatsAppGroup = "Haushaltsplaner"

func main() {
	fmt.Println("[INFO] Starting Household Planner...")
	debug := len(os.Args) > 1 && os.Args[1] == "-d"
	useWhatsApp := len(os.Args) > 1 && os.Args[1] == "-w"

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
			// fmt.Printf("%# v\n", pretty.Formatter(myHousehold))
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
			client := planner.NewClient()

			//NOTE: Needs QR Login via terminal on first startup
			planner.Login(client)
			JIDs := planner.PhoneNumbersToJIDs

			if !planner.GroupExists(client, members, whatsAppGroup) {
				client.CreateGroup(whatsAppGroup, JIDs)
			}

			for i, member := range myHouseHold.Members {
				assignedTasks := myHouseHold.GetAssignedTasks(member)
				message := planner.CreateTaskMessage(member, assignedTasks)

				JID := JIDs[i]
				_, err := client.SendMessage(JID, message)
				if err != nil {
					fmt.Printf("Failed to deliver message to %s.\n", member.Name)
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
