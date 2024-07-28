package models

type Statistic struct {
	ProcessedMessages     int     `json:"total_processed"`
	AverageProcessingTime float32 `json:"average_processing_time"`
}
