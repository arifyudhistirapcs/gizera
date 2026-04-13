package models

import (
	"time"
)

// Supplier represents a vendor that supplies ingredients
type Supplier struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SPPGID          *uint     `gorm:"index" json:"sppg_id"`
	Name            string    `gorm:"size:200;not null;index" json:"name" validate:"required"`
	ContactPerson   string    `gorm:"size:100" json:"contact_person"`
	PhoneNumber     string    `gorm:"size:20" json:"phone_number"`
	Email           string    `gorm:"size:100" json:"email" validate:"omitempty,email"`
	Address         string    `gorm:"type:text" json:"address"`
	Latitude        float64   `gorm:"default:0" json:"latitude"`
	Longitude       float64   `gorm:"default:0" json:"longitude"`
	ProductCategory string    `gorm:"size:100;index" json:"product_category"`
	IsActive        bool      `gorm:"default:true;index" json:"is_active"`
	OnTimeDelivery  float64   `gorm:"default:0" json:"on_time_delivery"` // percentage
	QualityRating   float64   `gorm:"default:0" json:"quality_rating"`   // 1-5 scale
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Yayasans        []Yayasan `gorm:"many2many:supplier_yayasans" json:"yayasans,omitempty"`
}

// PurchaseOrder represents an order placed with a supplier
type PurchaseOrder struct {
	ID                    uint                `gorm:"primaryKey" json:"id"`
	SPPGID                *uint               `gorm:"index" json:"sppg_id"`
	PONumber              string              `gorm:"uniqueIndex;size:50;not null" json:"po_number" validate:"required"`
	SupplierID            uint                `gorm:"index;not null" json:"supplier_id"`
	OrderDate             time.Time           `gorm:"index;not null" json:"order_date"`
	ExpectedDelivery      time.Time           `gorm:"index" json:"expected_delivery"`
	Status                string              `gorm:"size:30;not null;index" json:"status" validate:"required,oneof=pending revision_by_supplier approved shipping received cancelled"` // pending, revision_by_supplier, approved, shipping, received, cancelled
	TotalAmount           float64             `gorm:"not null" json:"total_amount"`
	SupplierRevisionNotes string              `gorm:"type:text" json:"supplier_revision_notes"`
	ApprovedBy            *uint               `gorm:"index" json:"approved_by"`
	ApprovedAt            *time.Time          `json:"approved_at"`
	CreatedBy             uint                `gorm:"not null;index" json:"created_by"`
	CreatedAt             time.Time           `json:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at"`
	YayasanID             *uint               `gorm:"index" json:"yayasan_id"`
	RABID                 *uint               `gorm:"index" json:"rab_id"`
	TargetSPPGID          *uint               `gorm:"index" json:"target_sppg_id"`
	Supplier              Supplier            `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Approver              *User               `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Creator               User                `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	POItems               []PurchaseOrderItem `gorm:"foreignKey:POID" json:"po_items,omitempty"`
	RAB                   *RAB                `gorm:"foreignKey:RABID" json:"rab,omitempty"`
	TargetSPPG            *SPPG               `gorm:"foreignKey:TargetSPPGID" json:"target_sppg,omitempty"`
	Yayasan               *Yayasan            `gorm:"foreignKey:YayasanID" json:"yayasan,omitempty"`
}

// PurchaseOrderItem represents a line item in a purchase order
type PurchaseOrderItem struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	POID         uint       `gorm:"index;not null" json:"po_id"`
	IngredientID uint       `gorm:"index;not null" json:"ingredient_id"`
	Quantity     float64    `gorm:"not null" json:"quantity" validate:"required,gt=0"`
	UnitPrice    float64    `gorm:"not null" json:"unit_price" validate:"required,gte=0"`
	Subtotal     float64    `gorm:"not null" json:"subtotal"`
	PO           PurchaseOrder `gorm:"foreignKey:POID" json:"po,omitempty"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

// GoodsReceipt represents the receipt of goods from a supplier
type GoodsReceipt struct {
	ID            uint               `gorm:"primaryKey" json:"id"`
	SPPGID        *uint              `gorm:"index" json:"sppg_id"`
	GRNNumber     string             `gorm:"uniqueIndex;size:50;not null" json:"grn_number" validate:"required"`
	POID          uint               `gorm:"index;not null" json:"po_id"`
	ReceiptDate   time.Time          `gorm:"index;not null" json:"receipt_date"`
	InvoicePhoto  string             `gorm:"size:500" json:"invoice_photo"` // cloud storage URL
	ReceivedBy    uint               `gorm:"not null;index" json:"received_by"`
	Notes         string             `gorm:"type:text" json:"notes"`
	QualityRating float64            `gorm:"default:0" json:"quality_rating"` // 0-5 rating for this delivery
	CreatedAt     time.Time          `json:"created_at"`
	PurchaseOrder PurchaseOrder      `gorm:"foreignKey:POID" json:"purchase_order,omitempty"`
	Receiver      User               `gorm:"foreignKey:ReceivedBy" json:"receiver,omitempty"`
	GRNItems      []GoodsReceiptItem `gorm:"foreignKey:GRNID" json:"grn_items,omitempty"`
}

// GoodsReceiptItem represents a line item in a goods receipt
type GoodsReceiptItem struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	GRNID            uint       `gorm:"index;not null" json:"grn_id"`
	IngredientID     uint       `gorm:"index;not null" json:"ingredient_id"`
	OrderedQuantity  float64    `gorm:"not null" json:"ordered_quantity"`
	ReceivedQuantity float64    `gorm:"not null" json:"received_quantity" validate:"required,gte=0"`
	ExpiryDate       *time.Time `gorm:"index" json:"expiry_date"`
	GRN              GoodsReceipt `gorm:"foreignKey:GRNID" json:"grn,omitempty"`
	Ingredient       Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

// InventoryItem represents current stock levels for an ingredient
type InventoryItem struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	SPPGID       *uint      `gorm:"index" json:"sppg_id"`
	IngredientID uint       `gorm:"uniqueIndex;not null" json:"ingredient_id"`
	Quantity     float64    `gorm:"not null" json:"quantity"`
	MinThreshold float64    `gorm:"not null" json:"min_threshold"`
	LastUpdated  time.Time  `gorm:"index;not null" json:"last_updated"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

// InventoryMovement tracks all inventory changes
type InventoryMovement struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	SPPGID       *uint      `gorm:"index" json:"sppg_id"`
	IngredientID uint       `gorm:"index;not null" json:"ingredient_id"`
	MovementType string     `gorm:"size:20;not null;index" json:"movement_type" validate:"required,oneof=in out adjustment"` // in, out, adjustment
	Quantity     float64    `gorm:"not null" json:"quantity"`
	Reference    string     `gorm:"size:100;index" json:"reference"` // GRN number, recipe ID, etc.
	MovementDate time.Time  `gorm:"index;not null" json:"movement_date"`
	CreatedBy    uint       `gorm:"not null;index" json:"created_by"`
	Notes        string     `gorm:"type:text" json:"notes"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
	Creator      User       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// StokOpnameForm represents a physical inventory count form
type StokOpnameForm struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	SPPGID          *uint             `gorm:"index" json:"sppg_id"`
	FormNumber      string            `gorm:"uniqueIndex;size:50;not null" json:"form_number" validate:"required"`
	CreatedBy       uint              `gorm:"not null;index" json:"created_by"`
	CreatedAt       time.Time         `gorm:"index;not null" json:"created_at"`
	Status          string            `gorm:"size:20;not null;index" json:"status" validate:"required,oneof=pending approved rejected"` // pending, approved, rejected
	Notes           string            `gorm:"type:text" json:"notes"`
	ApprovedBy      *uint             `gorm:"index" json:"approved_by"`
	ApprovedAt      *time.Time        `json:"approved_at"`
	RejectionReason string            `gorm:"type:text" json:"rejection_reason"`
	IsProcessed     bool              `gorm:"default:false;index" json:"is_processed"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Creator         User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Approver        *User             `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Items           []StokOpnameItem  `gorm:"foreignKey:FormID" json:"items,omitempty"`
}

// StokOpnameItem represents a line item in a stok opname form
type StokOpnameItem struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	FormID        uint           `gorm:"index;not null" json:"form_id"`
	IngredientID  uint           `gorm:"index;not null" json:"ingredient_id"`
	SystemStock   float64        `gorm:"not null" json:"system_stock"`
	PhysicalCount float64        `gorm:"not null" json:"physical_count" validate:"required,gte=0"`
	Difference    float64        `gorm:"not null" json:"difference"` // physical_count - system_stock
	ItemNotes     string         `gorm:"type:text" json:"item_notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Form          StokOpnameForm `gorm:"foreignKey:FormID" json:"form,omitempty"`
	Ingredient    Ingredient     `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

// RAB - Rencana Anggaran Belanja
type RAB struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	RABNumber         string     `gorm:"uniqueIndex;size:50;not null" json:"rab_number"`
	MenuPlanID        uint       `gorm:"index;not null" json:"menu_plan_id"`
	SPPGID            *uint      `gorm:"index" json:"sppg_id"`
	YayasanID         *uint      `gorm:"index" json:"yayasan_id"`
	Status            string     `gorm:"size:30;not null;index" json:"status" validate:"required,oneof=draft approved_sppg approved_yayasan revision_requested completed"`
	TotalAmount       float64    `gorm:"not null" json:"total_amount"`
	RevisionNotes     string     `gorm:"type:text" json:"revision_notes"`
	ApprovedBySPPG    *uint      `gorm:"index" json:"approved_by_sppg"`
	ApprovedAtSPPG    *time.Time `json:"approved_at_sppg"`
	ApprovedByYayasan *uint      `gorm:"index" json:"approved_by_yayasan"`
	ApprovedAtYayasan *time.Time `json:"approved_at_yayasan"`
	CreatedBy         uint       `gorm:"not null;index" json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	MenuPlan          MenuPlan   `gorm:"foreignKey:MenuPlanID" json:"menu_plan,omitempty"`
	SPPG              *SPPG      `gorm:"foreignKey:SPPGID" json:"sppg,omitempty"`
	SPPGApprover      *User      `gorm:"foreignKey:ApprovedBySPPG" json:"sppg_approver,omitempty"`
	YayasanApprover   *User      `gorm:"foreignKey:ApprovedByYayasan" json:"yayasan_approver,omitempty"`
	Creator           User       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Items             []RABItem  `gorm:"foreignKey:RABID" json:"items,omitempty"`
}

// RABItem - Line item in RAB
type RABItem struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	RABID                 uint           `gorm:"index;not null" json:"rab_id"`
	IngredientID          uint           `gorm:"index;not null" json:"ingredient_id"`
	Quantity              float64        `gorm:"not null" json:"quantity" validate:"required,gt=0"`
	Unit                  string         `gorm:"size:20;not null" json:"unit" validate:"required"`
	UnitPrice             float64        `gorm:"not null" json:"unit_price" validate:"gte=0"`
	Subtotal              float64        `gorm:"not null" json:"subtotal"`
	RecommendedSupplierID *uint          `gorm:"index" json:"recommended_supplier_id"`
	POID                  *uint          `gorm:"index" json:"po_id"`
	GRNID                 *uint          `gorm:"index" json:"grn_id"`
	Status                string         `gorm:"size:20;not null;index" json:"status" validate:"required,oneof=pending po_created grn_received"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	RAB                   RAB            `gorm:"foreignKey:RABID" json:"rab,omitempty"`
	Ingredient            Ingredient     `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
	RecommendedSupplier   *Supplier      `gorm:"foreignKey:RecommendedSupplierID" json:"recommended_supplier,omitempty"`
	PurchaseOrder         *PurchaseOrder `gorm:"foreignKey:POID" json:"purchase_order,omitempty"`
	GoodsReceipt          *GoodsReceipt  `gorm:"foreignKey:GRNID" json:"goods_receipt,omitempty"`
}

// SupplierProduct - Supplier product catalog
type SupplierProduct struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	SupplierID    uint       `gorm:"uniqueIndex:idx_supplier_ingredient;not null" json:"supplier_id"`
	IngredientID  uint       `gorm:"uniqueIndex:idx_supplier_ingredient;not null" json:"ingredient_id"`
	UnitPrice     float64    `gorm:"not null" json:"unit_price" validate:"required,gte=0"`
	MinOrderQty   float64    `gorm:"default:0" json:"min_order_qty" validate:"gte=0"`
	IsAvailable   bool       `gorm:"default:true;index" json:"is_available"`
	StockQuantity float64    `gorm:"default:0" json:"stock_quantity" validate:"gte=0"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Supplier      Supplier   `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Ingredient    Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

// SupplierYayasan - Junction table for supplier-yayasan many-to-many
type SupplierYayasan struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SupplierID uint      `gorm:"uniqueIndex:idx_supplier_yayasan;not null" json:"supplier_id"`
	YayasanID  uint      `gorm:"uniqueIndex:idx_supplier_yayasan;not null" json:"yayasan_id"`
	CreatedAt  time.Time `json:"created_at"`
	Supplier   Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Yayasan    Yayasan   `gorm:"foreignKey:YayasanID" json:"yayasan,omitempty"`
}

// Invoice - Supplier invoice for PO
type Invoice struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	InvoiceNumber string        `gorm:"uniqueIndex;size:50;not null" json:"invoice_number"`
	POID          uint          `gorm:"index;not null" json:"po_id"`
	SupplierID    uint          `gorm:"index;not null" json:"supplier_id"`
	YayasanID     uint          `gorm:"index;not null" json:"yayasan_id"`
	Amount        float64       `gorm:"not null" json:"amount" validate:"required,gt=0"`
	Status        string        `gorm:"size:20;not null;index" json:"status" validate:"required,oneof=pending paid"`
	DueDate       time.Time     `gorm:"index" json:"due_date"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	PurchaseOrder PurchaseOrder `gorm:"foreignKey:POID" json:"purchase_order,omitempty"`
	Supplier      Supplier      `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Payment       *Payment      `gorm:"foreignKey:InvoiceID" json:"payment,omitempty"`
}

// Payment - Payment for invoice
type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	InvoiceID     uint      `gorm:"uniqueIndex;not null" json:"invoice_id"`
	PaymentDate   time.Time `gorm:"index;not null" json:"payment_date"`
	Amount        float64   `gorm:"not null" json:"amount" validate:"required,gt=0"`
	ProofURL      string    `gorm:"size:500" json:"proof_url"`
	PaymentMethod string    `gorm:"size:50;not null" json:"payment_method" validate:"required,oneof=bank_transfer"`
	PaidBy        uint      `gorm:"not null;index" json:"paid_by"`
	CreatedAt     time.Time `json:"created_at"`
	Invoice       Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	Payer         User      `gorm:"foreignKey:PaidBy" json:"payer,omitempty"`
}
