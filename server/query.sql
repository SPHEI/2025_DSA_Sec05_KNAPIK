-- name: AddUser :exec
INSERT INTO User (name, password, email, phone, role_id) VALUES(?, ?, ?, ?, ?);

-- name: GetUserInfo :one
SELECT id, name, phone, role_id FROM User 
WHERE id = ?;

-- name: GetUserId :one
SELECT id FROM User 
WHERE email = ?;

-- name: GetUserPassword :one
SELECT password FROM User 
WHERE id = ?;

-- name: GetUserPasswordEmail :one
SELECT id, password FROM User 
WHERE email = ?;

-- name: GetUserRole :one
SELECT role_id FROM User 
WHERE id = ?;

-- name: GetApartmentID :one
SELECT apartment_id FROM renting_history 
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

-- name: GetSubcontractors :many
SELECT Subcontractor.*, User.name FROM Subcontractor
INNER JOIN User ON Subcontractor.user_id = User.id;
;

-- name: AddApartment :exec
INSERT INTO Apartament (
  name, street, building_number, building_name, flat_number, owner_id
) VALUES(
  ?, ?, ?, ?, ?, ?
);

-- name: GetApartments :many
SELECT id, name street, building_number, building_name, flat_number, owner_id
FROM Apartament;

-- name: AddOwner :exec
INSERT INTO Owner (
  name, email, phone
) VALUES (
  ?, ?, ?
);

-- name: GetOwners :many
SELECT id, name, email, phone
FROM Owner;

-- name: GetRent :one
SELECT price FROM pricinghistory 
WHERE is_current = 0 AND apartment_id = ?;

-- name: ChangeRent1 :exec
UPDATE pricinghistory
	SET is_current = 1
	WHERE is_current = 0 AND apartment_id = ?;

-- name: ChangeRent2 :exec
INSERT INTO pricinghistory (apartment_id, price) VALUES(?, ?);

-- name: GetActiveRenting :many
SELECT id, apartment_id, user_id, start_date FROM renting_history WHERE end_date IS NULL;

-- name: AddNewRenting :exec
INSERT INTO renting_history (apartment_id, user_id, start_date) VALUES(?, ?, ?);

-- name: SetEndDate :exec
UPDATE renting_history SET end_date = ? WHERE id = ?;

-- name: GetFaultReports :many
SELECT faultreport.*, Apartament.name FROM faultreport
INNER JOIN Apartament ON Apartament.id = faultreport.apartment_id;

-- name: GetFaultReportsUser :many
SELECT faultreport.*, Apartament.name FROM faultreport
INNER JOIN Apartament ON Apartament.id = faultreport.apartment_id
WHERE faultreport.apartment_id = ?;

-- name: AddFault :exec
INSERT INTO faultreport (title, description, date_reported, status_id, apartment_id) VALUES(?, ?, ?, ?, ?);

-- name: UpdateFaultStatus :one
UPDATE faultreport
SET status_id = ?
WHERE id = ?
RETURNING *;

-- name: GetTenets :many
SELECT id, name, email, phone, role_id FROM User WHERE role_id = "2";

-- name: GetSubcontractorSpec :many
SELECT id, name FROM Speciality;

-- name: AddSpec :exec
INSERT INTO Speciality (name) VALUES(?);

-- name: AddRepair :exec
INSERT INTO repair (
  title, fault_report_id, date_assigned
) VALUES (
  ?, ?, ?
);


-- name: GetRepair :many
SELECT repair.*, User.name FROM repair
LEFT JOIN Subcontractor ON repair.subcontractor_id = Subcontractor.id
LEFT JOIN User ON Subcontractor.user_id = User.id;

-- name: GetRepairSub :many
SELECT repair.*, User.name FROM repair
LEFT JOIN Subcontractor ON repair.subcontractor_id = Subcontractor.id
LEFT JOIN User ON Subcontractor.user_id = User.id
WHERE subcontractor_id = (SELECT id FROM Subcontractor WHERE Subcontractor.user_id = ?);

-- name: GetRepairApart :many
SELECT repair.*, User.name FROM repair
LEFT JOIN Subcontractor ON repair.subcontractor_id = Subcontractor.id
LEFT JOIN User ON Subcontractor.user_id = User.id
WHERE fault_report_id = (SELECT id FROM FaultReport WHERE apartment_id = ?);

-- name: UpdateSubToRepair :one
UPDATE repair
SET subcontractor_id = ?
WHERE id = ?
RETURNING *;

-- name: UpdateRepairData :one
UPDATE repair
SET status_id = (SELECT id FROM RepairStatus WHERE name = ?), date_completed = ?
WHERE repair.id = ?
RETURNING *;
