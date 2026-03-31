package services

import (
	"errors"
	"fmt"

	"github.com/erp-sppg/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Provisioning service errors
var (
	ErrProvisioningForbidden    = errors.New("peran Anda tidak memiliki wewenang untuk membuat akun pengguna")
	ErrRoleNotAllowed           = errors.New("peran yang dipilih di luar cakupan wewenang Anda")
	ErrSPPGNotInYayasan         = errors.New("SPPG yang dipilih tidak berada di bawah Yayasan Anda")
	ErrSPPGRequired             = errors.New("sppg_id wajib diisi untuk peran tingkat SPPG")
	ErrYayasanRequired          = errors.New("yayasan_id wajib diisi untuk peran kepala_yayasan")
	ErrPasswordRequired         = errors.New("password wajib diisi")
	ErrProvisioningUserNotFound = errors.New("pengguna tidak ditemukan")
)

// operationalRoles lists all SPPG-level operational roles
var operationalRoles = map[string]bool{
	"akuntan":          true,
	"ahli_gizi":        true,
	"pengadaan":        true,
	"chef":             true,
	"packing":          true,
	"driver":           true,
	"asisten_lapangan": true,
	"kebersihan":       true,
}

// allValidRoles lists every valid role in the system
var allValidRoles = map[string]bool{
	"superadmin":       true,
	"admin_bgn":        true,
	"kepala_yayasan":   true,
	"kepala_sppg":      true,
	"akuntan":          true,
	"ahli_gizi":        true,
	"pengadaan":        true,
	"chef":             true,
	"packing":          true,
	"driver":           true,
	"asisten_lapangan": true,
	"kebersihan":       true,
}

// CreateUserRequest holds the data needed to create a new user
type CreateUserRequest struct {
	NIK         string `json:"nik" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FullName    string `json:"full_name" validate:"required"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role" validate:"required"`
	SPPGID      *uint  `json:"sppg_id"`
	YayasanID   *uint  `json:"yayasan_id"`
}

// UpdateUserRequest holds the data for updating a user
type UpdateUserRequest struct {
	Email       *string `json:"email"`
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Role        *string `json:"role"`
	SPPGID      *uint   `json:"sppg_id"`
	YayasanID   *uint   `json:"yayasan_id"`
}

// ProvisioningService handles delegated user provisioning
type ProvisioningService struct {
	db           *gorm.DB
	auditService *AuditTrailService
}

// NewProvisioningService creates a new ProvisioningService
func NewProvisioningService(db *gorm.DB, auditService *AuditTrailService) *ProvisioningService {
	return &ProvisioningService{
		db:           db,
		auditService: auditService,
	}
}

// isOperationalRole returns true if the role is an SPPG-level operational role
func isOperationalRole(role string) bool {
	return operationalRoles[role]
}

// isSPPGLevelRole returns true if the role requires an sppg_id (kepala_sppg + operational)
func isSPPGLevelRole(role string) bool {
	return role == "kepala_sppg" || operationalRoles[role]
}

// CreateUser creates a new user with delegated provisioning rules.
// The creator's role determines which roles they can create and the scope.
func (s *ProvisioningService) CreateUser(req *CreateUserRequest, creator *models.User) (*models.User, error) {
	// Validate password
	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	// Validate target role is a known role
	if !allValidRoles[req.Role] {
		return nil, ErrRoleNotAllowed
	}

	// Enforce delegated provisioning rules
	if err := s.validateProvisioningRules(req, creator); err != nil {
		return nil, err
	}

	// Auto-fill sppg_id and yayasan_id based on creator context
	s.autoFillTenantContext(req, creator)

	// Validate tenant consistency for the target role
	if err := s.validateTenantConsistency(req); err != nil {
		return nil, err
	}

	// Validate SPPG belongs to the correct Yayasan if both are set
	if req.SPPGID != nil && *req.SPPGID > 0 {
		var sppg models.SPPG
		if err := s.db.First(&sppg, *req.SPPGID).Error; err != nil {
			return nil, ErrSPPGNotFound
		}
		if req.YayasanID != nil && *req.YayasanID > 0 && sppg.YayasanID != *req.YayasanID {
			return nil, ErrSPPGNotInYayasan
		}
		// Auto-fill yayasan_id from SPPG's parent
		yID := sppg.YayasanID
		req.YayasanID = &yID
	}

	// Check NIK uniqueness
	var existingByNIK models.User
	if err := s.db.Where("nik = ?", req.NIK).First(&existingByNIK).Error; err == nil {
		return nil, ErrDuplicateNIK
	}

	// Check email uniqueness
	var existingByEmail models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingByEmail).Error; err == nil {
		return nil, ErrDuplicateEmail
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal hash password: %w", err)
	}

	user := models.User{
		NIK:          req.NIK,
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		PhoneNumber:  req.PhoneNumber,
		Role:         req.Role,
		SPPGID:       req.SPPGID,
		YayasanID:    req.YayasanID,
		IsActive:     true,
		CreatedBy:    &creator.ID,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Record audit trail
	if s.auditService != nil {
		s.auditService.RecordAction(creator.ID, "create", "user", fmt.Sprintf("%d", user.ID), nil, &user, "")
	}

	// Reload with associations
	var created models.User
	s.db.Preload("SPPG").Preload("Yayasan").First(&created, user.ID)

	return &created, nil
}

// validateProvisioningRules checks if the creator is allowed to create a user with the given role.
func (s *ProvisioningService) validateProvisioningRules(req *CreateUserRequest, creator *models.User) error {
	switch creator.Role {
	case "superadmin":
		// Superadmin can create any role
		return nil

	case "admin_bgn":
		// Admin BGN is NOT allowed to create any user
		return ErrProvisioningForbidden

	case "kepala_yayasan":
		// Kepala Yayasan can create: kepala_sppg and operational roles
		// Only for SPPGs under their Yayasan
		if req.Role != "kepala_sppg" && !isOperationalRole(req.Role) {
			return ErrRoleNotAllowed
		}
		// Validate SPPG belongs to creator's Yayasan
		if req.SPPGID != nil && *req.SPPGID > 0 {
			if err := s.validateSPPGBelongsToYayasan(*req.SPPGID, creator.YayasanID); err != nil {
				return err
			}
		}
		return nil

	case "kepala_sppg":
		// Kepala SPPG can only create operational roles for their own SPPG
		if !isOperationalRole(req.Role) {
			return ErrRoleNotAllowed
		}
		return nil

	default:
		// All other roles cannot create users
		return ErrProvisioningForbidden
	}
}

// validateSPPGBelongsToYayasan checks that the given SPPG is under the given Yayasan.
func (s *ProvisioningService) validateSPPGBelongsToYayasan(sppgID uint, yayasanID *uint) error {
	if yayasanID == nil || *yayasanID == 0 {
		return ErrSPPGNotInYayasan
	}
	var sppg models.SPPG
	if err := s.db.First(&sppg, sppgID).Error; err != nil {
		return ErrSPPGNotFound
	}
	if sppg.YayasanID != *yayasanID {
		return ErrSPPGNotInYayasan
	}
	return nil
}

// autoFillTenantContext fills sppg_id and yayasan_id based on the creator's context.
func (s *ProvisioningService) autoFillTenantContext(req *CreateUserRequest, creator *models.User) {
	switch creator.Role {
	case "kepala_yayasan":
		// Auto-fill yayasan_id from creator's yayasan
		if req.YayasanID == nil && creator.YayasanID != nil {
			yID := *creator.YayasanID
			req.YayasanID = &yID
		}

	case "kepala_sppg":
		// Auto-fill sppg_id and yayasan_id from creator's context
		if req.SPPGID == nil && creator.SPPGID != nil {
			sID := *creator.SPPGID
			req.SPPGID = &sID
		}
		if req.YayasanID == nil && creator.YayasanID != nil {
			yID := *creator.YayasanID
			req.YayasanID = &yID
		}
	}
}

// validateTenantConsistency ensures the target role has the required tenant fields.
func (s *ProvisioningService) validateTenantConsistency(req *CreateUserRequest) error {
	if isSPPGLevelRole(req.Role) {
		if req.SPPGID == nil || *req.SPPGID == 0 {
			return ErrSPPGRequired
		}
	}
	if req.Role == "kepala_yayasan" {
		if req.YayasanID == nil || *req.YayasanID == 0 {
			return ErrYayasanRequired
		}
	}
	return nil
}

// GetUsers returns users scoped by the requester's tenant context.
func (s *ProvisioningService) GetUsers(requester *models.User) ([]models.User, error) {
	var users []models.User
	query := s.db.Preload("SPPG").Preload("Yayasan").Order("id ASC")

	switch requester.Role {
	case "superadmin":
		// See all users
	case "admin_bgn":
		// See all users (read-only, but can list)
	case "kepala_yayasan":
		// See users in SPPGs under their Yayasan
		if requester.YayasanID == nil {
			return nil, ErrYayasanRequired
		}
		query = query.Where(
			"yayasan_id = ? OR sppg_id IN (?)",
			*requester.YayasanID,
			s.db.Table("sppgs").Select("id").Where("yayasan_id = ?", *requester.YayasanID),
		)
	case "kepala_sppg":
		// See users in their own SPPG
		if requester.SPPGID == nil {
			return nil, ErrSPPGRequired
		}
		query = query.Where("sppg_id = ?", *requester.SPPGID)
	default:
		// Operational roles cannot list users
		return nil, ErrProvisioningForbidden
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID returns a single user by ID, scoped by the requester's tenant context.
func (s *ProvisioningService) GetUserByID(id uint, requester *models.User) (*models.User, error) {
	var user models.User
	query := s.db.Preload("SPPG").Preload("Yayasan")

	if err := query.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProvisioningUserNotFound
		}
		return nil, err
	}

	// Validate the requester can see this user (tenant scoping)
	if err := s.validateUserVisibility(&user, requester); err != nil {
		// Return 404 to prevent enumeration (cross-tenant access prevention)
		return nil, ErrProvisioningUserNotFound
	}

	return &user, nil
}

// validateUserVisibility checks if the requester can see the target user.
func (s *ProvisioningService) validateUserVisibility(target *models.User, requester *models.User) error {
	switch requester.Role {
	case "superadmin", "admin_bgn":
		return nil
	case "kepala_yayasan":
		if requester.YayasanID == nil {
			return ErrProvisioningForbidden
		}
		// Target must be in the same Yayasan or in an SPPG under the Yayasan
		if target.YayasanID != nil && *target.YayasanID == *requester.YayasanID {
			return nil
		}
		if target.SPPGID != nil {
			var sppg models.SPPG
			if err := s.db.First(&sppg, *target.SPPGID).Error; err == nil {
				if sppg.YayasanID == *requester.YayasanID {
					return nil
				}
			}
		}
		return ErrProvisioningForbidden
	case "kepala_sppg":
		if requester.SPPGID == nil {
			return ErrProvisioningForbidden
		}
		if target.SPPGID != nil && *target.SPPGID == *requester.SPPGID {
			return nil
		}
		return ErrProvisioningForbidden
	default:
		return ErrProvisioningForbidden
	}
}

// UpdateUser updates a user's data with tenant scoping.
func (s *ProvisioningService) UpdateUser(id uint, req *UpdateUserRequest, requester *models.User) (*models.User, error) {
	// Get existing user with visibility check
	existing, err := s.GetUserByID(id, requester)
	if err != nil {
		return nil, err
	}

	oldValue := *existing

	updates := map[string]interface{}{}

	if req.Email != nil {
		// Check email uniqueness
		var dup models.User
		if err := s.db.Where("email = ? AND id != ?", *req.Email, id).First(&dup).Error; err == nil {
			return nil, ErrDuplicateEmail
		}
		updates["email"] = *req.Email
	}
	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Role != nil {
		if !allValidRoles[*req.Role] {
			return nil, ErrRoleNotAllowed
		}
		updates["role"] = *req.Role
	}
	if req.SPPGID != nil {
		updates["sppg_id"] = *req.SPPGID
	}
	if req.YayasanID != nil {
		updates["yayasan_id"] = *req.YayasanID
	}

	if len(updates) == 0 {
		return existing, nil
	}

	if err := s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload
	var updated models.User
	s.db.Preload("SPPG").Preload("Yayasan").First(&updated, id)

	// Audit trail
	if s.auditService != nil {
		s.auditService.RecordAction(requester.ID, "update", "user", fmt.Sprintf("%d", id), &oldValue, &updated, "")
	}

	return &updated, nil
}

// SetUserStatus activates or deactivates a user with tenant scoping.
func (s *ProvisioningService) SetUserStatus(id uint, isActive bool, requester *models.User) (*models.User, error) {
	// Get existing user with visibility check
	existing, err := s.GetUserByID(id, requester)
	if err != nil {
		return nil, err
	}

	oldValue := *existing

	if err := s.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", isActive).Error; err != nil {
		return nil, err
	}

	// Reload
	var updated models.User
	s.db.Preload("SPPG").Preload("Yayasan").First(&updated, id)

	action := "activate_user"
	if !isActive {
		action = "deactivate_user"
	}

	if s.auditService != nil {
		s.auditService.RecordAction(requester.ID, action, "user", fmt.Sprintf("%d", id), &oldValue, &updated, "")
	}

	return &updated, nil
}
