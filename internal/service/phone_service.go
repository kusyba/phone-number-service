package service

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "time"
    
    "phone-number-service/internal/database"
    "phone-number-service/internal/models"
    "phone-number-service/pkg/logger"
)

type PhoneService struct {
    db      *sql.DB
    queries *database.Queries
}

func NewPhoneService(db *sql.DB) *PhoneService {
    return &PhoneService{
        db:      db,
        queries: database.New(db),
    }
}

type ImportStats struct {
    Accepted int
    Skipped  int
    Errors   int
}

// formatPhoneForDisplay безопасно форматирует номер для отображения
func formatPhoneForDisplay(number string) string {
    withoutPlus := strings.TrimPrefix(number, "+")
    
    if len(withoutPlus) == 11 && strings.HasPrefix(withoutPlus, "7") {
        return "+7 (" + withoutPlus[1:4] + ") " + withoutPlus[4:7] + "-" + withoutPlus[7:9] + "-" + withoutPlus[9:11]
    }
    
    if len(withoutPlus) == 11 && strings.HasPrefix(withoutPlus, "8") {
        return "+7 (" + withoutPlus[1:4] + ") " + withoutPlus[4:7] + "-" + withoutPlus[7:9] + "-" + withoutPlus[9:11]
    }
    
    if len(withoutPlus) == 10 {
        return "+7 (" + withoutPlus[0:3] + ") " + withoutPlus[3:6] + "-" + withoutPlus[6:8] + "-" + withoutPlus[8:10]
    }
    
    if len(withoutPlus) >= 10 {
        codeLen := 1
        if strings.HasPrefix(withoutPlus, "44") {
            codeLen = 2
        }
        remaining := withoutPlus[codeLen:]
        formatted := "+" + withoutPlus[:codeLen]
        for i := 0; i < len(remaining); i += 3 {
            end := i + 3
            if end > len(remaining) {
                end = len(remaining)
            }
            if i > 0 {
                formatted += " "
            }
            formatted += remaining[i:end]
        }
        return formatted
    }
    
    return number
}

func (s *PhoneService) ProcessNumbers(ctx context.Context, numbers []string, source string) *ImportStats {
    stats := &ImportStats{}
    
    select {
    case <-ctx.Done():
        logger.Error("Context cancelled", "error", ctx.Err())
        return stats
    default:
    }
    
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        stats.Errors = len(numbers)
        logger.Error("Failed to begin transaction", "error", err)
        return stats
    }
    defer tx.Rollback()
    
    for _, rawNumber := range numbers {
        select {
        case <-ctx.Done():
            logger.Error("Context cancelled during processing", "error", ctx.Err())
            return stats
        default:
        }
        
        if strings.TrimSpace(rawNumber) == "" {
            stats.Errors++
            continue
        }
        
        normalized := rawNumber
        normalized = strings.ReplaceAll(normalized, " ", "")
        normalized = strings.ReplaceAll(normalized, "-", "")
        normalized = strings.ReplaceAll(normalized, "(", "")
        normalized = strings.ReplaceAll(normalized, ")", "")
        
        isValid := true
        for i, ch := range normalized {
            if i == 0 && ch == '+' {
                continue
            }
            if ch < '0' || ch > '9' {
                isValid = false
                break
            }
        }
        
        if !isValid {
            stats.Errors++
            logger.Error("Invalid characters in number", "number", rawNumber)
            continue
        }
        
        digitsCount := len(normalized)
        if strings.HasPrefix(normalized, "+") {
            digitsCount--
        }
        
        if digitsCount < 10 {
            stats.Errors++
            logger.Error("Number too short (min 10 digits)", "number", rawNumber, "digits", digitsCount)
            continue
        }
        
        if digitsCount > 15 {
            stats.Errors++
            logger.Error("Number too long (max 15 digits)", "number", rawNumber, "digits", digitsCount)
            continue
        }
        
        if len(normalized) == 10 {
            normalized = "+7" + normalized
        } else if strings.HasPrefix(normalized, "8") && len(normalized) == 11 {
            normalized = "+7" + normalized[1:]
        } else if !strings.HasPrefix(normalized, "+") && len(normalized) == 11 {
            normalized = "+" + normalized
        }
        
        var exists bool
        checkQuery := `SELECT EXISTS(SELECT 1 FROM phones WHERE number = $1)`
        err := tx.QueryRowContext(ctx, checkQuery, normalized).Scan(&exists)
        if err != nil {
            stats.Errors++
            logger.Error("Failed to check existence", "number", normalized, "error", err)
            continue
        }
        
        if exists {
            stats.Skipped++
            logger.Info("Number already exists, skipped", "number", normalized)
            continue
        }
        
        country := "Россия"
        region, provider := "", ""
        if strings.HasPrefix(normalized, "+7") {
            region = "Москва"
            provider = "МТС"
        }
        
        _, err = tx.ExecContext(ctx,
            `INSERT INTO phones (number, country, region, provider, source, created_at)
             VALUES ($1, $2, $3, $4, $5, NOW())`,
            normalized, country, region, provider, source)
        
        if err != nil {
            stats.Errors++
            logger.Error("Failed to insert number", "number", normalized, "error", err)
            continue
        }
        
        stats.Accepted++
        logger.Info("Insert successful", "number", normalized)
    }
    
    if err := tx.Commit(); err != nil {
        stats.Errors += len(numbers)
        logger.Error("Failed to commit transaction", "error", err)
        return stats
    }
    
    logger.Info("Import completed",
        "accepted", stats.Accepted,
        "skipped", stats.Skipped,
        "errors", stats.Errors)
    
    return stats
}

type SearchFilters struct {
    Number   string
    Country  string
    Region   string
    Provider string
    Limit    int
    Offset   int
}

func (s *PhoneService) SearchPhones(ctx context.Context, filters SearchFilters) ([]models.Phone, int64, error) {
    select {
    case <-ctx.Done():
        return nil, 0, ctx.Err()
    default:
    }
    
    if filters.Limit <= 0 {
        filters.Limit = 10
    }
    if filters.Limit > 100 {
        filters.Limit = 100
    }
    if filters.Offset < 0 {
        filters.Offset = 0
    }
    
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    dbPhones, err := s.queries.GetPhonesWithFilters(ctx, database.GetPhonesWithFiltersParams{
        Column1: filters.Number,
        Column2: filters.Country,
        Column3: filters.Region,
        Column4: filters.Provider,
        Limit:   int32(filters.Limit),
        Offset:  int32(filters.Offset),
    })
    if err != nil {
        return nil, 0, err
    }
    
    total, err := s.queries.CountPhonesWithFilters(ctx, database.CountPhonesWithFiltersParams{
        Column1: filters.Number,
        Column2: filters.Country,
        Column3: filters.Region,
        Column4: filters.Provider,
    })
    if err != nil {
        return nil, 0, err
    }
    
    phones := make([]models.Phone, len(dbPhones))
    for i, p := range dbPhones {
        phones[i] = models.Phone{
            ID:          int(p.ID),
            Number:      p.Number,
            Country:     p.Country.String,
            Region:      p.Region.String,
            Provider:    p.Provider.String,
            Source:      p.Source.String,
            CreatedAt:   p.CreatedAt.Time,
        }
    }
    
    return phones, total, nil
}

func (s *PhoneService) FormatPhoneNumber(ctx context.Context, id int) (string, string, error) {
    select {
    case <-ctx.Done():
        return "", "", ctx.Err()
    default:
    }
    
    if id <= 0 {
        return "", "", fmt.Errorf("invalid phone id: %d", id)
    }
    
    var number string
    query := `SELECT number FROM phones WHERE id = $1`
    err := s.db.QueryRowContext(ctx, query, id).Scan(&number)
    if err != nil {
        return "", "", err
    }
    
    return formatPhoneForDisplay(number), "", nil
}

func (s *PhoneService) FormatPhoneNumberByValue(ctx context.Context, phoneNumber string) (string, string, error) {
    select {
    case <-ctx.Done():
        return "", "", ctx.Err()
    default:
    }
    
    if strings.TrimSpace(phoneNumber) == "" {
        return "", "", fmt.Errorf("empty phone number")
    }
    
    normalized := phoneNumber
    normalized = strings.ReplaceAll(normalized, " ", "")
    normalized = strings.ReplaceAll(normalized, "-", "")
    normalized = strings.ReplaceAll(normalized, "(", "")
    normalized = strings.ReplaceAll(normalized, ")", "")
    
    if len(normalized) == 10 {
        normalized = "+7" + normalized
    } else if strings.HasPrefix(normalized, "8") && len(normalized) == 11 {
        normalized = "+7" + normalized[1:]
    } else if !strings.HasPrefix(normalized, "+") && len(normalized) == 11 {
        normalized = "+" + normalized
    }
    
    return formatPhoneForDisplay(normalized), "", nil
}
