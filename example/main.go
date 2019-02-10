package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/twiglab/sqlt"
)

type StaffHandler struct {
	Staffs []*Staff
}

func (sh *StaffHandler) Extract(rs sqlt.Rows) (err error) {
	for rs.Next() {
		staff := new(Staff)
		if err = rs.StructScan(staff); err != nil {
			return
		}

		sh.Staffs = append(sh.Staffs, staff)
	}

	return rs.Err()
}

func main() {
	dbx := sqlt.MustConnect("postgres", "dbname=testdb sslmode=disable")
	tpl := sqlt.NewSqlTemplate("tpl/*.tpl")
	tpl.SetDebug(true)

	dbop := sqlt.New(dbx, tpl)

	staff1 := new(Staff)

	staff1.CreatedAt = time.Now()
	staff1.UpdatedAt = time.Now()
	staff1.StaffId = 12345
	staff1.StaffName = "MikeWang"
	sqlt.MustExec(dbop, context.Background(), "Staff.insert", staff1)

	staff2 := new(Staff)
	staff2.CreatedAt = time.Now()
	staff2.UpdatedAt = time.Now()
	staff2.StaffId = 67890
	staff2.StaffName = "It512"
	sqlt.MustExec(dbop, context.Background(), "Staff.insert", staff2)

	staff := new(Staff)
	staff.StaffId = 67890
	h := new(StaffHandler)
	sqlt.MustQuery(dbop, context.Background(), "Staff.select", staff, h)

	for _, s := range h.Staffs {
		fmt.Println(s)
	}

	staff = new(Staff)
	staff.StaffId = 12345
	staff.StaffName = "god"
	staff.UpdatedAt = time.Now()
	sqlt.MustExec(dbop, context.Background(), "Staff.update", staff)

	sqlt.MustExec(dbop, context.Background(), "Staff.delete", nil)
}
