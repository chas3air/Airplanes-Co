package service

import (
	"os"
	"strconv"
	"time"

	"github.com/chas3air/Airplanes-Co/Core/Management_customers/internal/config"
)

func GetLimitTime(name string) time.Duration {
	limitTimeEnvS := os.Getenv(name)
	limitTimeEnv, err := strconv.Atoi(limitTimeEnvS)
	if err != nil {
		return time.Duration(config.DEFAULT_LIMIT_TIME)
	}
	return time.Duration(limitTimeEnv) * time.Second
}
