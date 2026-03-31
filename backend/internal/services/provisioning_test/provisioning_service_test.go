package provisioning_test

import (
	"testing"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Yayasan{}, &models.SPPG{}, &models.User{}, &models.AuditTrail{})
	require.NoError(t, err)
	return db
}

func seed(t *testing.T, db *gorm.DB) (*models.Yayasan, *models.SPPG) {
	y := models.Yayasan{Kode: "YYS-0001", Nama: "Yayasan Test", IsActive: true}
	require.NoError(t, db.Create(&y).Error)
	s := models.SPPG{Kode: "SPPG-0001", Nama: "SPPG Test", YayasanID: y.ID, IsActive: true}
	require.NoError(t, db.Create(&s).Error)
	return &y, &s
}

func superadmin() *models.User {
	return &models.User{ID: 1, Role: "superadmin", FullName: "Super Admin"}
}

func adminBGN() *models.User {
	return &models.User{ID: 2, Role: "admin_bgn", FullName: "Admin BGN"}
}

func kepalaYayasan(yID uint) *models.User {
	return &models.User{ID: 3, Role: "kepala_yayasan", YayasanID: &yID, FullName: "Kepala Yayasan"}
}

func kepalaSPPG(sID, yID uint) *models.User {
	return &models.User{ID: 4, Role: "kepala_sppg", SPPGID: &sID, YayasanID: &yID, FullName: "Kepala SPPG"}
}

func TestCreateUser_SuperadminCanCreateAllRoles(t *testing.T) {
	db := setupDB(t)
	yayasan, sppg := seed(t, db)
	svc := services.NewProvisioningService(db, nil)
	creator := superadmin()

	cases := []struct {
		role      string
		sppgID    *uint
		yayasanID *uint
	}{
		{"superadmin", nil, nil},
		{"admin_bgn", nil, nil},
		{"kepala_yayasan", nil, &yayasan.ID},
		{"kepala_sppg", &sppg.ID, &yayasan.ID},
		{"chef", &sppg.ID, &yayasan.ID},
		{"driver", &sppg.ID, &yayasan.ID},
	}

	for _, tc := range cases {
		t.Run(tc.role, func(t *testing.T) {
			req := &services.CreateUserRequest{
				NIK:       "NIK-" + tc.role,
				Email:     tc.role + "@test.com",
				Password:  "password123",
				FullName:  "User " + tc.role,
				Role:      tc.role,
				SPPGID:    tc.sppgID,
				YayasanID: tc.yayasanID,
			}
			user, err := svc.CreateUser(req, creator)
			require.NoError(t, err)
			assert.Equal(t, tc.role, user.Role)
			require.NotNil(t, user.CreatedBy)
			assert.Equal(t, creator.ID, *user.CreatedBy)
		})
	}
}

func TestCreateUser_AdminBGNCannotCreate(t *testing.T) {
	db := setupDB(t)
	seed(t, db)
	svc := services.NewProvisioningService(db, nil)

	req := &services.CreateUserRequest{
		NIK: "NIK-test", Email: "test@test.com", Password: "password123",
		FullName: "Test", Role: "chef",
	}
	_, err := svc.CreateUser(req, adminBGN())
	assert.ErrorIs(t, err, services.ErrProvisioningForbidden)
}

func TestCreateUser_KepalaYayasanAllowedRoles(t *testing.T) {
	db := setupDB(t)
	yayasan, sppg := seed(t, db)
	svc := services.NewProvisioningService(db, nil)
	creator := kepalaYayasan(yayasan.ID)

	// kepala_sppg
	u, err := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-ks", Email: "ks@test.com", Password: "password123",
		FullName: "KS", Role: "kepala_sppg", SPPGID: &sppg.ID,
	}, creator)
	require.NoError(t, err)
	assert.Equal(t, "kepala_sppg", u.Role)
	assert.Equal(t, yayasan.ID, *u.YayasanID)

	// operational
	u2, err := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-chef", Email: "chef@test.com", Password: "password123",
		FullName: "Chef", Role: "chef", SPPGID: &sppg.ID,
	}, creator)
	require.NoError(t, err)
	assert.Equal(t, "chef", u2.Role)
}

func TestCreateUser_KepalaYayasanCannotCreateHigherRoles(t *testing.T) {
	db := setupDB(t)
	yayasan, _ := seed(t, db)
	svc := services.NewProvisioningService(db, nil)
	creator := kepalaYayasan(yayasan.ID)

	for _, role := range []string{"superadmin", "admin_bgn", "kepala_yayasan"} {
		t.Run(role, func(t *testing.T) {
			_, err := svc.CreateUser(&services.CreateUserRequest{
				NIK: "NIK-" + role, Email: role + "@test.com", Password: "password123",
				FullName: "U", Role: role,
			}, creator)
			assert.ErrorIs(t, err, services.ErrRoleNotAllowed)
		})
	}
}

func TestCreateUser_KepalaYayasanCannotCreateForOtherYayasan(t *testing.T) {
	db := setupDB(t)
	yayasan, _ := seed(t, db)

	otherY := models.Yayasan{Kode: "YYS-0002", Nama: "Other", Email: "other-y@test.com", NPWP: "NPWP-002", IsActive: true}
	require.NoError(t, db.Create(&otherY).Error)
	otherS := models.SPPG{Kode: "SPPG-0002", Nama: "Other SPPG", Email: "other-s@test.com", YayasanID: otherY.ID, IsActive: true}
	require.NoError(t, db.Create(&otherS).Error)

	svc := services.NewProvisioningService(db, nil)
	_, err := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-x", Email: "x@test.com", Password: "password123",
		FullName: "X", Role: "chef", SPPGID: &otherS.ID,
	}, kepalaYayasan(yayasan.ID))
	assert.ErrorIs(t, err, services.ErrSPPGNotInYayasan)
}

func TestCreateUser_KepalaSPPGOnlyOperational(t *testing.T) {
	db := setupDB(t)
	yayasan, sppg := seed(t, db)
	svc := services.NewProvisioningService(db, nil)
	creator := kepalaSPPG(sppg.ID, yayasan.ID)

	// operational auto-fills sppg_id
	u, err := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-ak", Email: "ak@test.com", Password: "password123",
		FullName: "Akuntan", Role: "akuntan",
	}, creator)
	require.NoError(t, err)
	assert.Equal(t, sppg.ID, *u.SPPGID)

	// cannot create kepala_sppg
	_, err = svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-ks2", Email: "ks2@test.com", Password: "password123",
		FullName: "KS2", Role: "kepala_sppg",
	}, creator)
	assert.ErrorIs(t, err, services.ErrRoleNotAllowed)
}

func TestCreateUser_DuplicateNIK(t *testing.T) {
	db := setupDB(t)
	_, sppg := seed(t, db)
	svc := services.NewProvisioningService(db, nil)
	creator := superadmin()

	svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-dup", Email: "d1@test.com", Password: "password123",
		FullName: "U1", Role: "chef", SPPGID: &sppg.ID,
	}, creator)

	_, err := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-dup", Email: "d2@test.com", Password: "password123",
		FullName: "U2", Role: "chef", SPPGID: &sppg.ID,
	}, creator)
	assert.ErrorIs(t, err, services.ErrDuplicateNIK)
}

func TestGetUsers_TenantScoping(t *testing.T) {
	db := setupDB(t)
	yayasan, sppg := seed(t, db)

	otherY := models.Yayasan{Kode: "YYS-0002", Nama: "Other", Email: "other-y2@test.com", NPWP: "NPWP-003", IsActive: true}
	require.NoError(t, db.Create(&otherY).Error)
	otherS := models.SPPG{Kode: "SPPG-0002", Nama: "Other SPPG", Email: "other-s2@test.com", YayasanID: otherY.ID, IsActive: true}
	require.NoError(t, db.Create(&otherS).Error)

	svc := services.NewProvisioningService(db, nil)
	creator := superadmin()
	svc.CreateUser(&services.CreateUserRequest{NIK: "N1", Email: "u1@t.com", Password: "pass1234", FullName: "U1", Role: "chef", SPPGID: &sppg.ID}, creator)
	svc.CreateUser(&services.CreateUserRequest{NIK: "N2", Email: "u2@t.com", Password: "pass1234", FullName: "U2", Role: "chef", SPPGID: &otherS.ID}, creator)

	// Superadmin sees all
	users, err := svc.GetUsers(creator)
	require.NoError(t, err)
	assert.Len(t, users, 2)

	// Kepala Yayasan sees only their yayasan
	users, err = svc.GetUsers(kepalaYayasan(yayasan.ID))
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "N1", users[0].NIK)

	// Kepala SPPG sees only their SPPG
	users, err = svc.GetUsers(kepalaSPPG(sppg.ID, yayasan.ID))
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "N1", users[0].NIK)
}

func TestSetUserStatus(t *testing.T) {
	db := setupDB(t)
	_, sppg := seed(t, db)
	svc := services.NewProvisioningService(db, nil)
	creator := superadmin()

	u, _ := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-s", Email: "s@t.com", Password: "password123",
		FullName: "S", Role: "chef", SPPGID: &sppg.ID,
	}, creator)
	assert.True(t, u.IsActive)

	updated, err := svc.SetUserStatus(u.ID, false, creator)
	require.NoError(t, err)
	assert.False(t, updated.IsActive)

	updated, err = svc.SetUserStatus(u.ID, true, creator)
	require.NoError(t, err)
	assert.True(t, updated.IsActive)
}

func TestGetUserByID_CrossTenantReturns404(t *testing.T) {
	db := setupDB(t)
	yayasan, sppg := seed(t, db)

	otherY := models.Yayasan{Kode: "YYS-0002", Nama: "Other", Email: "other-y3@test.com", NPWP: "NPWP-004", IsActive: true}
	require.NoError(t, db.Create(&otherY).Error)
	otherS := models.SPPG{Kode: "SPPG-0002", Nama: "Other SPPG", Email: "other-s3@test.com", YayasanID: otherY.ID, IsActive: true}
	require.NoError(t, db.Create(&otherS).Error)

	svc := services.NewProvisioningService(db, nil)
	other, _ := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-o", Email: "o@t.com", Password: "password123",
		FullName: "O", Role: "chef", SPPGID: &otherS.ID,
	}, superadmin())

	_, err := svc.GetUserByID(other.ID, kepalaSPPG(sppg.ID, yayasan.ID))
	assert.ErrorIs(t, err, services.ErrProvisioningUserNotFound)
}
