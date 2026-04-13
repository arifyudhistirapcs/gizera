package models

// AllModels returns a slice of all model types for migration
func AllModels() []interface{} {
	return []interface{}{
		// Organization
		&Yayasan{},
		&SPPG{},

		// User & Authentication
		&User{},
		&AuditTrail{},
		
		// Recipe & Menu Planning - Ingredients & Semi-Finished Goods
		&Ingredient{},
		&SemiFinishedGoods{},
		&SemiFinishedRecipe{},
		&SemiFinishedRecipeIngredient{},
		&SemiFinishedInventory{},
		&SemiFinishedProductionLog{},
		&SemiFinishedMovement{},
		&Recipe{},
		&RecipeItem{},
		&RecipeVersion{},
		&MenuPlan{},
		&MenuItem{},
		&MenuItemSchoolAllocation{},
		
		// Supply Chain & Inventory
		&Supplier{},
		&PurchaseOrder{},
		&PurchaseOrderItem{},
		&GoodsReceipt{},
		&GoodsReceiptItem{},
		&InventoryItem{},
		&InventoryMovement{},
		&StokOpnameForm{},
		&StokOpnameItem{},
		&RAB{},
		&RABItem{},
		&SupplierProduct{},
		&SupplierYayasan{},
		&Invoice{},
		&Payment{},
		
		// Logistics & Distribution
		&School{},
		&DeliveryTask{},
		&DeliveryMenuItem{},
		&ElectronicPOD{},
		&OmprengTracking{},
		&OmprengInventory{},
		&DeliveryRecord{},
		&StatusTransition{},
		&OmprengCleaning{},
		&PickupTask{},
		&DeliveryReview{},
		
		// Human Resources
		&Employee{},
		&Attendance{},
		&WiFiConfig{},
		&GPSConfig{},
		
		// Financial & Asset Management
		&KitchenAsset{},
		&AssetMaintenance{},
		&CashFlowEntry{},
		&BudgetTarget{},
		
		// System Configuration
		&SystemConfig{},
		&Notification{},

		// Risk Assessment
		&SOPCategory{},
		&SOPChecklistItem{},
		&RiskAssessmentForm{},
		&RiskAssessmentItem{},
		&RiskAssessmentCategoryScore{},
		&SPPGOperationalSnapshot{},
	}
}
