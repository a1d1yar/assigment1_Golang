package data

type DepartmentInfo struct {
	ID                 int    `json:"id"`
	DepartmentName     string `json:"department_name"`
	StaffQuantity      int    `json:"staff_quantity"`
	DepartmentDirector string `json:"department_director"`
	ModuleID           int    `json:"module_id"`
}
