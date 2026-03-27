package schedules_repository

import "time"

type ScheduleModel struct {
	ID         string
	RoomID     string
	DaysOfWeek []int
	StartTime  time.Time
	EndTime    time.Time
}
