package service

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/google/uuid"
)

func GetLimitTime() time.Duration {
	limitTimeEnvS := os.Getenv("LIMIT_RESPONSE_TIME")
	limitTimeEnv, err := strconv.Atoi(limitTimeEnvS)
	if err != nil {
		return time.Duration(config.DEFAULT_LIMIT_TIME)
	}
	return time.Duration(limitTimeEnv) * time.Second
}

func IsCustomerEmpty(c models.Customer) bool {
	return c.Id == uuid.Nil && c.Login == "" && c.Password == "" && c.Role == "" && c.Surname == "" && c.Name == ""
}

func GetInput(scanner *bufio.Scanner, prompt string) string {
	fmt.Println(prompt)
	scanner.Scan()
	return scanner.Text()
}

func ParseTime(time_s string, stencil string) (time.Time, error) {
	time_t, err := time.Parse(stencil, time_s)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}, err
	}
	return time_t, nil
}

func GetInt(scanner *bufio.Scanner, query string) (int, error) {
	entry_s := GetInput(scanner, query)
	entry_i, err := strconv.Atoi(entry_s)
	if err != nil {
		fmt.Println("Error parsing int:", err)
		return 0, err
	}
	return entry_i, nil
}

// /flight/logout
func Logout() error {
	_, err := http.NewRequest(http.MethodDelete, config.Backend_url+"/flight/logout", nil)
	return err
}

