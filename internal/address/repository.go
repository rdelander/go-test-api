package address

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-test-api/internal/address/db"
	userdb "go-test-api/internal/user/db"

	"github.com/jackc/pgx/v5/pgtype"
)

// Repo defines the interface for address data access
type Repo interface {
	Create(ctx context.Context, req *CreateAddressRequest) (*AddressResponse, error)
	Get(ctx context.Context, id int32) (*AddressResponse, error)
	ListByEntity(ctx context.Context, entityType, entityID string) ([]*AddressResponse, error)
	ListByEntityAndType(ctx context.Context, entityType, entityID, addressType string) ([]*AddressResponse, error)
	Update(ctx context.Context, id int32, req *UpdateAddressRequest) (*AddressResponse, error)
	Delete(ctx context.Context, id int32) error
}

// Repository handles address data access
type Repository struct {
	queries     *db.Queries
	userQueries *userdb.Queries
}

// NewRepository creates a new Repository
func NewRepository(queries *db.Queries, userQueries *userdb.Queries) *Repository {
	return &Repository{
		queries:     queries,
		userQueries: userQueries,
	}
}

// Create creates a new address
func (r *Repository) Create(ctx context.Context, req *CreateAddressRequest) (*AddressResponse, error) {
	// Validate that the entity exists
	if req.EntityType == "user" {
		_, err := r.userQueries.GetUser(ctx, req.EntityID)
		if err != nil {
			return nil, fmt.Errorf("user with id %d does not exist: %w", req.EntityID, err)
		}
	}

	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	addr, err := r.queries.CreateAddress(ctx, db.CreateAddressParams{
		EntityType:  db.EntityType(req.EntityType),
		EntityID:    req.EntityID,
		AddressType: db.AddressType(req.AddressType),
		StreetLine1: req.StreetLine1,
		StreetLine2: pgtype.Text{String: req.StreetLine2, Valid: req.StreetLine2 != ""},
		City:        req.City,
		State:       req.State,
		PostalCode:  req.PostalCode,
		Country:     req.Country,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	return toAddressResponse(addr), nil
}

// Get retrieves an address by ID
func (r *Repository) Get(ctx context.Context, id int32) (*AddressResponse, error) {
	addr, err := r.queries.GetAddress(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	return toAddressResponse(addr), nil
}

// ListByEntity retrieves all addresses for an entity
func (r *Repository) ListByEntity(ctx context.Context, entityType, entityID string) ([]*AddressResponse, error) {
	entityIdInt, err := stringToInt32(entityID)
	if err != nil {
		return nil, err
	}

	addrs, err := r.queries.ListAddressesByEntity(ctx, db.ListAddressesByEntityParams{
		EntityType: db.EntityType(entityType),
		EntityID:   int32(entityIdInt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}

	res := make([]*AddressResponse, len(addrs))
	for i, a := range addrs {
		res[i] = toAddressResponse(a)
	}
	return res, nil
}

// ListByEntityAndType retrieves addresses for an entity filtered by type
func (r *Repository) ListByEntityAndType(ctx context.Context, entityType, entityID, addressType string) ([]*AddressResponse, error) {
	entityIdInt, err := stringToInt32(entityID)
	if err != nil {
		return nil, err
	}
	addrs, err := r.queries.ListAddressesByEntityAndType(ctx, db.ListAddressesByEntityAndTypeParams{
		EntityType:  db.EntityType(entityType),
		EntityID:    int32(entityIdInt),
		AddressType: db.AddressType(addressType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}

	res := make([]*AddressResponse, len(addrs))
	for i, a := range addrs {
		res[i] = toAddressResponse(a)
	}
	return res, nil
}

// Update updates an existing address
func (r *Repository) Update(ctx context.Context, id int32, req *UpdateAddressRequest) (*AddressResponse, error) {
	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	addr, err := r.queries.UpdateAddress(ctx, db.UpdateAddressParams{
		ID:          id,
		StreetLine1: req.StreetLine1,
		StreetLine2: pgtype.Text{String: req.StreetLine2, Valid: req.StreetLine2 != ""},
		City:        req.City,
		State:       req.State,
		PostalCode:  req.PostalCode,
		Country:     req.Country,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}
	return toAddressResponse(addr), nil
}

// Delete deletes an address
func (r *Repository) Delete(ctx context.Context, id int32) error {
	if err := r.queries.DeleteAddress(ctx, id); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	return nil
}

func toAddressResponse(addr db.Address) *AddressResponse {
	return &AddressResponse{
		ID:          fmt.Sprintf("%d", addr.ID),
		EntityType:  string(addr.EntityType),
		EntityID:    fmt.Sprintf("%d", addr.EntityID),
		AddressType: string(addr.AddressType),
		StreetLine1: addr.StreetLine1,
		StreetLine2: addr.StreetLine2.String,
		City:        addr.City,
		State:       addr.State,
		PostalCode:  addr.PostalCode,
		Country:     addr.Country,
	}
}

func stringToInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid integer string %q: %w", s, err)
	}
	return int32(i), nil
}
