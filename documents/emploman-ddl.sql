SET TIME ZONE 'Asia/Jakarta';

drop schema achmadnr;
create schema achmadnr;
create schema auth;


drop table achmadnr.employee_assignments;

drop table achmadnr.employees;
drop table achmadnr.role_promotions;
drop table achmadnr.roles;
drop table achmadnr.religions;
drop table achmadnr.grades;
drop table achmadnr.echelons;
drop table achmadnr.units;
drop table achmadnr.positions;


--DDL 

create table achmadnr.roles(
	id char(3) unique primary key,
	name varchar(255) unique not null,
	level int not null,
	description text default 'no desc',
	can_add_role boolean default false,
	can_add_employee boolean default false,
	can_add_unit boolean default false,
	can_add_position boolean default false,
	can_add_echelon boolean default false,
	can_add_religion boolean default false,
	can_add_grade boolean default false,
	can_assign_employee_internal boolean default false,
	can_assign_employee_global boolean default false,
	created_at timestamp default now(),
	modified_at timestamp default now()
);

CREATE TABLE achmadnr.role_promotions (
	promoter_role_id CHAR(3) not null,
	from_role_id CHAR(3) NOT NULL,
	to_role_id CHAR(3) NOT NULL,
	PRIMARY KEY (promoter_role_id, from_role_id, to_role_id),
	foreign key (promoter_role_id) references achmadnr.roles(id) on delete cascade,
	FOREIGN KEY (from_role_id) REFERENCES achmadnr.roles(id) ON DELETE CASCADE,
	FOREIGN KEY (to_role_id) REFERENCES achmadnr.roles(id) ON DELETE CASCADE
);


create table achmadnr.religions(
	id char(3) unique primary key,
	name varchar(100) unique not null,
	created_at timestamp default now(),
	modified_at timestamp default now()
);


create table achmadnr.grades(
	id SERIAL unique primary key,
	code varchar(100) unique not null,
	created_at timestamp default now(),
	modified_at timestamp default now()
);


create table achmadnr.echelons(
	id SERIAL unique primary key,
	code varchar(100) unique not null,
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
	photo_url text default 'noimg',
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
	name varchar(255) unique not null,
	address text not null,
	description text default 'no desc',
	created_at timestamp default now(),
	modified_at timestamp default now()
);

create table achmadnr.positions(
	id SERIAL primary key,
	name varchar(255) unique not null,
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

