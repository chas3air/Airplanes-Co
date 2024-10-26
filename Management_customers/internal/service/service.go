package service

import (
	"os"
	"strconv"
	"time"

	"github.com/chas3air/Airplanes-Co/Management_customers/internal/config"
)

func GetLimitTime() time.Duration {
	limitTimeEnvS := os.Getenv("LIMIT_RESPONSE_TIME")
	limitTimeEnv, err := strconv.Atoi(limitTimeEnvS)
	if err != nil {
		return time.Duration(config.DEFAULT_LIMIT_TIME)
	}
	return time.Duration(limitTimeEnv) * time.Second
}
