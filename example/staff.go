package main

import "time"

type Staff struct {
	StaffId   int       `db:"staff_id"`   //
	StaffName string    `db:"staff_name"` //
	CreatedAt time.Time `db:"created_at"` //
	UpdatedAt time.Time `db:"updated_at"` //
	Age       int       `db:"age"`        //
}
