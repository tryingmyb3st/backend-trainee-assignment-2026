package domain

import (
	"errors"
)

var (
	INVALID_REQUEST     = NewCustomError("INVALID_REQUEST", "invalid request")
	INTERNAL_ERROR      = NewCustomError("INTERNAL_ERROR", "internal server error")
	SCHEDULE_EXISTS     = NewCustomError("SCHEDULE_EXISTS", "schedule for this room already exists and cannot be changed")
	UNAUTHORIZED        = NewCustomError("UNAUTHORIZED", "wrong authorization")
	ROOM_NOT_FOUND      = NewCustomError("ROOM_NOT_FOUND", "room is not found")
	SLOT_NOT_FOUND      = NewCustomError("SLOT_NOT_FOUND", "slot is not found")
	BOOKING_NOT_FOUND   = NewCustomError("BOOKING_NOT_FOUND", "booking is not found")
	SLOT_ALREADY_BOOKED = NewCustomError("SLOT_ALREADY_BOOKED", "slot is already booked")
	FORBIDDEN           = NewCustomError("FORBIDDEN", "cannot cancel another user's booking")
	NOT_FOUND           = errors.New("not found")
)

type CustomError struct {
	Code    string `json:"code" enums:"INVALID_REQUEST, INTERNAL_ERROR, SCHEDULE_EXISTS, SLOT_ALREADY_BOOKED, FORBIDEN, NOT_FOUND"`
	Message string `json:"message" example:"invalid request"`
}

type InternalError struct {
	Code    string `json:"code" enums:"INTERNAL_ERROR"`
	Message string `json:"message" example:"internal server error"`
}

func NewCustomError(code string, msg string) CustomError {
	return CustomError{
		Code:    code,
		Message: msg,
	}
}

func (e CustomError) Error() string {
	return e.Message
}
