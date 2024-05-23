package crons

import "github.com/robfig/cron/v3"

func Verify(cronStr string) (bool, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(cronStr)
	if err != nil {
		return false, err
	}
	return true, nil
}
