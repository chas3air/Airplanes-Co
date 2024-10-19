package service

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetLimitTime() (time.Duration, error) {
	limitTimeEnvS := os.Getenv("PSQL_LIMIT_RESPONSE_TIME")
	limitTimeEnv, err := strconv.Atoi(limitTimeEnvS)
	if err != nil {
		return 0, fmt.Errorf("wrong environment limit time variable: %w", err)
	}
	return time.Duration(limitTimeEnv) * time.Second, nil
}
