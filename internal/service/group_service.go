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

type GroupService struct {
    queries *database.Queries
    db      *sql.DB
}

func NewGroupService(db *sql.DB) *GroupService {
    return &GroupService{
        queries: database.New(db),
        db:      db,
    }
}

func (s *GroupService) CreateGroup(ctx context.Context, name, description string, flags int) (*models.Group, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    name = strings.TrimSpace(name)
    if name == "" {
        return nil, fmt.Errorf("group name cannot be empty")
    }
    
    if len(name) > 100 {
        return nil, fmt.Errorf("group name too long (max 100 characters)")
    }
    
    if flags < 0 || flags > 2047 {
        return nil, fmt.Errorf("flags must be between 0 and 2047")
    }
    
    group, err := s.queries.CreateGroup(ctx, database.CreateGroupParams{
        Name:        name,
        Description: sql.NullString{String: description, Valid: description != ""},
        Flags:       sql.NullInt32{Int32: int32(flags), Valid: true},
    })
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key") {
            return nil, fmt.Errorf("group with name '%s' already exists", name)
        }
        return nil, fmt.Errorf("failed to create group: %w", err)
    }
    
    logger.Info("Group created", "id", group.ID, "name", group.Name)
    
    return &models.Group{
        ID:          int(group.ID),
        Name:        group.Name,
        Description: group.Description.String,
        Flags:       int(group.Flags.Int32),
        CreatedAt:   group.CreatedAt.Time,
        UpdatedAt:   group.UpdatedAt.Time,
    }, nil
}

func (s *GroupService) GetGroups(ctx context.Context, filters models.GroupFilters) ([]models.Group, int64, error) {
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
    
    sort := filters.Sort
    validSorts := map[string]bool{
        "name_asc": true, "name_desc": true,
        "created_at_asc": true, "created_at_desc": true,
        "": true,
    }
    if !validSorts[sort] {
        sort = "created_at_desc"
    }
    
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    dbGroups, err := s.queries.GetGroupsWithFilters(ctx, database.GetGroupsWithFiltersParams{
        Column1: filters.Name,
        Column2: filters.Search,
        Column3: sort,
        Limit:   int32(filters.Limit),
        Offset:  int32(filters.Offset),
    })
    if err != nil {
        return nil, 0, err
    }
    
    total, err := s.queries.CountGroupsWithFilters(ctx, database.CountGroupsWithFiltersParams{
        Column1: filters.Name,
        Column2: filters.Search,
    })
    if err != nil {
        return nil, 0, err
    }
    
    groups := make([]models.Group, len(dbGroups))
    for i, g := range dbGroups {
        groups[i] = models.Group{
            ID:          int(g.ID),
            Name:        g.Name,
            Description: g.Description.String,
            Flags:       int(g.Flags.Int32),
            CreatedAt:   g.CreatedAt.Time,
            UpdatedAt:   g.UpdatedAt.Time,
        }
    }
    
    return groups, total, nil
}

func (s *GroupService) GetGroupByID(ctx context.Context, id int) (*models.Group, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    if id <= 0 {
        return nil, fmt.Errorf("invalid group id: %d", id)
    }
    
    group, err := s.queries.GetGroupByID(ctx, int32(id))
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("group with id %d not found", id)
        }
        return nil, err
    }
    
    return &models.Group{
        ID:          int(group.ID),
        Name:        group.Name,
        Description: group.Description.String,
        Flags:       int(group.Flags.Int32),
        CreatedAt:   group.CreatedAt.Time,
        UpdatedAt:   group.UpdatedAt.Time,
    }, nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, id int, name, description string, flags int) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    if id <= 0 {
        return fmt.Errorf("invalid group id: %d", id)
    }
    
    if _, err := s.GetGroupByID(ctx, id); err != nil {
        return err
    }
    
    name = strings.TrimSpace(name)
    if name == "" {
        return fmt.Errorf("group name cannot be empty")
    }
    
    if flags < 0 || flags > 2047 {
        return fmt.Errorf("flags must be between 0 and 2047")
    }
    
    err := s.queries.UpdateGroup(ctx, database.UpdateGroupParams{
        ID:          int32(id),
        Name:        name,
        Description: sql.NullString{String: description, Valid: description != ""},
        Flags:       sql.NullInt32{Int32: int32(flags), Valid: true},
    })
    if err != nil {
        return fmt.Errorf("failed to update group: %w", err)
    }
    
    logger.Info("Group updated", "id", id, "name", name)
    
    return nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, id int) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    if id <= 0 {
        return fmt.Errorf("invalid group id: %d", id)
    }
    
    if _, err := s.GetGroupByID(ctx, id); err != nil {
        return err
    }
    
    err := s.queries.DeleteGroup(ctx, int32(id))
    if err != nil {
        return fmt.Errorf("failed to delete group: %w", err)
    }
    
    logger.Info("Group deleted", "id", id)
    
    return nil
}

func (s *GroupService) AddUserToGroup(ctx context.Context, userID, groupID int) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    if userID <= 0 {
        return fmt.Errorf("invalid user id: %d", userID)
    }
    if groupID <= 0 {
        return fmt.Errorf("invalid group id: %d", groupID)
    }
    
    if _, err := s.GetGroupByID(ctx, groupID); err != nil {
        return fmt.Errorf("group not found: %w", err)
    }
    
    err := s.queries.AddUserToGroup(ctx, database.AddUserToGroupParams{
        UserID:  int32(userID),
        GroupID: int32(groupID),
    })
    if err != nil {
        return fmt.Errorf("failed to add user to group: %w", err)
    }
    
    logger.Info("User added to group", "user_id", userID, "group_id", groupID)
    
    return nil
}

func (s *GroupService) RemoveUserFromGroup(ctx context.Context, userID, groupID int) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    if userID <= 0 {
        return fmt.Errorf("invalid user id: %d", userID)
    }
    if groupID <= 0 {
        return fmt.Errorf("invalid group id: %d", groupID)
    }
    
    err := s.queries.RemoveUserFromGroup(ctx, database.RemoveUserFromGroupParams{
        UserID:  int32(userID),
        GroupID: int32(groupID),
    })
    if err != nil {
        return fmt.Errorf("failed to remove user from group: %w", err)
    }
    
    logger.Info("User removed from group", "user_id", userID, "group_id", groupID)
    
    return nil
}

func (s *GroupService) GetUserGroups(ctx context.Context, userID int) ([]models.Group, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    if userID <= 0 {
        return nil, fmt.Errorf("invalid user id: %d", userID)
    }
    
    dbGroups, err := s.queries.GetUserGroups(ctx, int32(userID))
    if err != nil {
        return nil, err
    }
    
    groups := make([]models.Group, len(dbGroups))
    for i, g := range dbGroups {
        groups[i] = models.Group{
            ID:          int(g.ID),
            Name:        g.Name,
            Description: g.Description.String,
            Flags:       int(g.Flags.Int32),
            CreatedAt:   g.CreatedAt.Time,
            UpdatedAt:   g.UpdatedAt.Time,
        }
    }
    
    return groups, nil
}

func (s *GroupService) GetUserFlags(ctx context.Context, userID int) (int, error) {
    groups, err := s.GetUserGroups(ctx, userID)
    if err != nil {
        return 0, err
    }
    return models.MergeFlags(groups), nil
}
