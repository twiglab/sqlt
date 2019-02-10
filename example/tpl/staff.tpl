{{ define "Staff.insert"}}
insert into Staff (
	created_at 
	{{if .StaffId}} ,staff_id {{end}}
	{{if .StaffName}} ,staff_name {{end}}
) values (
	:created_at
	{{if .StaffId}} ,:staff_id {{end}}
	{{if .StaffName}} ,:staff_name {{end}}
)
{{end}}

{{ define "Staff.select"}}
select
	staff_id,
	staff_name,
	created_at
from
	Staff
where
	{{if .StaffId}} staff_id = :staff_id {{end}}
	{{if .StaffName}}and  staff_name = :staff_name {{end}}
{{end}}

{{ define "Staff.update"}}
update
	Staff
set
	updated_at = :updated_at
	{{if .StaffName}},staff_name = :staff_name {{end}}

where
	staff_id = :staff_id
{{end}}

{{ define "Staff.delete"}}
delete from staff
{{end}}
