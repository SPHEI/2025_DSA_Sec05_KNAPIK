-- name: GetUserInfo :one
SELECT id, name, phone, role_id FROM User 
WHERE id = ?;

-- name: GetUserId :one
SELECT id FROM User 
WHERE email = ?;

-- name: GetUserPassword :one
SELECT password FROM User 
WHERE id = ?;

-- name: GetUserRole :one
SELECT role_id FROM User 
WHERE id = ?;

-- name: GetApartmentID :one
SELECT apartment_id FROM Renting_History 
WHERE end_date IS NULL AND user_id = ?;

-- name: GetSubconInfo :one
SELECT address, NIP, speciality_id FROM Subcontractor 
WHERE user_id = ?;

-- name: AddSubcontractor :exec
INSERT INTO Subcontractor (
  user_id, address, NIP, speciality_id
) VALUES (
  ?, ?, ?, ?
);

-- name: AddApartment :exec
INSERT INTO Apartament (
  name, street, building_number, building_name, flat_number, owner_id
) VALUES(
  ?, ?, ?, ?, ?, ?
);

-- name: AddOwner :exec
INSERT INTO Owner (
  name, email, phone
) VALUES (
  ?, ?, ?
);

-- name: GetRent :one
SELECT price FROM Pricing_History 
WHERE is_current = 0 AND apartment_id = ?;

-- name: AddRepair :exec
INSERT INTO repair (
  fault_report_id, date_assigned
) VALUES (
  ?, ?
);

-- name: GetRepair :many
SELECT id, fault_report_id, date_assigned, date_completed, status_id, subcontractor_id FROM repair;

-- name: UpdateSubToRepair :one
UPDATE repair
SET subcontractor_id = ?
WHERE id = ?
RETURNING *;
