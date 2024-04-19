package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) Insert(moduleInfo *ModuleInfo) error {
	_, err := m.DB.Exec("INSERT INTO module_info (created_at, updated_at, module_name, module_duration, exam_type, version) VALUES ($1, $2, $3, $4, $5, $6)",
		moduleInfo.CreatedAt, moduleInfo.UpdatedAt, moduleInfo.ModuleName, moduleInfo.ModuleDuration, moduleInfo.ExamType, moduleInfo.Version)
	if err != nil {
		return fmt.Errorf("failed to insert module info: %w", err)
	}
	return nil
}

func (m *DBModel) Retrieve(id int) (*ModuleInfo, error) {
	var moduleInfo ModuleInfo
	row := m.DB.QueryRow("SELECT * FROM module_info WHERE id = $1", id)
	err := row.Scan(&moduleInfo.ID, &moduleInfo.CreatedAt, &moduleInfo.UpdatedAt, &moduleInfo.ModuleName, &moduleInfo.ModuleDuration, &moduleInfo.ExamType, &moduleInfo.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve module info: %w", err)
	}
	return &moduleInfo, nil
}

func (m *DBModel) Update(moduleInfo *ModuleInfo) error {
	_, err := m.DB.Exec("UPDATE module_info SET created_at = $1, updated_at = $2, module_name = $3, module_duration = $4, exam_type = $5, version = $6 WHERE id = $7",
		moduleInfo.CreatedAt, moduleInfo.UpdatedAt, moduleInfo.ModuleName, moduleInfo.ModuleDuration, moduleInfo.ExamType, moduleInfo.Version, moduleInfo.ID)
	if err != nil {
		return fmt.Errorf("failed to update module info: %w", err)
	}
	return nil
}

func (m *DBModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM module_info WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete module info: %w", err)
	}
	return nil
}
