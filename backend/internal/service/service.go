package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"dormitory/backend/config"
	"dormitory/backend/internal/dto"
	"dormitory/backend/internal/errs"
	"dormitory/backend/internal/middleware"
	"dormitory/backend/internal/model"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db  *sqlx.DB
	cfg config.Config
}

func New(db *sqlx.DB, cfg config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	var user model.User
	if err := s.db.GetContext(ctx, &user, `SELECT * FROM users WHERE username=$1`, req.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.LoginResponse{}, errs.ErrUnauthorized
		}
		return dto.LoginResponse{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.LoginResponse{}, errs.ErrUnauthorized
	}
	access, err := middleware.GenerateToken(s.cfg.JWT, user.ID, user.Role, "access", s.cfg.JWT.AccessExpiry)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	refresh, err := middleware.GenerateToken(s.cfg.JWT, user.ID, user.Role, "refresh", s.cfg.JWT.RefreshExpiry)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if err := s.storeRefreshToken(ctx, user.ID, refresh, time.Now().Add(s.cfg.JWT.RefreshExpiry)); err != nil {
		return dto.LoginResponse{}, err
	}
	return dto.LoginResponse{AccessToken: access, RefreshToken: refresh, TokenType: "Bearer", ExpiresIn: int64(s.cfg.JWT.AccessExpiry.Seconds())}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (dto.LoginResponse, error) {
	claims, err := middleware.ParseToken(s.cfg.JWT, refreshToken)
	if err != nil {
		return dto.LoginResponse{}, errs.ErrUnauthorized
	}
	if claims.Type != "refresh" {
		return dto.LoginResponse{}, errs.ErrUnauthorized
	}
	if err := s.requireValidRefreshToken(ctx, claims.UserID, refreshToken); err != nil {
		return dto.LoginResponse{}, err
	}
	var exists bool
	if err := s.db.GetContext(ctx, &exists, `SELECT EXISTS (SELECT 1 FROM users WHERE id=$1 AND role=$2)`, claims.UserID, claims.Role); err != nil {
		return dto.LoginResponse{}, err
	}
	if !exists {
		return dto.LoginResponse{}, errs.ErrUnauthorized
	}
	access, err := middleware.GenerateToken(s.cfg.JWT, claims.UserID, claims.Role, "access", s.cfg.JWT.AccessExpiry)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	refresh, err := middleware.GenerateToken(s.cfg.JWT, claims.UserID, claims.Role, "refresh", s.cfg.JWT.RefreshExpiry)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if err := s.rotateRefreshToken(ctx, claims.UserID, refreshToken, refresh, time.Now().Add(s.cfg.JWT.RefreshExpiry)); err != nil {
		return dto.LoginResponse{}, err
	}
	return dto.LoginResponse{AccessToken: access, RefreshToken: refresh, TokenType: "Bearer", ExpiresIn: int64(s.cfg.JWT.AccessExpiry.Seconds())}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	claims, err := middleware.ParseToken(s.cfg.JWT, refreshToken)
	if err != nil || claims.Type != "refresh" {
		return errs.ErrUnauthorized
	}
	res, err := s.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at=now() WHERE user_id=$1 AND token_hash=$2 AND revoked_at IS NULL`, claims.UserID, tokenHash(refreshToken))
	return requireChanged(res, err)
}

func (s *Service) CreateUser(ctx context.Context, req dto.CreateUserRequest) (model.User, error) {
	if !oneOf(req.Role, "student", "repair_staff", "cleaning_staff", "dormitory_manager", "system_admin") {
		return model.User{}, fmt.Errorf("%w: invalid role", errs.ErrBadRequest)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return model.User{}, err
	}
	var user model.User
	err = s.db.GetContext(ctx, &user, `
		INSERT INTO users (username, password_hash, role, name, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *`, req.Username, string(hash), req.Role, req.Name, req.Phone)
	return user, err
}

func (s *Service) User(ctx context.Context, id int64) (model.User, error) {
	var user model.User
	err := s.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id=$1`, id)
	return user, err
}

func (s *Service) ListBuildings(ctx context.Context) ([]model.Building, error) {
	var rows []model.Building
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM buildings ORDER BY id`)
	return rows, err
}

func (s *Service) RoomBalance(ctx context.Context, roomID int) (model.Room, error) {
	var room model.Room
	err := s.db.GetContext(ctx, &room, `SELECT * FROM rooms WHERE id=$1`, roomID)
	return room, err
}

func (s *Service) AvailableBeds(ctx context.Context, buildingID, floor *int) ([]model.AvailableBed, error) {
	query := `SELECT * FROM v_available_beds WHERE ($1::int IS NULL OR building_id=$1) AND ($2::int IS NULL OR floor=$2) ORDER BY building_id, room_number, bed_label`
	var rows []model.AvailableBed
	err := s.db.SelectContext(ctx, &rows, query, buildingID, floor)
	return rows, err
}

func (s *Service) SaveSurvey(ctx context.Context, studentID int64, req dto.SurveyRequest) (model.LifestyleSurvey, error) {
	var survey model.LifestyleSurvey
	err := s.db.GetContext(ctx, &survey, `
		INSERT INTO lifestyle_surveys (student_id, sleep_time, smoking, snoring, study_habit, remarks)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, student_id, sleep_time::text AS sleep_time, smoking, snoring, study_habit, remarks, submitted_at`,
		studentID, req.SleepTime, req.Smoking, req.Snoring, req.StudyHabit, req.Remarks)
	return survey, err
}

func (s *Service) LatestSurvey(ctx context.Context, studentID int64) (model.LifestyleSurvey, error) {
	var survey model.LifestyleSurvey
	err := s.db.GetContext(ctx, &survey, `
		SELECT id, student_id, sleep_time::text AS sleep_time, smoking, snoring, study_habit, remarks, submitted_at
		FROM lifestyle_surveys WHERE student_id=$1 ORDER BY submitted_at DESC LIMIT 1`, studentID)
	return survey, err
}

func (s *Service) MyRequests(ctx context.Context, studentID int64) ([]model.MyRequest, error) {
	var rows []model.MyRequest
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM v_my_requests WHERE student_id=$1 ORDER BY created_at DESC`, studentID)
	return rows, err
}

func (s *Service) Roommates(ctx context.Context, studentID int64) ([]model.Roommate, error) {
	var rows []model.Roommate
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM v_student_roommates WHERE student_id=$1 ORDER BY bed_label`, studentID)
	return rows, err
}

func (s *Service) CreateLeave(ctx context.Context, studentID int64, req dto.LeaveRequest) (int64, error) {
	typ := req.Type
	if typ == "" {
		typ = "normal"
	}
	if !oneOf(typ, "normal", "holiday") {
		return 0, fmt.Errorf("%w: invalid leave type", errs.ErrBadRequest)
	}
	return s.insertID(ctx, `INSERT INTO leave_applications (student_id, type, destination, emergency_contact, return_time, reason) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		studentID, typ, req.Destination, req.EmergencyContact, req.ReturnTime, req.Reason)
}

func (s *Service) CreateLateReturn(ctx context.Context, studentID int64, req dto.LateReturnRequest) (int64, error) {
	return s.insertID(ctx, `INSERT INTO late_return_records (student_id, return_date, reason) VALUES ($1,$2,$3) RETURNING id`, studentID, req.ReturnDate, req.Reason)
}

func (s *Service) CreateRepair(ctx context.Context, studentID int64, req dto.RepairRequest) (int64, error) {
	if err := s.requireStudentInRoom(ctx, studentID, req.RoomID); err != nil {
		return 0, err
	}
	return s.insertID(ctx, `INSERT INTO repair_requests (student_id, room_id, description) VALUES ($1,$2,$3) RETURNING id`, studentID, req.RoomID, req.Description)
}

func (s *Service) CreateCleaning(ctx context.Context, studentID int64, req dto.CleaningRequest) (int64, error) {
	return s.insertID(ctx, `INSERT INTO cleaning_requests (student_id, building_id, location_desc) VALUES ($1,$2,$3) RETURNING id`, studentID, req.BuildingID, req.LocationDesc)
}

func (s *Service) CreateOffCampus(ctx context.Context, studentID int64, req dto.OffCampusRequest) (int64, error) {
	if !oneOfInt(req.RetainBed, 0, 1) {
		return 0, fmt.Errorf("%w: retain_bed must be 0 or 1", errs.ErrBadRequest)
	}
	return s.insertID(ctx, `INSERT INTO off_campus_living_applications (student_id, retain_bed, reason, destination) VALUES ($1,$2,$3,$4) RETURNING id`,
		studentID, req.RetainBed, req.Reason, req.Destination)
}

func (s *Service) CreateRoomChange(ctx context.Context, studentID int64, req dto.RoomChangeRequest) (int64, error) {
	if (req.TargetRoomID == nil) != (req.TargetBedID == nil) {
		return 0, fmt.Errorf("%w: target_room_id and target_bed_id must be provided together", errs.ErrBadRequest)
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var from struct {
		BedID      int `db:"bed_id"`
		BuildingID int `db:"building_id"`
	}
	if err := tx.GetContext(ctx, &from, `
		SELECT bed.id AS bed_id, r.building_id
		FROM beds bed
		JOIN rooms r ON r.id = bed.room_id
		WHERE bed.student_id=$1 AND bed.status='occupied'
		FOR UPDATE OF bed`, studentID); err != nil {
		return 0, err
	}
	var recommendedBedID *int
	if req.TargetBedID == nil {
		var bedID int
		if err := tx.GetContext(ctx, &bedID, `
			SELECT bed.id
			FROM beds bed
			JOIN rooms r ON r.id = bed.room_id
			WHERE bed.status='available'
			  AND r.building_id=$1
			  AND NOT EXISTS (
			      SELECT 1 FROM room_change_requests rcr
			      WHERE rcr.status='pending'
			        AND COALESCE(rcr.target_bed_id, rcr.recommended_bed_id)=bed.id
			  )
			ORDER BY r.floor, r.room_number, bed.bed_label
			LIMIT 1
			FOR UPDATE OF bed SKIP LOCKED`, from.BuildingID); err != nil {
			return 0, err
		}
		recommendedBedID = &bedID
	}
	var id int64
	if err := tx.GetContext(ctx, &id, `
		INSERT INTO room_change_requests (student_id, from_bed_id, target_room_id, target_bed_id, recommended_bed_id, reason)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id`, studentID, from.BedID, req.TargetRoomID, req.TargetBedID, recommendedBedID, req.Reason); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (s *Service) CreatePayment(ctx context.Context, payerID int64, req dto.PaymentRequest) (int64, error) {
	if !oneOf(req.PaymentType, "water", "electricity", "both") {
		return 0, fmt.Errorf("%w: invalid payment_type", errs.ErrBadRequest)
	}
	if err := s.requireStudentInRoom(ctx, payerID, req.RoomID); err != nil {
		return 0, err
	}
	return s.insertID(ctx, `INSERT INTO utility_payments (room_id, payer_id, amount, payment_type) VALUES ($1,$2,$3,$4) RETURNING id`,
		req.RoomID, payerID, req.Amount, req.PaymentType)
}

func (s *Service) CreateAllocation(ctx context.Context, studentID int64) (int64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var hasOccupiedBed bool
	if err := tx.GetContext(ctx, &hasOccupiedBed, `SELECT EXISTS (SELECT 1 FROM beds WHERE student_id=$1 AND status='occupied')`, studentID); err != nil {
		return 0, err
	}
	if hasOccupiedBed {
		return 0, fmt.Errorf("%w: student already has an occupied bed", errs.ErrConflict)
	}
	var hasPending bool
	if err := tx.GetContext(ctx, &hasPending, `SELECT EXISTS (SELECT 1 FROM allocation_requests WHERE student_id=$1 AND status='pending')`, studentID); err != nil {
		return 0, err
	}
	if hasPending {
		return 0, fmt.Errorf("%w: student already has a pending allocation", errs.ErrConflict)
	}

	var bed struct {
		BedID  int `db:"bed_id"`
		RoomID int `db:"room_id"`
	}
	if err := tx.GetContext(ctx, &bed, `
		SELECT bed.id AS bed_id, bed.room_id
		FROM beds bed
		JOIN rooms r ON r.id = bed.room_id
		WHERE bed.status='available'
		  AND NOT EXISTS (
		      SELECT 1 FROM allocation_requests ar
		      WHERE ar.recommended_bed_id = bed.id
		        AND ar.status = 'pending'
		  )
		ORDER BY r.building_id, r.floor, r.room_number, bed.bed_label
		LIMIT 1
		FOR UPDATE OF bed SKIP LOCKED`); err != nil {
		return 0, err
	}
	var id int64
	if err := tx.GetContext(ctx, &id, `
		INSERT INTO allocation_requests (student_id, recommended_room_id, recommended_bed_id)
		VALUES ($1,$2,$3)
		RETURNING id`, studentID, bed.RoomID, bed.BedID); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (s *Service) insertID(ctx context.Context, query string, args ...any) (int64, error) {
	var id int64
	err := s.db.GetContext(ctx, &id, query, args...)
	return id, err
}

func (s *Service) PendingRepairs(ctx context.Context) ([]model.PendingRepair, error) {
	var rows []model.PendingRepair
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM v_pending_repairs ORDER BY created_at`)
	return rows, err
}

func (s *Service) AcceptRepair(ctx context.Context, id, staffID int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE repair_requests SET status='accepted', repair_staff_id=$2 WHERE id=$1 AND status='pending'`, id, staffID)
	return requireChanged(res, err)
}

func (s *Service) CompleteRepair(ctx context.Context, id, staffID int64, desc string) error {
	res, err := s.db.ExecContext(ctx, `UPDATE repair_requests SET status='repaired', repair_description=$3 WHERE id=$1 AND repair_staff_id=$2 AND status='accepted'`, id, staffID, desc)
	return requireChanged(res, err)
}

func (s *Service) ReviewRepair(ctx context.Context, id, reviewerID int64, req dto.ReviewRequest) error {
	if req.Status != "completed" && req.Status != "rejected" {
		return fmt.Errorf("%w: status must be completed or rejected", errs.ErrBadRequest)
	}
	res, err := s.db.ExecContext(ctx, `UPDATE repair_requests SET status=$3, reviewer_id=$2, review_comment=$4 WHERE id=$1 AND status='repaired'`, id, reviewerID, req.Status, req.Comment)
	return requireChanged(res, err)
}

func (s *Service) PendingCleanings(ctx context.Context) ([]model.PendingCleaning, error) {
	var rows []model.PendingCleaning
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM v_pending_cleanings ORDER BY created_at`)
	return rows, err
}

func (s *Service) PendingAllocations(ctx context.Context) ([]model.AllocationRequest, error) {
	var rows []model.AllocationRequest
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM allocation_requests WHERE status='pending' ORDER BY created_at`)
	return rows, err
}

func (s *Service) AcceptCleaning(ctx context.Context, id, cleanerID int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE cleaning_requests SET status='accepted', cleaner_id=$2 WHERE id=$1 AND status='pending'`, id, cleanerID)
	return requireChanged(res, err)
}

func (s *Service) CompleteCleaning(ctx context.Context, id, cleanerID int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE cleaning_requests SET status='cleaned' WHERE id=$1 AND cleaner_id=$2 AND status='accepted'`, id, cleanerID)
	return requireChanged(res, err)
}

func (s *Service) ReviewCleaning(ctx context.Context, id, reviewerID int64, req dto.ReviewRequest) error {
	if req.Status != "completed" && req.Status != "rejected" {
		return fmt.Errorf("%w: status must be completed or rejected", errs.ErrBadRequest)
	}
	res, err := s.db.ExecContext(ctx, `UPDATE cleaning_requests SET status=$3, reviewer_id=$2, review_comment=$4 WHERE id=$1 AND status='cleaned'`, id, reviewerID, req.Status, req.Comment)
	return requireChanged(res, err)
}

func (s *Service) ReviewApplication(ctx context.Context, table string, id, reviewerID int64, status string) error {
	if status != "approved" && status != "rejected" {
		return fmt.Errorf("%w: status must be approved or rejected", errs.ErrBadRequest)
	}
	managerColumn := map[string]string{
		"leave_applications":             "manager_id",
		"late_return_records":            "manager_id",
		"off_campus_living_applications": "manager_id",
		"allocation_requests":            "admin_id",
	}[table]
	if managerColumn == "" {
		return errs.ErrBadRequest
	}
	res, err := s.db.ExecContext(ctx, fmt.Sprintf(`UPDATE %s SET status=$2, %s=$3 WHERE id=$1 AND status='pending'`, table, managerColumn), id, status, reviewerID)
	return requireChanged(res, err)
}

func (s *Service) ApproveAllocation(ctx context.Context, id, adminID int64, status string) error {
	if status == "rejected" {
		return s.ReviewApplication(ctx, "allocation_requests", id, adminID, status)
	}
	if status != "approved" {
		return errs.ErrBadRequest
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var req struct {
		StudentID int64 `db:"student_id"`
		BedID     int   `db:"recommended_bed_id"`
	}
	if err := tx.GetContext(ctx, &req, `SELECT student_id, recommended_bed_id FROM allocation_requests WHERE id=$1 AND status='pending' FOR UPDATE`, id); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `UPDATE beds SET status='occupied', student_id=$2, occupied_since=CURRENT_DATE WHERE id=$1 AND status='available'`, req.BedID, req.StudentID)
	if err := requireChanged(res, err); err != nil {
		return fmt.Errorf("%w: recommended bed is no longer available", err)
	}
	res, err = tx.ExecContext(ctx, `UPDATE allocation_requests SET status='approved', admin_id=$2 WHERE id=$1`, id, adminID)
	if err := requireChanged(res, err); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) ApproveRoomChange(ctx context.Context, id, managerID int64, status string) error {
	if status == "rejected" {
		return s.ReviewApplication(ctx, "room_change_requests", id, managerID, status)
	}
	if status != "approved" {
		return errs.ErrBadRequest
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var req struct {
		StudentID        int64 `db:"student_id"`
		FromBedID        int   `db:"from_bed_id"`
		TargetBedID      *int  `db:"target_bed_id"`
		RecommendedBedID *int  `db:"recommended_bed_id"`
	}
	if err := tx.GetContext(ctx, &req, `SELECT student_id, from_bed_id, target_bed_id, recommended_bed_id FROM room_change_requests WHERE id=$1 AND status='pending' FOR UPDATE`, id); err != nil {
		return err
	}
	newBedID := req.TargetBedID
	if newBedID == nil {
		newBedID = req.RecommendedBedID
	}
	if newBedID == nil {
		return fmt.Errorf("%w: no target or recommended bed", errs.ErrBadRequest)
	}
	res, err := tx.ExecContext(ctx, `UPDATE beds SET status='available', student_id=NULL, occupied_since=NULL WHERE id=$1 AND student_id=$2`, req.FromBedID, req.StudentID)
	if err := requireChanged(res, err); err != nil {
		return fmt.Errorf("%w: original bed is not occupied by applicant", err)
	}
	res, err = tx.ExecContext(ctx, `UPDATE beds SET status='occupied', student_id=$2, occupied_since=CURRENT_DATE WHERE id=$1 AND status='available'`, *newBedID, req.StudentID)
	if err := requireChanged(res, err); err != nil {
		return fmt.Errorf("%w: target bed is no longer available", err)
	}
	res, err = tx.ExecContext(ctx, `UPDATE room_change_requests SET status='approved', manager_id=$2 WHERE id=$1`, id, managerID)
	if err := requireChanged(res, err); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) ReviewOffCampus(ctx context.Context, id, managerID int64, status string) error {
	if status == "rejected" {
		return s.ReviewApplication(ctx, "off_campus_living_applications", id, managerID, status)
	}
	if status != "approved" {
		return errs.ErrBadRequest
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var req struct {
		StudentID int64 `db:"student_id"`
		RetainBed int   `db:"retain_bed"`
	}
	if err := tx.GetContext(ctx, &req, `SELECT student_id, retain_bed FROM off_campus_living_applications WHERE id=$1 AND status='pending' FOR UPDATE`, id); err != nil {
		return err
	}
	if req.RetainBed == 0 {
		res, err := tx.ExecContext(ctx, `UPDATE beds SET status='available', student_id=NULL, occupied_since=NULL WHERE student_id=$1`, req.StudentID)
		if err := requireChanged(res, err); err != nil {
			return fmt.Errorf("%w: applicant has no occupied bed to release", err)
		}
	}
	res, err := tx.ExecContext(ctx, `UPDATE off_campus_living_applications SET status='approved', manager_id=$2 WHERE id=$1`, id, managerID)
	if err := requireChanged(res, err); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) Notifications(ctx context.Context, userID int64) ([]model.Notification, error) {
	var rows []model.Notification
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM notifications WHERE recipient_id=$1 ORDER BY created_at DESC`, userID)
	return rows, err
}

func (s *Service) MarkNotificationRead(ctx context.Context, userID, id int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE notifications SET is_read=1 WHERE id=$1 AND recipient_id=$2`, id, userID)
	return requireChanged(res, err)
}

func (s *Service) DashboardSummary(ctx context.Context) ([]model.DormitorySummary, error) {
	var rows []model.DormitorySummary
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM v_dormitory_summary ORDER BY building_id`)
	return rows, err
}

func (s *Service) LowBalanceRooms(ctx context.Context) ([]model.LowBalanceRoom, error) {
	var rows []model.LowBalanceRoom
	err := s.db.SelectContext(ctx, &rows, `SELECT * FROM v_low_balance_rooms ORDER BY building_id, room_number`)
	return rows, err
}

func (s *Service) SaveAttachment(ctx context.Context, actorID int64, role string, req dto.AttachmentUploadRequest, fileName, contentType string, data []byte) (model.AttachmentMeta, error) {
	if err := s.requireAttachmentWrite(ctx, actorID, role, req.OwnerType, req.OwnerID, req.Category); err != nil {
		return model.AttachmentMeta{}, err
	}
	var meta model.AttachmentMeta
	err := s.db.GetContext(ctx, &meta, `
		INSERT INTO attachments (owner_type, owner_id, category, sort_order, file_name, content_type, file_data)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (owner_id) WHERE owner_type='user_avatar' AND category='avatar'
		DO UPDATE SET file_name=EXCLUDED.file_name, content_type=EXCLUDED.content_type, file_data=EXCLUDED.file_data, uploaded_at=now()
		RETURNING id, owner_type, owner_id, category, sort_order, content_type, file_name, uploaded_at`,
		req.OwnerType, req.OwnerID, req.Category, req.SortOrder, fileName, contentType, data)
	return meta, err
}

func (s *Service) Attachment(ctx context.Context, actorID int64, role string, id int64) (model.AttachmentData, error) {
	var meta model.AttachmentMeta
	if err := s.db.GetContext(ctx, &meta, `SELECT id, owner_type, owner_id, category, sort_order, content_type, file_name, uploaded_at FROM attachments WHERE id=$1`, id); err != nil {
		return model.AttachmentData{}, err
	}
	if err := s.requireAttachmentRead(ctx, actorID, role, meta.OwnerType, meta.OwnerID); err != nil {
		return model.AttachmentData{}, err
	}
	var data model.AttachmentData
	err := s.db.GetContext(ctx, &data, `SELECT id, content_type, file_name, file_data FROM attachments WHERE id=$1`, id)
	return data, err
}

func (s *Service) AttachmentMetadata(ctx context.Context, actorID int64, role, ownerType string, ownerID int64, category *string) ([]model.AttachmentMeta, error) {
	if err := s.requireAttachmentRead(ctx, actorID, role, ownerType, ownerID); err != nil {
		return nil, err
	}
	var rows []model.AttachmentMeta
	err := s.db.SelectContext(ctx, &rows, `
		SELECT id, owner_type, owner_id, category, sort_order, content_type, file_name, uploaded_at
		FROM v_attachment_metadata
		WHERE owner_type=$1 AND owner_id=$2 AND ($3::text IS NULL OR category=$3)
		ORDER BY sort_order, id`, ownerType, ownerID, category)
	return rows, err
}

func (s *Service) requireStudentInRoom(ctx context.Context, studentID int64, roomID int) error {
	var ok bool
	if err := s.db.GetContext(ctx, &ok, `SELECT EXISTS (SELECT 1 FROM beds WHERE room_id=$1 AND student_id=$2 AND status='occupied')`, roomID, studentID); err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: student is not assigned to this room", errs.ErrForbidden)
	}
	return nil
}

func (s *Service) requireAttachmentWrite(ctx context.Context, actorID int64, role, ownerType string, ownerID int64, category string) error {
	if role == "system_admin" || role == "dormitory_manager" {
		return nil
	}
	switch ownerType {
	case "user_avatar":
		if category == "avatar" && actorID == ownerID {
			return nil
		}
	case "repair":
		if category != "after" {
			break
		}
		var ok bool
		if err := s.db.GetContext(ctx, &ok, `SELECT EXISTS (SELECT 1 FROM repair_requests WHERE id=$1 AND repair_staff_id=$2)`, ownerID, actorID); err != nil {
			return err
		}
		if ok && role == "repair_staff" {
			return nil
		}
	case "cleaning":
		var ok bool
		if category == "before" && role == "student" {
			if err := s.db.GetContext(ctx, &ok, `SELECT EXISTS (SELECT 1 FROM cleaning_requests WHERE id=$1 AND student_id=$2)`, ownerID, actorID); err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
		if category == "after" && role == "cleaning_staff" {
			if err := s.db.GetContext(ctx, &ok, `SELECT EXISTS (SELECT 1 FROM cleaning_requests WHERE id=$1 AND cleaner_id=$2)`, ownerID, actorID); err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
	}
	return errs.ErrForbidden
}

func (s *Service) requireAttachmentRead(ctx context.Context, actorID int64, role, ownerType string, ownerID int64) error {
	if role == "system_admin" || role == "dormitory_manager" {
		return nil
	}
	switch ownerType {
	case "user_avatar":
		return nil
	case "repair":
		var ok bool
		if err := s.db.GetContext(ctx, &ok, `
			SELECT EXISTS (
				SELECT 1 FROM repair_requests rr
				WHERE rr.id=$1
				  AND (rr.student_id=$2 OR rr.repair_staff_id=$2)
			)`, ownerID, actorID); err != nil {
			return err
		}
		if ok {
			return nil
		}
	case "cleaning":
		var ok bool
		if err := s.db.GetContext(ctx, &ok, `
			SELECT EXISTS (
				SELECT 1 FROM cleaning_requests cr
				WHERE cr.id=$1
				  AND (cr.student_id=$2 OR cr.cleaner_id=$2)
			)`, ownerID, actorID); err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return errs.ErrForbidden
}

func requireChanged(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (s *Service) storeRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3)`, userID, tokenHash(token), expiresAt)
	return err
}

func (s *Service) requireValidRefreshToken(ctx context.Context, userID int64, token string) error {
	var ok bool
	if err := s.db.GetContext(ctx, &ok, `
		SELECT EXISTS (
			SELECT 1 FROM refresh_tokens
			WHERE user_id=$1
			  AND token_hash=$2
			  AND revoked_at IS NULL
			  AND expires_at > now()
		)`, userID, tokenHash(token)); err != nil {
		return err
	}
	if !ok {
		return errs.ErrUnauthorized
	}
	return nil
}

func (s *Service) rotateRefreshToken(ctx context.Context, userID int64, oldToken, newToken string, expiresAt time.Time) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at=now()
		WHERE user_id=$1 AND token_hash=$2 AND revoked_at IS NULL AND expires_at > now()`, userID, tokenHash(oldToken))
	if err := requireChanged(res, err); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3)`, userID, tokenHash(newToken), expiresAt); err != nil {
		return err
	}
	return tx.Commit()
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func oneOf(v string, allowed ...string) bool {
	for _, item := range allowed {
		if v == item {
			return true
		}
	}
	return false
}

func oneOfInt(v int, allowed ...int) bool {
	for _, item := range allowed {
		if v == item {
			return true
		}
	}
	return false
}
