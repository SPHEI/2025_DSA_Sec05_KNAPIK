CREATE TABLE User(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(15) NOT NULL,
    role VARCHAR(10) DEFAULT 'user' NOT NULL 
        CHECK (role IN ('admin', 'user')),
    password VARCHAR(255) NOT NULL
);

CREATE TABLE Subcontractor(
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    address VARCHAR(255) NOT NULL,
    NIP VARCHAR(10) NOT NULL,
    speciality VARCHAR(100) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(id)
);

CREATE TABLE FaultReport(
    id INT PRIMARY KEY AUTO_INCREMENT,
    description TEXT NOT NULL,
    date_reported DATE NOT NULL,
    status VARCHAR(10) DEFAULT 'open' NOT NULL 
        CHECK (status IN ('open', 'closed'))
);
CREATE INDEX date_reported ON FaultReport(date_reported);

CREATE TABLE Repair(
    id INT PRIMARY KEY AUTO_INCREMENT,
    fault_report_id INT,
    date_assigned DATE NOT NULL,
    date_completed DATE,
    status VARCHAR(10) DEFAULT 'pending' NOT NULL 
        CHECK (status IN ('pending', 'in_progress', 'completed')),
    subcontractor_id INT,
    FOREIGN KEY (fault_report_id) REFERENCES FaultReport(id),
    FOREIGN KEY (subcontractor_id) REFERENCES Subcontractor(id)
);
CREATE INDEX date_assigned ON Repair(date_assigned);

CREATE TABLE Apartament(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    street VARCHAR(100) NOT NULL,
    building_number VARCHAR(10) NOT NULL,
    building_name VARCHAR(100),
    flat_number VARCHAR(10) NOT NULL,
    pricing DECIMAL(10, 2) NOT NULL,
    owner_name VARCHAR(100) NOT NULL,
    owner_email VARCHAR(100) NOT NULL UNIQUE,
    owner_phone VARCHAR(15) NOT NULL
);

CREATE TABLE Renting_history(
    id INT PRIMARY KEY AUTO_INCREMENT,
    apartment_id INT,
    user_id INT,
    start_date DATE NOT NULL,
    end_date DATE,
    fault_report_id INT,
    FOREIGN KEY (fault_report_id) REFERENCES FaultReport(id),
    FOREIGN KEY (apartment_id) REFERENCES Apartament(id),
    FOREIGN KEY (user_id) REFERENCES User(id)
);
CREATE INDEX start_date ON Renting_history(start_date);

CREATE TABLE Pricing_History(
    id INT PRIMARY KEY AUTO_INCREMENT,
    apartment_id INT,
    date DATE NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (apartment_id) REFERENCES Apartament(id)
);
CREATE INDEX date ON Pricing_History(date);

CREATE TABLE payments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    apartament_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_date DATE NOT NULL,
    status VARCHAR(10) DEFAULT 'pending' NOT NULL 
        CHECK (status IN ('pending', 'completed', 'failed')),
    payment_method VARCHAR(10) NOT NULL
        CHECK (payment_method IN ('credit_card', 'bank_transfer', 'cash')),
    transaction_reference VARCHAR(100),
    FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (apartament_id) REFERENCES Apartament(id)
);
CREATE INDEX user_id ON payments(user_id);
CREATE INDEX payment_date ON payments(payment_date);

CREATE TABLE Expenses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    amount DECIMAL(10, 2) NOT NULL,
    expense_date DATE NOT NULL,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    repair_id INT NULL,
    FOREIGN KEY (repair_id) REFERENCES Repair(id)
);
CREATE INDEX expense_date ON expenses(expense_date);
CREATE INDEX category ON expenses(category);

CREATE TABLE Financial_Records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    type VARCHAR(10) NOT NULL 
        CHECK (type IN ('income', 'expense')), 
    amount DECIMAL(10, 2) NOT NULL,
    record_date DATE NOT NULL,
    description VARCHAR(255),
    related_payment_id INT NULL,
    related_expense_id INT NULL,
    FOREIGN KEY (related_payment_id) REFERENCES payments(id),
    FOREIGN KEY (related_expense_id) REFERENCES expenses(id)
);
CREATE INDEX record_date ON financial_records(record_date);
CREATE INDEX type ON financial_records(type);