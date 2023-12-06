package exporter

func NewReport() *Report {
	return &Report{
		Succeded: 0,
		Failed:   0,
		Status:   true,
		Results:  []*TestResult{},
	}
}

func (r *Report) Add(tr *TestResult) {

	r.Results = append(r.Results, tr)
	if !tr.Status {
		r.Failed++
	} else {
		r.Succeded++
	}

	if !tr.Status && r.Status {
		r.Status = false
	}
}
