package services

import "gorm.io/gorm"

// WithDB methods return new service instances with a different (tenant-scoped) DB.
// This allows handlers to pass a tenant-scoped DB without modifying the original service.

// WithDB returns a new RecipeService with the given DB.
func (s *RecipeService) WithDB(db *gorm.DB) *RecipeService {
	return &RecipeService{
		db:                 db,
		nutritionStandards: s.nutritionStandards,
	}
}

// WithDB returns a new MenuPlanningService with the given DB.
func (s *MenuPlanningService) WithDB(db *gorm.DB) *MenuPlanningService {
	return &MenuPlanningService{
		db:            db,
		recipeService: s.recipeService.WithDB(db),
	}
}

// WithDB returns a new SupplierService with the given DB.
func (s *SupplierService) WithDB(db *gorm.DB) *SupplierService {
	return &SupplierService{
		db: db,
	}
}

// WithDB returns a new PurchaseOrderService with the given DB.
func (s *PurchaseOrderService) WithDB(db *gorm.DB) *PurchaseOrderService {
	return &PurchaseOrderService{
		db: db,
	}
}

// WithDB returns a new InventoryService with the given DB.
func (s *InventoryService) WithDB(db *gorm.DB) *InventoryService {
	return &InventoryService{
		db: db,
	}
}

// WithDB returns a new GoodsReceiptService with the given DB.
func (s *GoodsReceiptService) WithDB(db *gorm.DB) *GoodsReceiptService {
	return &GoodsReceiptService{
		db:               db,
		inventoryService: s.inventoryService.WithDB(db),
		cashFlowService:  s.cashFlowService.WithDB(db),
	}
}

// WithDB returns a new SchoolService with the given DB.
func (s *SchoolService) WithDB(db *gorm.DB) *SchoolService {
	return &SchoolService{
		db: db,
	}
}

// WithDB returns a new DeliveryTaskService with the given DB.
func (s *DeliveryTaskService) WithDB(db *gorm.DB) *DeliveryTaskService {
	return &DeliveryTaskService{
		db: db,
	}
}

// WithDB returns a new AssetService with the given DB.
func (s *AssetService) WithDB(db *gorm.DB) *AssetService {
	return &AssetService{
		db: db,
	}
}

// WithDB returns a new CashFlowService with the given DB.
func (s *CashFlowService) WithDB(db *gorm.DB) *CashFlowService {
	return &CashFlowService{
		db: db,
	}
}

// WithDB returns a new FinancialReportService with the given DB.
func (s *FinancialReportService) WithDB(db *gorm.DB) *FinancialReportService {
	return &FinancialReportService{
		db:              db,
		cashFlowService: s.cashFlowService.WithDB(db),
		assetService:    s.assetService.WithDB(db),
	}
}

// WithDB returns a new EmployeeService with the given DB.
func (s *EmployeeService) WithDB(db *gorm.DB) *EmployeeService {
	return &EmployeeService{
		db:          db,
		authService: s.authService,
	}
}

// WithDB returns a new AttendanceService with the given DB.
func (s *AttendanceService) WithDB(db *gorm.DB) *AttendanceService {
	return &AttendanceService{
		db:              db,
		employeeService: s.employeeService.WithDB(db),
	}
}

// WithDB returns a new ReviewService with the given DB.
func (s *ReviewService) WithDB(db *gorm.DB) *ReviewService {
	return &ReviewService{
		db: db,
	}
}

// WithDB returns a new SemiFinishedService with the given DB.
func (s *SemiFinishedService) WithDB(db *gorm.DB) *SemiFinishedService {
	return &SemiFinishedService{
		db: db,
	}
}

// WithDB returns a new OmprengTrackingService with the given DB.
func (s *OmprengTrackingService) WithDB(db *gorm.DB) *OmprengTrackingService {
	return &OmprengTrackingService{
		db: db,
	}
}

// WithDB returns a new EPODService with the given DB.
func (s *EPODService) WithDB(db *gorm.DB) *EPODService {
	return &EPODService{
		db:                  db,
		deliveryTaskService: s.deliveryTaskService.WithDB(db),
	}
}

// WithDB returns a new stokOpnameServiceImpl with the given DB.
func (s *stokOpnameServiceImpl) WithDB(db *gorm.DB) StokOpnameService {
	return &stokOpnameServiceImpl{
		db:                  db,
		inventoryService:    s.inventoryService.WithDB(db),
		notificationService: s.notificationService,
	}
}

// WithDB returns a new PickupTaskService with the given DB.
func (s *PickupTaskService) WithDB(db *gorm.DB) *PickupTaskService {
	return &PickupTaskService{
		db:                     db,
		activityTrackerService: s.activityTrackerService,
	}
}

