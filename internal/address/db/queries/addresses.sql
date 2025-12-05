-- name: CreateAddress :one
INSERT INTO addresses (
    entity_type, entity_id, address_type,
    street_line1, street_line2, city, state, postal_code, country,
    created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, entity_type, entity_id, address_type, street_line1, street_line2, 
    city, state, postal_code, country, created_at, updated_at;

-- name: GetAddress :one
SELECT id, entity_type, entity_id, address_type, street_line1, street_line2, 
    city, state, postal_code, country, created_at, updated_at
FROM addresses
WHERE id = $1;

-- name: ListAddressesByEntity :many
SELECT id, entity_type, entity_id, address_type, street_line1, street_line2, 
    city, state, postal_code, country, created_at, updated_at
FROM addresses
WHERE entity_type = $1 AND entity_id = $2
ORDER BY address_type, id;

-- name: ListAddressesByEntityAndType :many
SELECT id, entity_type, entity_id, address_type, street_line1, street_line2, 
    city, state, postal_code, country, created_at, updated_at
FROM addresses
WHERE entity_type = $1 AND entity_id = $2 AND address_type = $3
ORDER BY id;

-- name: UpdateAddress :one
UPDATE addresses
SET street_line1 = $2,
    street_line2 = $3,
    city = $4,
    state = $5,
    postal_code = $6,
    country = $7,
    updated_at = $8
WHERE id = $1
RETURNING id, entity_type, entity_id, address_type, street_line1, street_line2, 
    city, state, postal_code, country, created_at, updated_at;

-- name: DeleteAddress :exec
DELETE FROM addresses
WHERE id = $1;
