package main

import "time"

type FeeReport struct {
	BlockNumber int64           `json:"curr_block_number"`
	BlockTime   time.Time       `json:"curr_block_time"`
	MaxFeeRate  int64           `json:"max_fee_rate"`
	AvgFeeRate  int64           `json:"avg_fee_rate"`
	MinFeeRate  int64           `json:"min_fee_rate"`
	MaxFee      int64           `json:"max_fee"`
	AvgFee      int64           `json:"avg_fee"`
	MinFee      int64           `json:"min_fee"`
	Suggestions SuggestionsInfo `json:"suggestions"`
}

type SuggestionsInfo struct {
	SuggestedBlock int64     `json:"suggested_block"`
	RequestTime    time.Time `json:"request_time"`
	API            *APIFee   `json:"api_fee"`
	Node           int64     `json:"node_fee"`
}

type APIFee struct {
	Limits struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"limits"`
	Regular  int `json:"regular"`
	Priority int `json:"priority"`
}
