package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"dormitory/backend/config"
	"dormitory/backend/internal/dto"
	"dormitory/backend/internal/errs"
	"dormitory/backend/internal/middleware"
	"dormitory/backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc *service.Service
	cfg config.Config
	log *zap.Logger
}

func NewRouter(svc *service.Service, cfg config.Config, log *zap.Logger) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	h := &Handler{svc: svc, cfg: cfg, log: log}
	r := gin.New()
	r.Use(gin.Recovery(), cors.Default())
	r.Use(func(c *gin.Context) {
		c.Next()
		log.Info("http_request", zap.String("method", c.Request.Method), zap.String("path", c.FullPath()), zap.Int("status", c.Writer.Status()))
	})

	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) { ok(c, gin.H{"status": "ok"}) })
	api.POST("/auth/login", h.login)
	api.POST("/auth/refresh", h.refresh)
	api.POST("/auth/logout", h.logout)

	auth := api.Group("", middleware.AuthRequired(cfg.JWT))
	auth.GET("/students/me", h.me)
	auth.GET("/students/me/survey", h.latestSurvey)
	auth.POST("/students/me/survey", h.saveSurvey)
	auth.GET("/students/me/requests", h.myRequests)
	auth.GET("/students/me/roommates", h.roommates)
	auth.GET("/buildings", h.buildings)
	auth.GET("/rooms/:id/balance", h.roomBalance)
	auth.GET("/beds/available", h.availableBeds)
	auth.GET("/notifications", h.notifications)
	auth.PUT("/notifications/:id/read", h.markNotificationRead)
	auth.POST("/attachments", h.uploadAttachment)
	auth.GET("/attachments", h.attachmentMetadata)
	auth.GET("/attachments/:id", h.downloadAttachment)

	student := api.Group("", middleware.AuthRequired(cfg.JWT, "student"))
	student.POST("/leaves", h.createLeave)
	student.POST("/allocations", h.createAllocation)
	student.POST("/late-returns", h.createLateReturn)
	student.POST("/room-changes", h.createRoomChange)
	student.POST("/off-campus", h.createOffCampus)
	student.POST("/repairs", h.createRepair)
	student.POST("/cleanings", h.createCleaning)
	student.POST("/payments", h.createPayment)

	repairRead := api.Group("", middleware.AuthRequired(cfg.JWT, "repair_staff", "dormitory_manager"))
	repairRead.GET("/repairs/pending", h.pendingRepairs)

	repair := api.Group("", middleware.AuthRequired(cfg.JWT, "repair_staff"))
	repair.PUT("/repairs/:id/accept", h.acceptRepair)
	repair.PUT("/repairs/:id/repair", h.completeRepair)

	cleaningRead := api.Group("", middleware.AuthRequired(cfg.JWT, "cleaning_staff", "dormitory_manager"))
	cleaningRead.GET("/cleanings/pending", h.pendingCleanings)

	cleaning := api.Group("", middleware.AuthRequired(cfg.JWT, "cleaning_staff"))
	cleaning.PUT("/cleanings/:id/accept", h.acceptCleaning)
	cleaning.PUT("/cleanings/:id/clean", h.completeCleaning)

	manager := api.Group("", middleware.AuthRequired(cfg.JWT, "dormitory_manager"))
	manager.GET("/leaves/pending", h.pendingLeaves)
	manager.PUT("/leaves/:id/review", h.reviewLeave)
	manager.GET("/late-returns/pending", h.pendingLateReturns)
	manager.PUT("/late-returns/:id/review", h.reviewLateReturn)
	manager.GET("/room-changes/pending", h.pendingRoomChanges)
	manager.PUT("/room-changes/:id/review", h.reviewRoomChange)
	manager.GET("/off-campus/pending", h.pendingOffCampus)
	manager.PUT("/off-campus/:id/review", h.reviewOffCampus)
	manager.PUT("/repairs/:id/review", h.reviewRepair)
	manager.PUT("/cleanings/:id/review", h.reviewCleaning)
	manager.GET("/dashboard/summary", h.dashboardSummary)
	manager.GET("/dashboard/low-balance", h.lowBalanceRooms)

	admin := api.Group("", middleware.AuthRequired(cfg.JWT, "system_admin"))
	admin.POST("/users", h.createUser)
	admin.GET("/allocations/pending", h.pendingAllocations)
	admin.PUT("/allocations/:id/review", h.reviewAllocation)
	admin.GET("/dashboard/summary", h.dashboardSummary)
	admin.GET("/dashboard/low-balance", h.lowBalanceRooms)

	return r
}

func (h *Handler) login(c *gin.Context) {
	req, ok := bindJSON[dto.LoginRequest](c)
	if !ok {
		return
	}
	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		fail(c, err)
		return
	}
	okJSON(c, resp)
}

func (h *Handler) refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err)
		return
	}
	resp, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		fail(c, err)
		return
	}
	okJSON(c, resp)
}

func (h *Handler) logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err)
		return
	}
	respondNoContent(c, h.svc.Logout(c.Request.Context(), req.RefreshToken))
}

func (h *Handler) me(c *gin.Context) {
	user, err := h.svc.User(c.Request.Context(), middleware.UserID(c))
	respond(c, user, err)
}

func (h *Handler) createUser(c *gin.Context) {
	req, bound := bindJSON[dto.CreateUserRequest](c)
	if !bound {
		return
	}
	user, err := h.svc.CreateUser(c.Request.Context(), req)
	respondCreated(c, user, err)
}

func (h *Handler) buildings(c *gin.Context) {
	rows, err := h.svc.ListBuildings(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) roomBalance(c *gin.Context) {
	id, ok := paramInt(c, "id")
	if !ok {
		return
	}
	room, err := h.svc.RoomBalance(c.Request.Context(), id)
	respond(c, room, err)
}

func (h *Handler) availableBeds(c *gin.Context) {
	buildingID := queryIntPtr(c, "building_id")
	floor := queryIntPtr(c, "floor")
	rows, err := h.svc.AvailableBeds(c.Request.Context(), buildingID, floor)
	respond(c, rows, err)
}

func (h *Handler) saveSurvey(c *gin.Context) {
	req, bound := bindJSON[dto.SurveyRequest](c)
	if !bound {
		return
	}
	survey, err := h.svc.SaveSurvey(c.Request.Context(), middleware.UserID(c), req)
	respondCreated(c, survey, err)
}

func (h *Handler) latestSurvey(c *gin.Context) {
	survey, err := h.svc.LatestSurvey(c.Request.Context(), middleware.UserID(c))
	respond(c, survey, err)
}

func (h *Handler) myRequests(c *gin.Context) {
	rows, err := h.svc.MyRequests(c.Request.Context(), middleware.UserID(c))
	respond(c, rows, err)
}

func (h *Handler) roommates(c *gin.Context) {
	rows, err := h.svc.Roommates(c.Request.Context(), middleware.UserID(c))
	respond(c, rows, err)
}

func (h *Handler) createLeave(c *gin.Context) {
	req, bound := bindJSON[dto.LeaveRequest](c)
	if !bound {
		return
	}
	id, err := h.svc.CreateLeave(c.Request.Context(), middleware.UserID(c), req)
	respondID(c, id, err)
}

func (h *Handler) createAllocation(c *gin.Context) {
	id, err := h.svc.CreateAllocation(c.Request.Context(), middleware.UserID(c))
	respondID(c, id, err)
}

func (h *Handler) createLateReturn(c *gin.Context) {
	req, bound := bindJSON[dto.LateReturnRequest](c)
	if !bound {
		return
	}
	id, err := h.svc.CreateLateReturn(c.Request.Context(), middleware.UserID(c), req)
	respondID(c, id, err)
}

func (h *Handler) createRoomChange(c *gin.Context) {
	req, bound := bindJSON[dto.RoomChangeRequest](c)
	if !bound {
		return
	}
	id, err := h.svc.CreateRoomChange(c.Request.Context(), middleware.UserID(c), req)
	respondID(c, id, err)
}

func (h *Handler) createOffCampus(c *gin.Context) {
	req, bound := bindJSON[dto.OffCampusRequest](c)
	if !bound {
		return
	}
	id, err := h.svc.CreateOffCampus(c.Request.Context(), middleware.UserID(c), req)
	respondID(c, id, err)
}

func (h *Handler) createRepair(c *gin.Context) {
	req, bound := bindJSON[dto.RepairRequest](c)
	if !bound {
		return
	}
	id, err := h.svc.CreateRepair(c.Request.Context(), middleware.UserID(c), req)
	respondID(c, id, err)
}

func (h *Handler) createCleaning(c *gin.Context) {
	req, bound := bindJSON[dto.CleaningRequest](c)
	if !bound {
		return
	}
	id, err := h.svc.CreateCleaning(c.Request.Context(), middleware.UserID(c), req)
	respondID(c, id, err)
}

func (h *Handler) createPayment(c *gin.Context) {
	req, bound := bindJSON[dto.PaymentRequest](c)
	if !bound {
		return
	}
	id, err := h.svc.CreatePayment(c.Request.Context(), middleware.UserID(c), req)
	respondID(c, id, err)
}

func (h *Handler) pendingRepairs(c *gin.Context) {
	rows, err := h.svc.PendingRepairs(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) acceptRepair(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	respondNoContent(c, h.svc.AcceptRepair(c.Request.Context(), id, middleware.UserID(c)))
}

func (h *Handler) completeRepair(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	req, bound := bindJSON[dto.RepairCompleteRequest](c)
	if !bound {
		return
	}
	respondNoContent(c, h.svc.CompleteRepair(c.Request.Context(), id, middleware.UserID(c), req.RepairDescription))
}

func (h *Handler) pendingCleanings(c *gin.Context) {
	rows, err := h.svc.PendingCleanings(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) acceptCleaning(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	respondNoContent(c, h.svc.AcceptCleaning(c.Request.Context(), id, middleware.UserID(c)))
}

func (h *Handler) completeCleaning(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	respondNoContent(c, h.svc.CompleteCleaning(c.Request.Context(), id, middleware.UserID(c)))
}

func (h *Handler) reviewLeave(c *gin.Context) {
	h.reviewApplication(c, "leave_applications")
}

func (h *Handler) pendingLeaves(c *gin.Context) {
	rows, err := h.svc.PendingLeaves(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) reviewLateReturn(c *gin.Context) {
	h.reviewApplication(c, "late_return_records")
}

func (h *Handler) pendingLateReturns(c *gin.Context) {
	rows, err := h.svc.PendingLateReturns(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) reviewApplication(c *gin.Context, table string) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	req, bound := bindJSON[dto.ReviewRequest](c)
	if !bound {
		return
	}
	respondNoContent(c, h.svc.ReviewApplication(c.Request.Context(), table, id, middleware.UserID(c), req.Status))
}

func (h *Handler) reviewAllocation(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	req, bound := bindJSON[dto.ReviewRequest](c)
	if !bound {
		return
	}
	respondNoContent(c, h.svc.ApproveAllocation(c.Request.Context(), id, middleware.UserID(c), req.Status))
}

func (h *Handler) reviewRoomChange(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	req, bound := bindJSON[dto.ReviewRequest](c)
	if !bound {
		return
	}
	respondNoContent(c, h.svc.ApproveRoomChange(c.Request.Context(), id, middleware.UserID(c), req.Status))
}

func (h *Handler) pendingRoomChanges(c *gin.Context) {
	rows, err := h.svc.PendingRoomChanges(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) reviewOffCampus(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	req, bound := bindJSON[dto.ReviewRequest](c)
	if !bound {
		return
	}
	respondNoContent(c, h.svc.ReviewOffCampus(c.Request.Context(), id, middleware.UserID(c), req.Status))
}

func (h *Handler) pendingOffCampus(c *gin.Context) {
	rows, err := h.svc.PendingOffCampus(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) reviewRepair(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	req, bound := bindJSON[dto.ReviewRequest](c)
	if !bound {
		return
	}
	respondNoContent(c, h.svc.ReviewRepair(c.Request.Context(), id, middleware.UserID(c), req))
}

func (h *Handler) reviewCleaning(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	req, bound := bindJSON[dto.ReviewRequest](c)
	if !bound {
		return
	}
	respondNoContent(c, h.svc.ReviewCleaning(c.Request.Context(), id, middleware.UserID(c), req))
}

func (h *Handler) notifications(c *gin.Context) {
	rows, err := h.svc.Notifications(c.Request.Context(), middleware.UserID(c))
	respond(c, rows, err)
}

func (h *Handler) markNotificationRead(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	respondNoContent(c, h.svc.MarkNotificationRead(c.Request.Context(), middleware.UserID(c), id))
}

func (h *Handler) dashboardSummary(c *gin.Context) {
	rows, err := h.svc.DashboardSummary(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) lowBalanceRooms(c *gin.Context) {
	rows, err := h.svc.LowBalanceRooms(c.Request.Context())
	respond(c, rows, err)
}

func (h *Handler) uploadAttachment(c *gin.Context) {
	var req dto.AttachmentUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		fail(c, err)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		fail(c, err)
		return
	}
	if file.Size > h.cfg.Upload.MaxSize {
		fail(c, errs.ErrBadRequest)
		return
	}
	src, err := file.Open()
	if err != nil {
		fail(c, err)
		return
	}
	defer src.Close()
	data, err := io.ReadAll(io.LimitReader(src, h.cfg.Upload.MaxSize+1))
	if err != nil {
		fail(c, err)
		return
	}
	if int64(len(data)) > h.cfg.Upload.MaxSize {
		fail(c, errs.ErrBadRequest)
		return
	}
	contentType := http.DetectContentType(data)
	if contentType != "image/jpeg" && contentType != "image/png" {
		fail(c, errs.ErrBadRequest)
		return
	}
	meta, err := h.svc.SaveAttachment(c.Request.Context(), middleware.UserID(c), middleware.Role(c), req, file.Filename, contentType, data)
	respondCreated(c, meta, err)
}

func (h *Handler) attachmentMetadata(c *gin.Context) {
	ownerType := c.Query("owner_type")
	ownerIDRaw := c.Query("owner_id")
	if ownerType == "" || ownerIDRaw == "" {
		fail(c, errs.ErrBadRequest)
		return
	}
	ownerID, err := strconv.ParseInt(ownerIDRaw, 10, 64)
	if err != nil {
		fail(c, errs.ErrBadRequest)
		return
	}
	var category *string
	if raw := c.Query("category"); raw != "" {
		category = &raw
	}
	rows, err := h.svc.AttachmentMetadata(c.Request.Context(), middleware.UserID(c), middleware.Role(c), ownerType, ownerID, category)
	respond(c, rows, err)
}

func (h *Handler) downloadAttachment(c *gin.Context) {
	id, ok := paramInt64(c, "id")
	if !ok {
		return
	}
	data, err := h.svc.Attachment(c.Request.Context(), middleware.UserID(c), middleware.Role(c), id)
	if err != nil {
		fail(c, err)
		return
	}
	if data.FileName != nil {
		c.Header("Content-Disposition", `inline; filename="`+strings.ReplaceAll(*data.FileName, `"`, "")+`"`)
	}
	c.Data(http.StatusOK, data.ContentType, data.FileData)
}

func (h *Handler) pendingAllocations(c *gin.Context) {
	rows, err := h.svc.PendingAllocations(c.Request.Context())
	respond(c, rows, err)
}

func okJSON(c *gin.Context, data any) {
	ok(c, data)
}

func respond(c *gin.Context, data any, err error) {
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

func respondCreated(c *gin.Context, data any, err error) {
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, data)
}

func respondID(c *gin.Context, id int64, err error) {
	if err != nil {
		fail(c, err)
		return
	}
	createdID(c, id)
}

func respondNoContent(c *gin.Context, err error) {
	if err != nil {
		fail(c, err)
		return
	}
	noContent(c)
}

func paramInt(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil {
		fail(c, errs.ErrBadRequest)
		return 0, false
	}
	return id, true
}

func paramInt64(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil {
		fail(c, errs.ErrBadRequest)
		return 0, false
	}
	return id, true
}

func queryIntPtr(c *gin.Context, name string) *int {
	raw := c.Query(name)
	if raw == "" {
		return nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}
	return &v
}
