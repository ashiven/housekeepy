package planner

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func GetEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
	return os.Getenv(key)
}

func InitializeTwilioClient() *twilio.RestClient {
	twilioAccountSid := GetEnvVar("TWILIO_ACCOUNT_SID")
	twilioAuthToken := GetEnvVar("TWILIO_AUTH_TOKEN")
	_ = twilioAccountSid
	_ = twilioAuthToken

	client := twilio.NewRestClient()
	return client
}

func taskNameOrNone(tasks []Assignable, index int) string {
	if index < len(tasks) {
		return tasks[index].GetName()
	}
	return "_"
}

func SendMessageWhatsapp(client *twilio.RestClient, receiver *Member, tasks []Assignable, debug bool) {
	templateSid := GetEnvVar("TEMPLATE_SID")
	sender := GetEnvVar("WHATSAPP_SENDER")
	serviceSid := GetEnvVar("SERVICE_SID")

	ContentVariables, err := json.Marshal(map[string]any{
		"1": receiver.Name,
		"2": "heutigen",
		"3": taskNameOrNone(tasks, 0),
		"4": taskNameOrNone(tasks, 1),
		"5": taskNameOrNone(tasks, 2),
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	if debug {
		fmt.Println("[DEBUG] Template SID:", templateSid)
		fmt.Println("[DEBUG] Service SID:", serviceSid)
		fmt.Println("[DEBUG] Sender:", sender)
		fmt.Println("[DEBUG] Receiver:", receiver.Name, "(", receiver.PhoneNumber, ")")
		fmt.Println("[DEBUG] Sending message with the following content variables:")
		fmt.Println(string(ContentVariables))
		return
	}

	params := &api.CreateMessageParams{}
	params.SetContentSid(templateSid)
	params.SetTo("whatsapp:" + receiver.PhoneNumber)
	params.SetFrom("whatsapp:" + sender)
	params.SetContentVariables(string(ContentVariables))
	params.SetMessagingServiceSid(serviceSid)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Body != nil {
			fmt.Println(*resp.Body)
		} else {
			fmt.Println(resp.Body)
		}
	}
}

func CreateDailyTaskMessage(tasks []Assignable, member *Member) string {
	message := fmt.Sprintf("%s! Deine heutigen Aufgaben sind:\n", member.Name)

	dailyTasks := "\n"
	for _, task := range tasks {
		dailyTasks += fmt.Sprintf("- %s\n", task.GetName())
	}

	return message + dailyTasks
}

func SendMessageSms(client *twilio.RestClient, receiver *Member, tasks []Assignable, debug bool) {
	sender := GetEnvVar("SMS_SENDER")

	message := CreateDailyTaskMessage(tasks, receiver)

	if debug {
		fmt.Println("[DEBUG] Sender:", sender)
		fmt.Println("[DEBUG] Receiver:", receiver.Name, "(", receiver.PhoneNumber, ")")
		fmt.Println("[DEBUG] Sending message:")
		fmt.Println(message)
		return
	}

	params := &api.CreateMessageParams{}
	params.SetTo(receiver.PhoneNumber)
	params.SetFrom(sender)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Body != nil {
			fmt.Println(*resp.Body)
		} else {
			fmt.Println(resp.Body)
		}
	}
}
