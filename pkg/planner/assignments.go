package planner

import (
	"math/rand"
)

func (h *Household) UpdateCurrentMember() {
	h.currentMember = h.Members[h.currentMemberIndex]
	h.currentMemberIndex++
	if h.currentMemberIndex >= len(h.Members) {
		h.currentMemberIndex = 0
	}
}

func (h *Household) AssignDailyTasks() {
	rand.Shuffle(len(h.DailyTasks), func(i, j int) {
		h.DailyTasks[i], h.DailyTasks[j] = h.DailyTasks[j], h.DailyTasks[i]
	})

	shuffledMembers := make([]*Member, len(h.Members))
	copy(shuffledMembers, (h.Members))
	rand.Shuffle(len(shuffledMembers), func(i, j int) {
		shuffledMembers[i], shuffledMembers[j] = shuffledMembers[j], shuffledMembers[i]
	})

	assigneeIndex := 0
	for _, task := range h.DailyTasks {
		task.SetAssignee(shuffledMembers[assigneeIndex])
		assigneeIndex++
		if assigneeIndex >= len(shuffledMembers) {
			assigneeIndex = 0
		}
	}
}

func (h *Household) AssignWeeklyTasks() {
	if h.remainingWeeklyTasks == 0 {
		h.remainingWeeklyTasks = len(h.WeeklyTasks)
	}

	amountAdded := 0
	weeklyTasksPerDay := max(len(h.WeeklyTasks)/len(h.Members), 1)
	for amountAdded < weeklyTasksPerDay && h.remainingWeeklyTasks > 0 {
		currentTaskIndex := max(len(h.WeeklyTasks) - h.remainingWeeklyTasks, 0)
		task := h.WeeklyTasks[currentTaskIndex]
		task.SetAssignee(h.currentMember)

		h.remainingWeeklyTasks--
		amountAdded++
	}
}

func (h *Household) AssignMonthlyTasks() {
	if h.remainingMonthlyTasks == 0 {
		h.remainingMonthlyTasks = len(h.MonthlyTasks)
	}

	randomMember := h.Members[rand.Intn(len(h.Members))]
	for randomMember.Name == h.currentMember.Name {
		randomMember = h.Members[rand.Intn(len(h.Members))]
	}

	taskIntervalMonth := 30 / len(h.MonthlyTasks)
	if h.dayOfTheMonth%taskIntervalMonth == 0 && h.remainingMonthlyTasks > 0 {
		currentTaskIndex := max(len(h.MonthlyTasks) - h.remainingMonthlyTasks, 0)
		task := h.MonthlyTasks[currentTaskIndex]
		task.SetAssignee(randomMember)

		h.remainingMonthlyTasks--
	}

	h.dayOfTheMonth++
	if h.dayOfTheMonth > 30 {
		h.dayOfTheMonth = 1
	}
}

func (h *Household) ClearAssignments() {
	for _, task := range h.DailyTasks {
		task.SetAssignee(nil)
	}
	for _, task := range h.WeeklyTasks {
		task.SetAssignee(nil)
	}
	for _, task := range h.MonthlyTasks {
		task.SetAssignee(nil)
	}
}

func (h *Household) GetAssignedTasks(member *Member) []Assignable {
	assignedTasks := []Assignable{}

	for _, task := range h.DailyTasks {
		assignee := task.GetAssignee()
		if assignee != nil && assignee.Name == member.Name {
			assignedTasks = append(assignedTasks, task)
		}
	}

	for _, task := range h.WeeklyTasks {
		assignee := task.GetAssignee()
		if assignee != nil && assignee.Name == member.Name {
			assignedTasks = append(assignedTasks, task)
		}
	}

	for _, task := range h.MonthlyTasks {
		assignee := task.GetAssignee()
		if assignee != nil && assignee.Name == member.Name {
			assignedTasks = append(assignedTasks, task)
		}
	}

	return assignedTasks
}
