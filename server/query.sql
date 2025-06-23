-- name: AddUser :exec
INSERT INTO user (name, password, email, phone, role_id) VALUES(?, ?, ?, ?, ?);

-- name: GetUserInfo :one
SELECT id, name, phone, role_id FROM user 
WHERE id = ?;

-- name: GetUserId :one
SELECT id FROM user 
WHERE email = ?;

-- name: GetUserPassword :one
SELECT password FROM user 
WHERE id = ?;

-- name: GetUserPasswordEmail :one
SELECT id, password FROM user 
WHERE email = ?;

-- name: GetUserRole :one
SELECT role.name FROM user 
JOIN role ON role.id = user.role_id
WHERE user.id = ?;

-- name: GetApartmentID :one
SELECT apartment_id FROM renting_history 
WHERE is_current IS 1 AND user_id = ?;

-- name: GetApartmentAll :one
SELECT apartment.* FROM renting_history 
LEFT JOIN apartment 
ON renting_history.apartment_id = apartment.id
WHERE renting_history.is_current IS 1 AND renting_history.user_id = ?;

-- name: GetSubconInfo :one
SELECT address, NIP, speciality_id FROM subcontractor 
WHERE user_id = ?;

-- name: AddSubcontractor :exec
INSERT INTO subcontractor (
  user_id, address, NIP, speciality_id
) VALUES (
  ?, ?, ?, ?
);

-- name: GetSubcontractors :many
SELECT subcontractor.*, user.name FROM subcontractor
INNER JOIN user ON subcontractor.user_id = user.id;
;

-- name: AddApartment :one
INSERT INTO apartment (
  name, street, building_number, building_name, flat_number, owner_id
) VALUES(
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetApartments :many
SELECT id, name, street, building_number, building_name, flat_number, owner_id
FROM apartment;

-- name: GetApartmentsAndRent :many
SELECT apartment.*, pricing_history.price FROM apartment
LEFT JOIN pricing_history ON pricing_history.apartment_id = apartment.id AND pricing_history.is_current = 1;

-- name: GetRent :one
SELECT price FROM pricing_history 
WHERE is_current = 1 AND apartment_id = ?;

-- name: ChangeRent1 :exec
UPDATE pricing_history
	SET is_current = 1
	WHERE is_current = 0 AND apartment_id = ?;

-- name: ChangeRent2 :exec
INSERT INTO pricing_history (apartment_id, price) VALUES(?, ?);

-- name: GetActiveRenting :many
SELECT id, apartment_id, user_id, start_date, end_date FROM renting_history WHERE is_current = 1;

-- name: GetActiveRentingID :one
SELECT id, apartment_id, user_id, start_date, end_date FROM renting_history WHERE is_current = 1 AND apartment_id = ?;

-- name: AddNewRenting :exec
INSERT INTO renting_history (apartment_id, user_id, start_date) VALUES(?, ?, ?);

-- name: SetEndDate :exec
UPDATE renting_history SET end_date = ? WHERE id = ?;

-- name: MakeAsEnd :exec
UPDATE renting_history SET is_current = 0 WHERE id = ?;

-- name: GetFaultReports :many
SELECT fault_report.*, apartment.name FROM fault_report
INNER JOIN apartment ON apartment.id = fault_report.apartment_id;

-- name: GetFaultReportsUser :many
SELECT fault_report.*, apartment.name FROM fault_report
INNER JOIN apartment ON apartment.id = fault_report.apartment_id
WHERE fault_report.apartment_id = ?;

-- name: AddFault :exec
INSERT INTO fault_report (title, description, status_id, apartment_id, user_id) VALUES(?, ?, ?, ?, ?);

-- name: UpdateFaultStatus :one
UPDATE fault_report
SET status_id = ?
WHERE id = ?
RETURNING *;

-- name: GetTenets :many
SELECT id, name, email, phone, role_id FROM user WHERE role_id = "2";

-- name: GetTenetsWithRent :many
SELECT user.id, user.name, user.email, user.phone, 
  apartment.id, apartment.name, pricing_history.price
FROM user 
LEFT JOIN renting_history ON renting_history.user_id = user.id 
LEFT JOIN apartment ON renting_history.apartment_id = apartment.id
AND renting_history.is_current = 1
LEFT JOIN pricing_history ON pricing_history.apartment_id = apartment.id
AND pricing_history.is_current = 1
WHERE role_id = "2";

-- name: GetSubcontractorSpec :many
SELECT id, name FROM speciality;

-- name: AddSpec :exec
INSERT INTO speciality (name) VALUES(?);

-- name: AddRepair :exec
INSERT INTO repair (
  title, fault_report_id, date_assigned
) VALUES (
  ?, ?, ?
);


-- name: GetRepair :many
SELECT repair.*, user.name FROM repair
LEFT JOIN subcontractor ON repair.subcontractor_id = subcontractor.id
LEFT JOIN user ON subcontractor.user_id = user.id;

-- name: GetRepairSub :many
SELECT repair.*, user.name FROM repair
LEFT JOIN subcontractor ON repair.subcontractor_id = subcontractor.id
LEFT JOIN user ON subcontractor.user_id = user.id
WHERE subcontractor_id = (SELECT id FROM subcontractor WHERE subcontractor.user_id = ?);

-- name: GetRepairApart :many
SELECT repair.*, user.name FROM repair
LEFT JOIN subcontractor ON repair.subcontractor_id = subcontractor.id
LEFT JOIN user ON subcontractor.user_id = user.id
WHERE fault_report_id = (SELECT id FROM fault_report WHERE apartment_id = ?);

-- name: UpdateSubToRepair :one
UPDATE repair
SET subcontractor_id = ?
WHERE id = ?
RETURNING *;

-- name: UpdateRepairData :one
UPDATE repair
SET status_id = (SELECT id FROM repair_status WHERE name = ?), date_completed = ?
WHERE repair.id = ?
RETURNING *;

-- name: GetAllPayment :many
SELECT * FROM payments;

-- name: AddPayment :exec
INSERT INTO payments (
  amount, due_date, renting_id
  ) VALUES (
  ?, ?, ?
);

-- name: UpdatePayment :one
UPDATE payments
SET status_id = 2, transaction_reference = ?, payment_date = ?
WHERE id = ?
RETURNING *;

-- name: GetPayments :many
SELECT *
FROM payments
WHERE renting_id = ?;

-- name: GetPaymentsId :many
SELECT *
FROM payments
WHERE renting_id = (SELECT id FROM renting_history WHERE user_id = ?);

-- name: GetPendingPaymants :many
SELECT *
FROM payments
WHERE status_id = 1;

-- name: GetOverduePayments :many
SELECT *
FROM payments
WHERE status_id = 3;

-- name: SetPaymanyOverdue :one
UPDATE payments
SET status_id = 3
WHERE id = ?
RETURNING *;
