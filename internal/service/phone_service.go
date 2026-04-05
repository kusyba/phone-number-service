package service

import (
    "context"
    "database/sql"
    "log"
    "time"
    
    "phone-number-service/internal/database"
    "phone-number-service/internal/models"
    "phone-number-service/internal/utils"
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

func (s *PhoneService) ProcessNumbers(ctx context.Context, numbers []string, source string) *ImportStats {
    stats := &ImportStats{}
    
    for _, rawNumber := range numbers {
        normalized := utils.NormalizeToE164(rawNumber)
        
        if normalized == "" || !utils.ValidateE164(normalized) {
            stats.Errors++
            log.Printf("Invalid number: %s", rawNumber)
            continue
        }
        
        exists, err := s.queries.CheckPhoneExists(ctx, normalized)
        if err != nil {
            stats.Errors++
            continue
        }
        
        if exists {
            stats.Skipped++
            continue
        }
        
        country := utils.GetCountryByCode(normalized)
        region, provider := "", ""
        if country == "Россия" {
            region, provider = utils.GetRussianRegionAndProvider(normalized)
        }
        
        // Используем sql.NullString
        err = s.queries.InsertPhone(ctx, database.InsertPhoneParams{
            Number:   normalized,
            Country:  sql.NullString{String: country, Valid: country != ""},
            Region:   sql.NullString{String: region, Valid: region != ""},
            Provider: sql.NullString{String: provider, Valid: provider != ""},
            Source:   sql.NullString{String: source, Valid: true},
        })
        
        if err != nil {
            stats.Errors++
            continue
        }
        
        stats.Accepted++
    }
    
    log.Printf("Import: accepted=%d, skipped=%d, errors=%d", stats.Accepted, stats.Skipped, stats.Errors)
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
            ID:        int(p.ID),
            Number:    p.Number,
            Country:   p.Country.String,
            Region:    p.Region.String,
            Provider:  p.Provider.String,
            Source:    p.Source.String,
            CreatedAt: p.CreatedAt.Time,
        }
    }
    
    return phones, total, nil
}
