package model

type SearchFilter struct {
	Warehouse string
	Cycle     string
	Status    RecordStatus
	SKU       string
	Published *bool
}

type SearchResult struct {
	Records []Record
	Total   int
	Page    int
	Pages   int
}

func (f SearchFilter) Matches(r Record) bool {
	if f.Warehouse != "" && r.Warehouse != f.Warehouse {
		return false
	}
	if f.Cycle != "" && r.Cycle != f.Cycle {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Published != nil && r.Published != *f.Published {
		return false
	}
	if f.SKU != "" {
		found := false
		for _, line := range r.Lines {
			if line.SKU == f.SKU {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func Paginate(records []Record, page, size int) SearchResult {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	total := len(records)
	pages := (total + size - 1) / size
	if pages == 0 {
		pages = 1
	}
	start := (page - 1) * size
	if start >= total {
		return SearchResult{Records: []Record{}, Total: total, Page: page, Pages: pages}
	}
	end := start + size
	if end > total {
		end = total
	}
	return SearchResult{Records: records[start:end], Total: total, Page: page, Pages: pages}
}
