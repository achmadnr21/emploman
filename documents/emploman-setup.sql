
--SETUP INSERT DATA PADA TABLE

--TABEL ROLE
insert into achmadnr.roles(id, name, level, can_add_role, can_add_employee, can_add_unit, can_add_position, can_add_echelon, can_add_religion, can_add_grade, can_assign_employee_internal, can_assign_employee_global)
values
('SUP', 'SUPERADMIN', 5, true, true, true, true, true, true, true, true, true), -- Super user (Developer)
('ADM', 'ADMIN', 4, true, true, true, true, true, true, true, true, false),  -- Admin dapat menambah semua kecuali assign employee global
('MGR', 'MANAGER', 3, false, false, true, true, false, false, false, true, false),  -- Manager dapat menambah unit dan posisi serta assign employee internal
('HRD', 'HUMAN RESOURCES', 2, false, true, false, false, false, false, false, true, true),  -- HRD dapat menambah user, assign employee internal dan global
('USR', 'USER', 1, false, false, false, false, false, false, false, false, false);  -- User hanya bisa mengakses data tanpa kemampuan menambah apa pun

select * from achmadnr.roles;


insert into achmadnr.role_promotions(promoter_role_id, from_role_id, to_role_id)
values
('SUP', 'USR','HRD'),
('SUP', 'HRD','MGR'),
('SUP', 'MGR','ADM'),
('SUP', 'ADM','SUP'),
('ADM', 'USR','HRD'),
('ADM', 'HRD','MGR'),
('MGR', 'USR','HRD');


select * from achmadnr.role_promotions;

SELECT promoter_role_id, from_role_id, to_role_id
FROM achmadnr.role_promotions 
WHERE promoter_role_id = 'SUP';




--TABEL AGAMA
insert into achmadnr.religions(id, name)
values
('ISL', 'ISLAM'),
('KRJ', 'KRISTEN'),
('KTD', 'KATHOLIK'),
('HND', 'HINDU'),
('BUD', 'BUDHA'),
('KON', 'KONGHUCU');

select * from achmadnr.religions;

--TABEL GRADES ATAU GOLONGAN
insert into achmadnr.grades(code)
values
('I/A'),
('I/B'),
('I/C'),
('I/D'),
('II/A'),
('II/B'),
('II/C'),
('II/D'),
('III/A'),
('III/B'),
('III/C'),
('III/D'),
('IV/A'),
('IV/B'),
('IV/C'),
('IV/D');
select * from achmadnr.grades;

--TABEL ECHELONS
insert into achmadnr.echelons(code)
values
('I'),
('II'),
('III'),
('IV'),
('V'),
('VI'),
('VII'),
('VIII');
select * from achmadnr.echelons;


-- MASUKKAN USER SUPERADMIN.
insert into achmadnr.employees(
    role_id, nip, password, full_name, place_of_birth, date_of_birth, gender, phone_number, address, npwp, grade_id, religion_id, echelon_id
)
values
(
    'SUP',  -- Role ID untuk Superadmin
    '000000000000000000',  -- Ganti dengan NIP yang sesuai
    '$2a$12$4dOIVNGQ9hBsaLoA9uhkpeHpyLm4yOyszBd/fkKG9VwCrlf2RpOR6',  -- Encrypted password
    'Rudy Traspac Developer',  -- Ganti dengan nama lengkap
    'Jakarta',  -- Ganti dengan tempat lahir
    '2001-01-01',  -- Ganti dengan tanggal lahir
    'L',  -- L atau P
    '081234567890',  -- Ganti dengan nomor telepon
    'Jl. Contoh No. 10, Jakarta',  -- Ganti dengan alamat lengkap
    NULL,  -- Ganti dengan NPWP jika ada, atau NULL
    16,  -- Ganti dengan ID grade sesuai
    'ISL',  -- Ganti dengan ID agama sesuai
    1  -- Ganti dengan ID echelon sesuai
),
(
    'SUP',  -- Role ID untuk Superadmin
    '000000000000000001',  -- Ganti dengan NIP yang sesuai
    '$2a$12$4dOIVNGQ9hBsaLoA9uhkpeHpyLm4yOyszBd/fkKG9VwCrlf2RpOR6',  -- Encrypted password
    'Rudy Traspac Developer II',  -- Ganti dengan nama lengkap
    'Jakarta',  -- Ganti dengan tempat lahir
    '2001-01-01',  -- Ganti dengan tanggal lahir
    'L',  -- L atau P
    '081234567890',  -- Ganti dengan nomor telepon
    'Jl. Contoh No. 10, Jakarta',  -- Ganti dengan alamat lengkap
    NULL,  -- Ganti dengan NPWP jika ada, atau NULL
    16,  -- Ganti dengan ID grade sesuai
    'ISL',  -- Ganti dengan ID agama sesuai
    1  -- Ganti dengan ID echelon sesuai
);

select * from achmadnr.employees;



-- WAKTUNYA REGION EMPLOYEE ASSIGNMENT BESERTA MAJOR ENTITY
-- UNITS
insert into achmadnr.units(name, address)
values ('Kantor Pusat', 'Jl. Panglima Sudirman, Jakarta Pusat');
select * from achmadnr.units;

insert into achmadnr.positions(name)
values
('Kepala Sekretariat Utama'),
('Penyusun Laporan Keuangan'),
('Surveyor Pemetaan Pertama'),
('Analis Data Survei dan Pemetaan'),
('Perancang Per-UU-an Utama IV/e');

select * from achmadnr.positions;

--sebelum melakukan assignment, maka tampilkan dulu keperluan.
select * from achmadnr.employees;
select * from achmadnr.units;
select * from achmadnr.positions;
select * from achmadnr.employee_assignments;

insert into achmadnr.employee_assignments(
	employee_id, unit_id, position_id, is_active
) values
('fd0336b6-b32d-4b42-8e8a-7f2b2d7a779c',
	1,
	2,
	true
);
select * from achmadnr.employee_assignments;

-- QUERY TESTS

-- tampilkan yang punya record di employee assignment saja
select te.full_name, tea.*
from achmadnr.employees te
right join achmadnr.employee_assignments tea on te.id = tea.employee_id;

-- tampilkan data unit dan posisi dari employee
select te.full_name, tu.name, tp.name
from achmadnr.employees te
right join achmadnr.employee_assignments tea on te.id = tea.employee_id
left join achmadnr.units tu on tea.unit_id = tu.id
left join achmadnr.positions tp on tea.position_id = tp.id;

-- tampilkan semua data dari employee

select
te.nip, te.full_name, te.place_of_birth, te.date_of_birth, te.gender,  tg.code, tec.code, tp.name, tu.name, tu.address, tr.name, te.phone_number, te.npwp
from achmadnr.employees te
right join achmadnr.employee_assignments tea on te.id = tea.employee_id
left join achmadnr.units tu on tea.unit_id = tu.id
left join achmadnr.positions tp on tea.position_id = tp.id
left join achmadnr.religions tr on te.religion_id = tr.id
left join achmadnr.echelons tec on te.echelon_id = tec.id
left join achmadnr.grades tg on te.grade_id = tg.id;

select coalesce(npwp, 'hh') as npwp
from achmadnr.employees e ;
select e.*
from achmadnr.employees e 
right join achmadnr.employee_assignments ea on e.id = ea.employee_id
where ea.unit_id = 1;


select 
ae.nip, ae.full_name, ae.place_of_birth, ae.address, ae.date_of_birth, ae.gender, ag.code, aec.code,
COALESCE(ap.name, '-') as position_name,
coalesce(au.address, '-') as tempat_kerja,
ar.name,
COALESCE(au.name, '-') AS unit_name,
ae.phone_number,  COALESCE(ae.npwp, '-') as npwp
from achmadnr.employees ae
left join achmadnr.employee_assignments aea on ae.id = aea.employee_id
left join achmadnr.units au on aea.unit_id = au.id
left join achmadnr.positions ap on aea.position_id = ap.id
left join achmadnr.grades ag on ae.grade_id = ag.id
left join achmadnr.echelons aec on ae.echelon_id = aec.id
left join achmadnr.religions ar on ae.religion_id = ar.id;


