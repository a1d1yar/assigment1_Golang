package data

type DepartmentInfo struct {
	ID                 int64  `json:"id"`
	DepartmentName     string `json:"department_name"`
	StaffQuantity      int    `json:"staff_quantity"`
	DepartmentDirector string `json:"department_director"`
	ModuleID           int64  `json:"module_id"`
}
