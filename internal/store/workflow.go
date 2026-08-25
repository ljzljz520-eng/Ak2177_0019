package store

import (
	"inventorychain/internal/model"
)

func (s *Store) SaveWorkflow(workflow model.Workflow) error {
	return s.put(WorkflowsBucket, workflow.ID, workflow)
}

func (s *Store) GetWorkflow(id string) (model.Workflow, error) {
	var workflow model.Workflow
	err := s.get(WorkflowsBucket, id, &workflow)
	return workflow, err
}

func (s *Store) ListWorkflows(recordID string) ([]model.Workflow, error) {
	workflows := make([]model.Workflow, 0)
	err := s.list(WorkflowsBucket, func(data []byte) error {
		var workflow model.Workflow
		if err := unmarshal(data, &workflow); err != nil {
			return err
		}
		if recordID == "" || workflow.RecordID == recordID {
			workflows = append(workflows, workflow)
		}
		return nil
	})
	return workflows, err
}
