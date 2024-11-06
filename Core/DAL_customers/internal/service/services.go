package service

import (
	"os"
	"strconv"
	"time"

	"github.com/chas3air/Airplanes-Co/Core/DAL_customers/internal/config"
)

func GetLimitTime(name string) time.Duration {
	limitTimeEnvS := os.Getenv(name)
	limitTimeEnv, err := strconv.Atoi(limitTimeEnvS)
	if err != nil {
		return config.DEFAULT_LIMIT_TIME * time.Second
	}
	return time.Duration(limitTimeEnv) * time.Second
}
