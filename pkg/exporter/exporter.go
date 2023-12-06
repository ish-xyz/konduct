package exporter

func NewStdoutExporter(printall bool) Exporter {
	return &StdoutExporter{
		PrintAll: printall,
	}
}
