---
inclusion: fileMatch
fileMatchPattern: "backend/internal/models/*.go,backend/internal/services/*.go,backend/internal/database/*.go"
---

# Database Patterns

## GORM Model Template
```go
type MyModel struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    SPPGID    *uint     `gorm:"index" json:"sppg_id"`          // tenant scoping
    YayasanID *uint     `gorm:"index" json:"yayasan_id"`       // tenant scoping
    Name      string    `gorm:"size:200;not null" json:"name" validate:"required"`
    IsActive  bool      `gorm:"default:true;index" json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (MyModel) TableName() string { return "my_models" }
```

## Tenant Scoping
Always filter by tenant in service layer:
```go
query := s.db.Model(&models.MyModel{})
if yayasanID > 0 {
    query = query.Where("yayasan_id = ?", yayasanID)
}
```

## Pagination Pattern
```go
var total int64
query.Count(&total)
query.Order("created_at DESC").Limit(pageSize).Offset((page-1)*pageSize).Find(&results)
```

## Transaction Pattern
```go
err := s.db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&parent).Error; err != nil { return err }
    if err := tx.Create(&child).Error; err != nil { return err }
    return nil
})
```

## Common Mistakes to Avoid
- Forgetting to add model to `AllModels()` in models.go
- Missing `gorm:"index"` on frequently queried columns
- Not using `Preload()` causing N+1 queries
- Forgetting tenant filter (yayasan_id/sppg_id) on queries
- Using `string` instead of `*string` for nullable fields
