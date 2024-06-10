// projects.go

package main

import (
	"fmt"
	"os"
)

// CreateProject creates the resources associated with the project.
func (p *Project) CreateProjectFolder(config AppConfig) error {
	// Create the Project data in the database
	if err := sqliteDB.CreateProject(p); err != nil {
		return err
	}

	// Create the project folder
	projectPath := fmt.Sprintf("%s/projects/%s", config.DataPath, p.Name)
	return os.MkdirAll(projectPath, os.ModePerm)
}

// DeleteProject deletes all resources associated with the project.
func (p *Project) DeleteProject(config AppConfig) error {
	// Delete the Project data from the database
	if err := sqliteDB.DeleteProject(p.Name); err != nil {
		return err
	}

	// Delete the project folder
	projectPath := fmt.Sprintf("%s/projects/%s", config.DataPath, p.Name)
	return os.RemoveAll(projectPath)
}
