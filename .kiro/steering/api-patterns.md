---
inclusion: fileMatch
fileMatchPattern: "backend/internal/handlers/*.go"
---

# API Handler Patterns

## Standard Response Format
```go
// Success
c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    result,
    "message": "Berhasil ...",
})

// Error
c.JSON(http.StatusBadRequest, gin.H{
    "success":    false,
    "error_code": "VALIDATION_ERROR",
    "message":    "Deskripsi error",
})
```

## Tenant Context Extraction
```go
yayasanID := getYayasanIDFromContext(c) // 0 = superadmin (no filter)
userID, ok := getUserIDFromContext(c)
```

## Error Code Mapping
| Service Error | HTTP | Error Code |
|---|---|---|
| ErrNotFound | 404 | `NOT_FOUND` |
| ErrAlreadySubmitted | 400 | `FORM_ALREADY_SUBMITTED` |
| ErrInvalidScore | 400 | `INVALID_SCORE` |
| ErrIncompleteScores | 400 | `INCOMPLETE_SCORES` |
| ErrSPPGNotFound | 404 | `SPPG_NOT_FOUND` |
| ErrUnauthorized | 401 | `UNAUTHORIZED` |
| ErrForbidden | 403 | `FORBIDDEN` |
| (generic) | 500 | `INTERNAL_ERROR` |

## ID Parsing
```go
id, err := strconv.ParseUint(c.Param("id"), 10, 32)
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"success": false, "error_code": "INVALID_ID", "message": "ID tidak valid"})
    return
}
```

## Pagination
```go
// Query params: page, page_size
// Response includes meta: { total, page, size }
```
