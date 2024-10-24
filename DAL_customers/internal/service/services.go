package service

import (
	"os"
	"strconv"
	"time"

	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/config"
)

func GetLimitTime() time.Duration {
	limitTimeEnvS := os.Getenv("PSQL_LIMIT_RESPONSE_TIME")
	limitTimeEnv, err := strconv.Atoi(limitTimeEnvS)
	if err != nil {
		return config.DEFAULT_LIMIT_TIME * time.Second
	}
	return time.Duration(limitTimeEnv) * time.Second
}
