-- RBAC model
-- This is a model for RBAC (Role-Based Access Control) system (We Are Not Implementing this but nice to know)
CREATE TABLE achmadnr.roles (
	id CHAR(3) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	description TEXT DEFAULT 'no desc',
	created_at TIMESTAMP DEFAULT now(),
	modified_at TIMESTAMP DEFAULT now()
);
CREATE TABLE achmadnr.permissions (
	id SERIAL PRIMARY KEY,
	code VARCHAR(100) UNIQUE NOT NULL,  
	description TEXT
);
CREATE TABLE achmadnr.role_permissions (
	role_id CHAR(3) NOT NULL,
	permission_id INT NOT NULL,
	PRIMARY KEY (role_id, permission_id),
	FOREIGN KEY (role_id) REFERENCES achmadnr.roles(id) ON DELETE CASCADE,
	FOREIGN KEY (permission_id) REFERENCES achmadnr.permissions(id) ON DELETE CASCADE
);

CREATE TABLE achmadnr.role_promotions (
	from_role_id CHAR(3) NOT NULL,
	to_role_id CHAR(3) NOT NULL,
	PRIMARY KEY (from_role_id, to_role_id),
	FOREIGN KEY (from_role_id) REFERENCES achmadnr.roles(id) ON DELETE CASCADE,
	FOREIGN KEY (to_role_id) REFERENCES achmadnr.roles(id) ON DELETE CASCADE
);


create table achmadnr.religions(
	id char(3) unique primary key,
	name varchar(100),
	created_at timestamp default now(),
	modified_at timestamp default now()
);


create table achmadnr.grades(
	id SERIAL unique primary key,
	code varchar(100),
	created_at timestamp default now(),
	modified_at timestamp default now()
);


create table achmadnr.echelons(
	id SERIAL unique primary key,
	code varchar(100),
	created_at timestamp default now(),
	modified_at timestamp default now()
);

--employees ddl

create table achmadnr.employees(
	id uuid unique default gen_random_uuid() primary key,
	role_id char(3) not null,
	nip varchar(20) unique not null,
	password varchar(255) not null,
	full_name varchar(255) not null,
	place_of_birth varchar(100) not null,
	date_of_birth DATE not null,
	gender char(1) check (gender in('L','P')),
	phone_number varchar(20) not null,
	photo_url text default '/image/profile_picture_default.png',
	address text not null,
	npwp varchar(25) null,
	grade_id int not null,
	religion_id char(3) not null,
	echelon_id int not null,
	created_at timestamp default now(),
	modified_at timestamp default now(),	
	foreign key (grade_id) references achmadnr.grades(id) on delete set null,
	foreign key (religion_id) references achmadnr.religions(id) on delete set null,
	foreign key (echelon_id) references achmadnr.echelons(id) on delete set null,
	foreign key (role_id) references achmadnr.roles(id) on delete set null
);

-- employee assignments

create table achmadnr.units(
	id SERIAL primary key,
	name varchar(255) not null,
	address text not null,
	description text default 'no desc',
	created_at timestamp default now(),
	modified_at timestamp default now()
);

create table achmadnr.positions(
	id SERIAL primary key,
	name varchar(255) not null,
	created_at timestamp default now(),
	modified_at timestamp default now()
);

create table achmadnr.employee_assignments(
	employee_id uuid not null,
	unit_id int not null,
	position_id int not null,
	is_active boolean not null,
	assigned_at timestamp default now(),
	created_at timestamp default now(),
	modified_at timestamp default now(),
	primary key(employee_id, position_id, unit_id),
	foreign key(employee_id) references achmadnr.employees(id),
	foreign key(unit_id) references achmadnr.units(id),
	foreign key(position_id) references achmadnr.positions(id)
);