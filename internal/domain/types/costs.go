package types

// CostAnalysis representa a análise de custos
type CostAnalysis struct {
	Current     *CostData      `json:"current"`
	Recommended *CostData      `json:"recommended"`
	Savings     *ResourceCosts `json:"savings"`
	Currency    string         `json:"currency"`
	Exchange    *ExchangeInfo  `json:"exchange"`
}

// CostData representa dados de custo
type CostData struct {
	Hourly  *ResourceCosts `json:"hourly"`
	Daily   *ResourceCosts `json:"daily"`
	Monthly *ResourceCosts `json:"monthly"`
}

// ResourceCosts representa custos de recursos
type ResourceCosts struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Total  float64 `json:"total"`
}

// ExchangeInfo representa informações de câmbio
type ExchangeInfo struct {
	Rate         float64 `json:"rate"`
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	UpdatedAt    string  `json:"updatedAt"`
}
