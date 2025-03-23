set time zone 'UTC';

DROP TABLE IF EXISTS college_locations;
DROP TABLE IF EXISTS btransactions;
DROP TABLE IF EXISTS btransfer_log;
DROP TABLE IF EXISTS foods;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS colleges ;
DROP TYPE IF EXISTS student_type;
DROP TYPE IF EXISTS admin_type;


CREATE TYPE student_type AS ENUM('NORMAL_STUDENT','DORM_STUDENT');
CREATE TYPE admin_type AS ENUM('COLLEGE_ADMIN','COLLEGE_OWNER','ADMIN','OWNER');


CREATE TABLE colleges (
	id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE PRIMARY KEY,
	college_name varchar(50)
);

CREATE TABLE college_locations (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name varchar(60) NOT NULL CHECK (LENGTH(name) > 3),
	college_id int4 NOT NULL,
	CONSTRAINT fk_location_college FOREIGN KEY (college_id) REFERENCES colleges
);


CREATE TABLE admins (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	username varchar(32) NOT NULL UNIQUE,
	password varchar(60) NOT NULL,
	first_name varchar(32) NOT NULL,
	last_name VARCHAR(50) NOT NULL,
	position admin_type NOT NULL,
	college_id int4 NULL,
	created_date TIMESTAMP DEFAULT now(),
	counter int2 NOT NULL DEFAULT 0,
	CONSTRAINT fk_admin_college FOREIGN KEY (college_id) REFERENCES colleges
);

CREATE TABLE students (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	username varchar(14) NOT NULL UNIQUE CHECK (LENGTH(username) = 14),
	password varchar(60) NOT NULL,
	first_name varchar(32) NOT NULL,
	last_name varchar(50) NOT NULL,
	phone_number varchar(11) NOT NULL CHECK (LENGTH(phone_number) = 11),
	position student_type NOT NULL,
	college_id int4 NOT NULL,
	created_date TIMESTAMP DEFAULT now(),
	added_by int4 NOT NULL,
	counter int2 NOT NULL DEFAULT 0,
	balance int2 NOT NULL DEFAULT 0 CHECK (balance >= 0),
	second_pass varchar(8) NULL CHECK (second_pass IS NULL OR LENGTH(second_pass) BETWEEN 4 AND 8),
	CONSTRAINT fk_student_addedby FOREIGN KEY (added_by) REFERENCES admins,
	CONSTRAINT fk_student_college FOREIGN KEY (college_id) REFERENCES colleges
);


CREATE TABLE btransactions (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	authority varchar(36) NOT NULL UNIQUE,
	ref_id INTEGER NOT NULL UNIQUE,
	payed_by int4,
	price int4,
	created_date TIMESTAMP DEFAULT now(),
	CONSTRAINT fk_btransaction_payedby FOREIGN KEY (payed_by) REFERENCES students
);

CREATE TABLE btransfer_log (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	from_who int4,
	to_who int4,
	price int4,
	created_date TIMESTAMP DEFAULT now(),
	CONSTRAINT fk_btransfer_fromWho FOREIGN KEY (from_who) REFERENCES students,
	CONSTRAINT fk_btransfer_toWho FOREIGN KEY (to_who) REFERENCES students
);

CREATE TABLE foods (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	food_name VARCHAR(60) NOT NULL,
	expensive BOOLEAN,
	college_id int4,
	added_by int4,
	created_date TIMESTAMP DEFAULT now(),
	CONSTRAINT fk_college_id FOREIGN KEY (college_id) REFERENCES colleges,
	CONSTRAINT fk_added_by FOREIGN KEY (added_by) REFERENCES admins
);