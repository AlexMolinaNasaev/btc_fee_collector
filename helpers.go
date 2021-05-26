package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteReport(report *FeeReport) error {
	f, err := os.OpenFile("report/report.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	reportJson, err := json.Marshal(report)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(f, string(reportJson))
	if err != nil {
		f.Close()
		return err
	}

	return nil
}
