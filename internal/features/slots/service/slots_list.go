package slots_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const MAX_INTERVALS = 48

type Interval struct {
	TimeStart time.Time
	TimeEnd   time.Time
}

func (s *SlotsService) GetActiveSlots(ctx context.Context, roomId string, date time.Time) ([]domain.Slot, error) {

	schedule, err := s.scheduleRepository.GetScheduleByRoomId(ctx, roomId)
	if err != nil {
		if _, err := s.roomsRepository.GetRoom(ctx, roomId); err != nil && errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no schedule for roomID: %w", domain.NOT_FOUND)
		}
		return []domain.Slot{}, nil
	}

	var slots []domain.Slot
	slots, err = s.slotsRepository.GetActiveSlots(ctx, roomId, date.Weekday())
	if err != nil {
		return nil, fmt.Errorf("get slots: %w", err)
	}

	busy, err := s.slotsRepository.IsSlotsBusy(ctx, roomId, date.Weekday())
	if err != nil {
		return nil, fmt.Errorf("get busy slots: %w", err)
	}
	if busy {
		return slots, nil
	}

	if len(slots) == 0 {
		createdSlots := s.makeSlotsFromSchedule(*schedule, date)
		slots, err = s.slotsRepository.SaveSlots(ctx, createdSlots)
		if err != nil {
			return nil, fmt.Errorf("save slots: %w", err)
		}
	}

	return slots, nil
}

func (s *SlotsService) makeSlotsFromSchedule(schedule domain.Schedule, date time.Time) []domain.Slot {

	intervals := s.makeSlotsTimeIntervals(schedule)

	slots := make([]domain.Slot, 0)

	for _, weekday := range schedule.DaysOfWeek {

		if time.Weekday(weekday) != date.Weekday() {
			if !(weekday == 7 && date.Weekday() == time.Sunday) {
				continue
			}
		}

		for _, interval := range intervals {
			var slot domain.Slot
			slot.ID = uuid.New().String()
			slot.RoomID = schedule.RoomID

			dateToday := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
			dateOfSlot := getNextDateByWeekday(weekday, dateToday)

			slot.StartTime = time.Date(
				dateOfSlot.Year(),
				date.Month(),
				dateOfSlot.Day(),
				interval.TimeStart.Hour(),
				interval.TimeStart.Minute(),
				0,
				0,
				time.UTC,
			)
			slot.EndTime = time.Date(
				dateOfSlot.Year(),
				date.Month(),
				dateOfSlot.Day(),
				interval.TimeEnd.Hour(),
				interval.TimeEnd.Minute(),
				0,
				0,
				time.UTC,
			)
			slots = append(slots, slot)
		}
	}

	return slots
}

func (s *SlotsService) makeSlotsTimeIntervals(schedule domain.Schedule) []Interval {
	intervals := make([]Interval, 0, MAX_INTERVALS)

	st := schedule.StartTime

	for st.Add(30*time.Minute).Before(schedule.EndTime) || st.Add(30*time.Minute).Equal(schedule.EndTime) {
		intervals = append(intervals, Interval{
			TimeStart: st,
			TimeEnd:   st.Add(30 * time.Minute),
		})

		st = st.Add(30 * time.Minute)
	}

	return intervals
}

func getNextDateByWeekday(weekday int, dateNow time.Time) time.Time {
	var desiredWeekday time.Weekday
	if weekday == 7 {
		weekday = 0
	}
	desiredWeekday = time.Weekday(weekday)

	daysUntil := int((desiredWeekday-dateNow.Weekday())+7) % 7
	return dateNow.AddDate(0, 0, daysUntil)
}
