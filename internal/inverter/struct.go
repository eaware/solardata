package inverter

type ConversionData []struct {
	Directory string                `json:"directory"`
	Items     []ConversionDataItems `json:"items"`
}

type ConversionDataItems struct {
	TitleEN      string         `json:"titleEN"`
	TitlePL      string         `json:"titlePL"`
	Registers    []string       `json:"registers"`
	DomoticzIdx  int            `json:"DomoticzIdx"`
	OptionRanges []OptionRanges `json:"optionRanges"`
	Ratio        float64        `json:"ratio"`
	Unit         string         `json:"unit"`
	Graph        int            `json:"graph"`
	MetricType   string         `json:"metric_type"`
	MetricName   string         `json:"metric_name"`
	LabelName    string         `json:"label_name"`
	LabelValue   string         `json:"label_value"`
}

type OptionRanges struct {
	Key     int    `json:"key"`
	ValueEN string `json:"valueEN"`
	ValuePL string `json:"valuePL"`
}

type OutputData struct {
	Key   string
	Value string
}
