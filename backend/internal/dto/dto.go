package dto

import "time"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required,max=32"`
	Password string  `json:"password" binding:"required,min=6"`
	Role     string  `json:"role" binding:"required"`
	Name     string  `json:"name" binding:"required,max=32"`
	Phone    *string `json:"phone"`
}

type SurveyRequest struct {
	SleepTime  *string `json:"sleep_time"`
	Smoking    int     `json:"smoking"`
	Snoring    int     `json:"snoring"`
	StudyHabit *string `json:"study_habit"`
	Remarks    *string `json:"remarks"`
}

type LeaveRequest struct {
	Type             string    `json:"type"`
	Destination      string    `json:"destination" binding:"required"`
	EmergencyContact string    `json:"emergency_contact" binding:"required"`
	ReturnTime       time.Time `json:"return_time" binding:"required"`
	Reason           string    `json:"reason" binding:"required"`
}

type LateReturnRequest struct {
	ReturnDate string `json:"return_date" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
}

type RoomChangeRequest struct {
	TargetRoomID *int   `json:"target_room_id"`
	TargetBedID  *int   `json:"target_bed_id"`
	Reason       string `json:"reason" binding:"required"`
}

type OffCampusRequest struct {
	RetainBed   int     `json:"retain_bed"`
	Reason      string  `json:"reason" binding:"required"`
	Destination *string `json:"destination"`
}

type RepairRequest struct {
	RoomID      int    `json:"room_id" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type CleaningRequest struct {
	BuildingID   int    `json:"building_id" binding:"required"`
	LocationDesc string `json:"location_desc" binding:"required"`
}

type PaymentRequest struct {
	RoomID      int     `json:"room_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	PaymentType string  `json:"payment_type" binding:"required"`
}

type ReviewRequest struct {
	Status  string  `json:"status" binding:"required"`
	Comment *string `json:"comment"`
}

type RepairCompleteRequest struct {
	RepairDescription string `json:"repair_description" binding:"required"`
}

type AttachmentUploadRequest struct {
	OwnerType string `form:"owner_type" binding:"required"`
	OwnerID   int64  `form:"owner_id" binding:"required"`
	Category  string `form:"category" binding:"required"`
	SortOrder int    `form:"sort_order"`
}
