package model

import "time"

type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	Name         string    `db:"name" json:"name"`
	Phone        *string   `db:"phone" json:"phone,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	HasSurvey    bool      `db:"has_survey" json:"has_survey"`
	HasBed       bool      `db:"has_bed" json:"has_bed"`
	BuildingID   *int      `db:"building_id" json:"building_id,omitempty"`
	BuildingName *string   `db:"building_name" json:"building_name,omitempty"`
	RoomID       *int      `db:"room_id" json:"room_id,omitempty"`
	RoomNumber   *string   `db:"room_number" json:"room_number,omitempty"`
	BedID        *int      `db:"bed_id" json:"bed_id,omitempty"`
	BedLabel     *string   `db:"bed_label" json:"bed_label,omitempty"`
}

type Building struct {
	ID       int     `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	Location *string `db:"location" json:"location,omitempty"`
}

type Room struct {
	ID                 int     `db:"id" json:"id"`
	BuildingID         int     `db:"building_id" json:"building_id"`
	RoomNumber         string  `db:"room_number" json:"room_number"`
	Floor              int     `db:"floor" json:"floor"`
	TotalBeds          int     `db:"total_beds" json:"total_beds"`
	WaterBalance       float64 `db:"water_balance" json:"water_balance"`
	ElectricityBalance float64 `db:"electricity_balance" json:"electricity_balance"`
}

type Bed struct {
	ID            int        `db:"id" json:"id"`
	RoomID        int        `db:"room_id" json:"room_id"`
	BedLabel      string     `db:"bed_label" json:"bed_label"`
	Status        string     `db:"status" json:"status"`
	StudentID     *int64     `db:"student_id" json:"student_id,omitempty"`
	OccupiedSince *time.Time `db:"occupied_since" json:"occupied_since,omitempty"`
}

type AvailableBed struct {
	BedID        int    `db:"bed_id" json:"bed_id"`
	RoomID       int    `db:"room_id" json:"room_id"`
	RoomNumber   string `db:"room_number" json:"room_number"`
	BedLabel     string `db:"bed_label" json:"bed_label"`
	BuildingID   int    `db:"building_id" json:"building_id"`
	BuildingName string `db:"building_name" json:"building_name"`
	Floor        int    `db:"floor" json:"floor"`
}

type DormitorySummary struct {
	BuildingID   int    `db:"building_id" json:"building_id"`
	BuildingName string `db:"building_name" json:"building_name"`
	TotalRooms   int    `db:"total_rooms" json:"total_rooms"`
	TotalBeds    int    `db:"total_beds" json:"total_beds"`
	OccupiedBeds int    `db:"occupied_beds" json:"occupied_beds"`
	FreeBeds     int    `db:"free_beds" json:"free_beds"`
}

type LowBalanceRoom struct {
	RoomID             int     `db:"room_id" json:"room_id"`
	BuildingID         int     `db:"building_id" json:"building_id"`
	RoomNumber         string  `db:"room_number" json:"room_number"`
	WaterBalance       float64 `db:"water_balance" json:"water_balance"`
	ElectricityBalance float64 `db:"electricity_balance" json:"electricity_balance"`
}

type Roommate struct {
	StudentID          int64   `db:"student_id" json:"student_id"`
	RoommateID         int64   `db:"roommate_id" json:"roommate_id"`
	RoommateName       string  `db:"roommate_name" json:"roommate_name"`
	RoommatePhone      *string `db:"roommate_phone" json:"roommate_phone,omitempty"`
	BedLabel           string  `db:"bed_label" json:"bed_label"`
	AvatarAttachmentID *int64  `db:"avatar_attachment_id" json:"avatar_attachment_id,omitempty"`
}

type LifestyleSurvey struct {
	ID          int64     `db:"id" json:"id"`
	StudentID   int64     `db:"student_id" json:"student_id"`
	SleepTime   *string   `db:"sleep_time" json:"sleep_time,omitempty"`
	Smoking     int       `db:"smoking" json:"smoking"`
	Snoring     int       `db:"snoring" json:"snoring"`
	StudyHabit  *string   `db:"study_habit" json:"study_habit,omitempty"`
	Remarks     *string   `db:"remarks" json:"remarks,omitempty"`
	SubmittedAt time.Time `db:"submitted_at" json:"submitted_at"`
}

type MyRequest struct {
	StudentID   int64     `db:"student_id" json:"student_id"`
	RequestType string    `db:"request_type" json:"request_type"`
	RequestID   int64     `db:"request_id" json:"request_id"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Detail      string    `db:"detail" json:"detail"`
}

type Notification struct {
	ID          int64     `db:"id" json:"id"`
	RecipientID int64     `db:"recipient_id" json:"recipient_id"`
	RoomID      *int      `db:"room_id" json:"room_id,omitempty"`
	Message     string    `db:"message" json:"message"`
	Type        string    `db:"type" json:"type"`
	IsRead      int       `db:"is_read" json:"is_read"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type AttachmentMeta struct {
	ID          int64     `db:"id" json:"id"`
	OwnerType   string    `db:"owner_type" json:"owner_type"`
	OwnerID     int64     `db:"owner_id" json:"owner_id"`
	Category    string    `db:"category" json:"category"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	ContentType string    `db:"content_type" json:"content_type"`
	FileName    *string   `db:"file_name" json:"file_name,omitempty"`
	UploadedAt  time.Time `db:"uploaded_at" json:"uploaded_at"`
}

type AttachmentData struct {
	ID          int64   `db:"id"`
	ContentType string  `db:"content_type"`
	FileName    *string `db:"file_name"`
	FileData    []byte  `db:"file_data"`
}

type PendingRepair struct {
	RequestID         int64      `db:"request_id" json:"request_id"`
	Status            string     `db:"status" json:"status"`
	StudentID         int64      `db:"student_id" json:"student_id"`
	StudentName       string     `db:"student_name" json:"student_name"`
	RoomID            int        `db:"room_id" json:"room_id"`
	RoomNumber        string     `db:"room_number" json:"room_number"`
	Description       string     `db:"description" json:"description"`
	RepairStaffID     *int64     `db:"repair_staff_id" json:"repair_staff_id,omitempty"`
	RepairStaffName   *string    `db:"repair_staff_name" json:"repair_staff_name,omitempty"`
	RepairDescription *string    `db:"repair_description" json:"repair_description,omitempty"`
	ReviewerID        *int64     `db:"reviewer_id" json:"reviewer_id,omitempty"`
	ReviewerName      *string    `db:"reviewer_name" json:"reviewer_name,omitempty"`
	ReviewComment     *string    `db:"review_comment" json:"review_comment,omitempty"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	AcceptedAt        *time.Time `db:"accepted_at" json:"accepted_at,omitempty"`
	RepairedAt        *time.Time `db:"repaired_at" json:"repaired_at,omitempty"`
	ReviewedAt        *time.Time `db:"reviewed_at" json:"reviewed_at,omitempty"`
}

type PendingCleaning struct {
	RequestID     int64      `db:"request_id" json:"request_id"`
	Status        string     `db:"status" json:"status"`
	StudentID     int64      `db:"student_id" json:"student_id"`
	StudentName   string     `db:"student_name" json:"student_name"`
	BuildingID    int        `db:"building_id" json:"building_id"`
	BuildingName  string     `db:"building_name" json:"building_name"`
	LocationDesc  string     `db:"location_desc" json:"location_desc"`
	CleanerID     *int64     `db:"cleaner_id" json:"cleaner_id,omitempty"`
	CleanerName   *string    `db:"cleaner_name" json:"cleaner_name,omitempty"`
	ReviewerID    *int64     `db:"reviewer_id" json:"reviewer_id,omitempty"`
	ReviewerName  *string    `db:"reviewer_name" json:"reviewer_name,omitempty"`
	ReviewComment *string    `db:"review_comment" json:"review_comment,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	AcceptedAt    *time.Time `db:"accepted_at" json:"accepted_at,omitempty"`
	CleanedAt     *time.Time `db:"cleaned_at" json:"cleaned_at,omitempty"`
	ReviewedAt    *time.Time `db:"reviewed_at" json:"reviewed_at,omitempty"`
}

type PendingLeave struct {
	ID               int64     `db:"id" json:"id"`
	StudentID        int64     `db:"student_id" json:"student_id"`
	StudentName      string    `db:"student_name" json:"student_name"`
	Type             string    `db:"type" json:"type"`
	Destination      string    `db:"destination" json:"destination"`
	EmergencyContact string    `db:"emergency_contact" json:"emergency_contact"`
	ReturnTime       time.Time `db:"return_time" json:"return_time"`
	Reason           string    `db:"reason" json:"reason"`
	Status           string    `db:"status" json:"status"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

type PendingLateReturn struct {
	ID          int64     `db:"id" json:"id"`
	StudentID   int64     `db:"student_id" json:"student_id"`
	StudentName string    `db:"student_name" json:"student_name"`
	ReturnDate  time.Time `db:"return_date" json:"return_date"`
	Reason      string    `db:"reason" json:"reason"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type PendingRoomChange struct {
	ID                      int64     `db:"id" json:"id"`
	StudentID               int64     `db:"student_id" json:"student_id"`
	StudentName             string    `db:"student_name" json:"student_name"`
	FromBedID               int       `db:"from_bed_id" json:"from_bed_id"`
	FromBuildingName        string    `db:"from_building_name" json:"from_building_name"`
	FromRoomNumber          string    `db:"from_room_number" json:"from_room_number"`
	FromBedLabel            string    `db:"from_bed_label" json:"from_bed_label"`
	TargetRoomID            *int      `db:"target_room_id" json:"target_room_id,omitempty"`
	TargetBedID             *int      `db:"target_bed_id" json:"target_bed_id,omitempty"`
	TargetBuildingName      *string   `db:"target_building_name" json:"target_building_name,omitempty"`
	TargetRoomNumber        *string   `db:"target_room_number" json:"target_room_number,omitempty"`
	TargetBedLabel          *string   `db:"target_bed_label" json:"target_bed_label,omitempty"`
	RecommendedBedID        *int      `db:"recommended_bed_id" json:"recommended_bed_id,omitempty"`
	RecommendedBuildingName *string   `db:"recommended_building_name" json:"recommended_building_name,omitempty"`
	RecommendedRoomNumber   *string   `db:"recommended_room_number" json:"recommended_room_number,omitempty"`
	RecommendedBedLabel     *string   `db:"recommended_bed_label" json:"recommended_bed_label,omitempty"`
	Reason                  string    `db:"reason" json:"reason"`
	Status                  string    `db:"status" json:"status"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
}

type PendingOffCampus struct {
	ID          int64     `db:"id" json:"id"`
	StudentID   int64     `db:"student_id" json:"student_id"`
	StudentName string    `db:"student_name" json:"student_name"`
	RetainBed   int       `db:"retain_bed" json:"retain_bed"`
	Reason      string    `db:"reason" json:"reason"`
	Destination *string   `db:"destination" json:"destination,omitempty"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type AllocationRequest struct {
	ID                int64      `db:"id" json:"id"`
	StudentID         int64      `db:"student_id" json:"student_id"`
	RecommendedRoomID int        `db:"recommended_room_id" json:"recommended_room_id"`
	RecommendedBedID  int        `db:"recommended_bed_id" json:"recommended_bed_id"`
	Status            string     `db:"status" json:"status"`
	AdminID           *int64     `db:"admin_id" json:"admin_id,omitempty"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	ResolvedAt        *time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
}
