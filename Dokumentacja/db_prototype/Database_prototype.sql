CREATE TABLE User(
    id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone INT(15) NOT NULL,
    role VARCHAR(10) DEFAULT 'user' NOT NULL 
        CHECK (role IN ('admin', 'user')),
    password VARCHAR(255) NOT NULL
);

CREATE TABLE Subcontractor(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(15) NOT NULL,
    address VARCHAR(255) NOT NULL,
    NIP VARCHAR(15) NOT NULL,
    speciality VARCHAR(100) NOT NULL
);

CREATE TABLE Repair(
    id INT PRIMARY KEY AUTO_INCREMENT,
    fault_report_id INT,
    date_assigned DATE NOT NULL,
    date_completed DATE,
    status VARCHAR(10) DEFAULT 'pending' NOT NULL 
        CHECK (status IN ('pending', 'in_progress', 'completed')),
    subcontractor_id INT,
    FOREIGN KEY (fault_report_id) REFERENCES FaultReport(id),
    FOREIGN KEY (subcontractor_id) REFERENCES Subcontractor(id),
    INDEX date_assigned
);

CREATE TABLE FaultReport(
    id INT PRIMARY KEY AUTO_INCREMENT,
    description TEXT NOT NULL,
    date_reported DATE NOT NULL,
    status VARCHAR(10) DEFAULT 'open' NOT NULL 
        CHECK (status IN ('open', 'closed')),
    ordered_by_user INT,
    apartament_id INT,
    FOREIGN KEY (apartament_id) REFERENCES Apartament(id),
    FOREIGN KEY (ordered_by_user) REFERENCES User(id),
);
CREATE INDEX date_reported ON FaultReport(date_reported);

CREATE TABLE Apartament(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
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
    FOREIGN KEY (apartment_id) REFERENCES Apartament(id),
    FOREIGN KEY (user_id) REFERENCES User(id)
);
CREATE INDEX start_date ON Renting_history(start_date);

CREATE TABLE Pricing_History(
    id INT PRIMARY KEY AUTO_INCREMENT,
    apartment_id INT,
    date DATE NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (apartament_id) REFERENCES Apartament(id)
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

CREATE TABLE financial_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    type ENUM('deficit', 'earning', 'loss') NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    record_date DATE NOT NULL,
    description VARCHAR(255),
    related_payment_id INT NULL,
    FOREIGN KEY (related_payment_id) REFERENCES payments(id)
);
CREATE INDEX record_date ON financial_records(record_date);
CREATE INDEX type ON financial_records(type);