package review

import "inventorychain/internal/model"

type Checklist struct {
	IdentityOK bool
	LinesOK    bool
	SlotsOK    bool
	Notes      []string
}

func BuildChecklist(record model.Record) Checklist {
	check := Checklist{IdentityOK: record.ID != "" && record.Warehouse != "" && record.Cycle != "", Notes: []string{}}
	check.LinesOK = len(record.Lines) > 0
	check.SlotsOK = len(record.Slots) > 0
	if !check.IdentityOK {
		check.Notes = append(check.Notes, "identity incomplete")
	}
	if !check.LinesOK {
		check.Notes = append(check.Notes, "lines missing")
	}
	if !check.SlotsOK {
		check.Notes = append(check.Notes, "slots missing")
	}
	return check
}

func (c Checklist) Ready() bool {
	return c.IdentityOK && c.LinesOK && c.SlotsOK
}
