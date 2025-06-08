-- roles
CREATE TABLE Role (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE User (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL,
    role_id INTEGER NOT NULL DEFAULT 2,  -- default'user'
    password TEXT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES Role(id)
);

CREATE TABLE Owner (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL
);

CREATE TABLE Apartament (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    street TEXT NOT NULL,
    building_number TEXT NOT NULL,
    building_name TEXT,
    flat_number TEXT NOT NULL,
    owner_id INTEGER NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES Owner(id)
);

CREATE TABLE Pricing_History (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    apartment_id INTEGER NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    price REAL NOT NULL,
    is_current BOOLEAN DEFAULT 0,
    FOREIGN KEY (apartment_id) REFERENCES Apartament(id)
);

CREATE TABLE FaultStatus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE RepairStatus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE PaymentStatus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE PaymentMethod (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE ExpenseCategory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE Speciality (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE Subcontractor (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    address TEXT NOT NULL,
    NIP TEXT NOT NULL,
    speciality_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (speciality_id) REFERENCES Speciality(id)
);

CREATE TABLE FaultReport (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    date_reported DATE NOT NULL,
    status_id INTEGER NOT NULL DEFAULT 1,  -- default'open'
    apartment_id INTEGER NOT NULL,
    FOREIGN KEY (status_id) REFERENCES FaultStatus(id),
    FOREIGN KEY (apartment_id) REFERENCES Apartament(id)
);

CREATE TABLE repair (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fault_report_id INTEGER,
    date_assigned DATE NOT NULL,
    date_completed DATE,
    status_id INTEGER NOT NULL DEFAULT 1,  -- default 'pending'
    subcontractor_id INTEGER,
    FOREIGN KEY (fault_report_id) REFERENCES FaultReport(id),
    FOREIGN KEY (subcontractor_id) REFERENCES Subcontractor(id),
    FOREIGN KEY (status_id) REFERENCES RepairStatus(id)
);

CREATE TABLE Renting_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    apartment_id INTEGER,
    user_id INTEGER,
    start_date DATE NOT NULL,
    end_date DATE,
    FOREIGN KEY (apartment_id) REFERENCES Apartament(id),
    FOREIGN KEY (user_id) REFERENCES User(id)
);

CREATE TABLE Repair_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    renting_history_id INTEGER,
    fault_report_id INTEGER,
    FOREIGN KEY (renting_history_id) REFERENCES Renting_history(id),
    FOREIGN KEY (fault_report_id) REFERENCES FaultReport(id)
);

CREATE TABLE payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    apartament_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    payment_date DATE NOT NULL,
    status_id INTEGER NOT NULL DEFAULT 1,  -- default'pending'
    payment_method_id INTEGER NOT NULL,
    transaction_reference TEXT,
    FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (apartament_id) REFERENCES Apartament(id),
    FOREIGN KEY (status_id) REFERENCES PaymentStatus(id),
    FOREIGN KEY (payment_method_id) REFERENCES PaymentMethod(id)
);

CREATE TABLE Expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    expense_date DATE NOT NULL,
    description TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    repair_id INTEGER NULL,
    FOREIGN KEY (repair_id) REFERENCES repair(id),
    FOREIGN KEY (category_id) REFERENCES ExpenseCategory(id)
);

-- view for financial records
CREATE VIEW Financial_Records AS
SELECT 
    'income' AS type,
    p.id AS source_id,
    p.amount,
    p.payment_date AS record_date,
    'Apartment Payment' AS description,
    p.id AS related_payment_id,
    NULL AS related_expense_id,
    u.name AS user_name,
    a.name AS apartment_name
FROM payments p
JOIN User u ON p.user_id = u.id
JOIN Apartament a ON p.apartament_id = a.id

UNION ALL

SELECT 
    'expense' AS type,
    e.id AS source_id,
    e.amount,
    e.expense_date AS record_date,
    e.description,
    NULL AS related_payment_id,
    e.id AS related_expense_id,
    NULL AS user_name,
    NULL AS apartment_name
FROM Expenses e;

-- roles
INSERT INTO Role (name) VALUES 
('admin'),
('tenent'),
('subcontractor');

-- users 
INSERT INTO User (name, email, phone, role_id, password) VALUES
('John Admin', 'admin@example.com', '123456789', 1, 'admin123'),
('Alice Beton', 'alice@example.com', '987654321', 2, 'alice123'),
('Bob Renter', 'bob@example.com', '555123456', 2, 'bob123'),
('Eve Subcontractor', 'eve@example.com', '555987654', 3, 'eve123'),
('Charlie Newman', 'chanew@example.com', '555654321', 3, 'charlie123');


-- owners
INSERT INTO Owner (name, email, phone) VALUES
('Property Owner LLC', 'owner@example.com', '111222333'),
('Jane Smith', 'jane.smith@example.com', '444555666');

-- apartments
INSERT INTO Apartament (name, street, building_number, building_name, flat_number, owner_id) VALUES
('Sunny Apartment', 'Main Street', '10', 'Sunshine Building', 'A5', 1),
('Cozy Studio', 'Oak Avenue', '25', NULL, '3B', 2),
('Luxury Penthouse', 'High Street', '1', 'Grand Tower', 'PH1', 1);

-- pricing history
INSERT INTO Pricing_History (apartment_id, date, price, is_current) VALUES
(1, '2025-01-01', 1200.00, 1),
(2, '2025-01-01', 850.00, 1),
(3, '2025-01-01', 2500.00, 1),
(1, '2022-06-01', 1100.00, 0),
(2, '2022-06-01', 800.00, 0);

-- subcontractors
INSERT INTO Speciality (name) VALUES
('Plumbing'),
('Electrical'),
('HVAC'),
('Cleaning');

INSERT INTO Subcontractor (user_id, address, NIP, speciality_id) VALUES
(4, '123 Contractor St, City', '1234567890', 1),
(5, '456 Repair Ave, Town', '0987654321', 2);


-- fault statuses
INSERT OR IGNORE INTO FaultStatus (name) VALUES 
('open'),
('closed');

-- fault reports
INSERT INTO FaultReport (description, date_reported, status_id, apartment_id) VALUES
('Leaky faucet in kitchen', '2025-05-10', 1, 1),
('Broken heater', '2025-05-15', 1, 2),
('Power outlet not working', '2025-06-01', 2, 2);

-- repair statuses
INSERT OR IGNORE INTO RepairStatus (name) VALUES 
('pending'),
('in_progress'),
('completed');

-- repairs
INSERT INTO repair (fault_report_id, date_assigned, date_completed, status_id, subcontractor_id) VALUES
(1, '2025-05-11', NULL, 2, 1),
(2, '2025-05-16', '2025-05-18', 3, 2),
(3, '2025-06-02', '2025-06-02', 3, 2);

-- renting history
INSERT INTO Renting_history (apartment_id, user_id, start_date, end_date) VALUES
(1, 3, '2025-01-15', '2025-06-30'),
(2, 2, '2025-02-01', NULL),
(3, 3, '2025-03-01', '2025-04-30');

INSERT INTO Repair_history (renting_history_id, fault_report_id) VALUES
(1, 1),
(2, 2),
(3, 3);

-- payment methods
INSERT OR IGNORE INTO PaymentMethod (name) VALUES 
('credit_card'),
('bank_transfer'),
('cash');

-- payment statuses
INSERT OR IGNORE INTO PaymentStatus (name) VALUES 
('pending'),
('completed'),
('failed');

-- payments
INSERT INTO payments (user_id, apartament_id, amount, payment_date, status_id, payment_method_id, transaction_reference) VALUES
(3, 1, 1200.00, '2025-01-01', 2, 1, 'PAY12345'),
(2, 2, 850.00, '2025-02-01', 2, 2, 'PAY12346'),
(3, 3, 2500.00, '2025-03-01', 2, 1, 'PAY12347'),
(3, 1, 1200.00, '2025-02-01', 2, 1, 'PAY12348');

-- categories
INSERT INTO ExpenseCategory (name) VALUES
('Plumbing'),
('Electrical'),
('Maintenance'),
('Cleaning');

-- example expenses
INSERT INTO Expenses (amount, expense_date, description, category_id, repair_id) VALUES
(150.00, '2025-05-18', 'Faucet replacement parts', 1, 1),
(200.00, '2025-05-18', 'Heater repair service', 3, 2),
(75.00, '2025-06-02', 'Outlet replacement', 2, 3),
(120.00, '2025-06-15', 'Monthly cleaning service', 4, NULL);
