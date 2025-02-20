package types

import "time"

// QueryResult representa o resultado de uma query pontual
type QueryResult struct {
	Value     float64
	Timestamp time.Time
}

// QueryRangeResult representa o resultado de uma query com range
type QueryRangeResult struct {
	Values    []QueryResult
	StartTime time.Time
	EndTime   time.Time
}
