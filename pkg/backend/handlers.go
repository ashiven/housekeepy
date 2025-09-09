package backend

import (
	"encoding/json"
	"fmt"
	"household-planner/pkg/planner"
	"net/http"
	"sync"
)

var (
	household     *planner.Household
	fileLock      sync.Mutex
	adminPassword = planner.GetEnvVar("ADMIN_PASSWORD")
)

func SetHousehold(householdToSet *planner.Household) {
	household = householdToSet
}

func checkAdminPassword(w http.ResponseWriter, r *http.Request) {
	type PasswordRequest struct {
		Password string `json:"password"`
	}
	var passwordRequest PasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&passwordRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if passwordRequest.Password != adminPassword {
		http.Error(w, "Password is incorrect", http.StatusUnauthorized)
		return
	}

	fmt.Println("[INFO] Admin password is correct, setting cookie")
	handleSetCookie(w)
	fmt.Println("[INFO] Cookie set successfully")
	w.WriteHeader(http.StatusOK)
}

func handleUpdate[T any](w http.ResponseWriter, r *http.Request, section string, setConfigFile func(option *T), setConfigMem func(updated []*T)) {
	fmt.Println("[INFO] Getting cookie for update operation")
	err := handleGetCookie(r)
	if err != nil {
		fmt.Println("[ERROR] Failed to read cookie:", err)
		http.Error(w, fmt.Sprintf("Failed to read cookie: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("[INFO] Cookie retrieved successfully")

	fileLock.Lock()
	defer fileLock.Unlock()

	var updatedOptions []*T
	if err := json.NewDecoder(r.Body).Decode(&updatedOptions); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	setConfigMem(updatedOptions)

	household.Config.RemoveSection(section)
	household.Config.AddSection(section)
	for _, option := range updatedOptions {
		setConfigFile(option)
	}
	if err := household.Config.SaveWithDelimiter(household.Configfile, ":"); err != nil {
		http.Error(w, "Error saving config: %v", http.StatusInternalServerError)
		return
	}
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(household.Members)
}

func updateMembers(w http.ResponseWriter, r *http.Request) {
	handleUpdate(w, r, "Members", func(member *planner.Member) {
		household.Config.Set("Members", member.Name, member.PhoneNumber)
	}, func(updatedMembers []*planner.Member) {
		household.Members = updatedMembers
	})
}

func getDailyTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(household.DailyTasks)
}

func updateDailyTasks(w http.ResponseWriter, r *http.Request) {
	handleUpdate(w, r, "Daily Tasks", func(task *planner.DailyTask) {
		household.Config.Set("Daily Tasks", task.Name, "")
	}, func(updatedTasks []*planner.DailyTask) {
		household.DailyTasks = updatedTasks
	})
}

func getWeeklyTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(household.WeeklyTasks)
}

func updateWeeklyTasks(w http.ResponseWriter, r *http.Request) {
	handleUpdate(w, r, "Weekly Tasks", func(task *planner.WeeklyTask) {
		household.Config.Set("Weekly Tasks", task.Name, "")
	}, func(updatedTasks []*planner.WeeklyTask) {
		household.WeeklyTasks = updatedTasks
	})
}

func getMonthlyTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(household.MonthlyTasks)
}

func updateMonthlyTasks(w http.ResponseWriter, r *http.Request) {
	handleUpdate(w, r, "Monthly Tasks", func(task *planner.MonthlyTask) {
		household.Config.Set("Monthly Tasks", task.Name, "")
	}, func(updatedTasks []*planner.MonthlyTask) {
		household.MonthlyTasks = updatedTasks
	})
}
