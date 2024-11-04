package customersadmininterface

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
	"github.com/google/uuid"
)

var limitTime = service.GetLimitTime()

// /customers/get
func GetAllCustomers() ([]models.Customer, error) {
	resp, err := http.Get(config.Backend_url + "/customers/get")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get customers: %s", resp.Status)
	}

	var customers []models.Customer
	err = json.NewDecoder(resp.Body).Decode(&customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// /customers/insert + body(json)
func AddCustomer() (models.Customer, error) {
	scanner := bufio.NewScanner(os.Stdin)

	login := service.GetInput(scanner, "Enter login")
	password := service.GetInput(scanner, "Enter password")
	role := service.GetInput(scanner, "Ernter role")
	surname := service.GetInput(scanner, "Enter surname")
	name := service.GetInput(scanner, "Enter name")

	customer := models.Customer{
		Login:    login,
		Password: password,
		Role:     role,
		Surname:  surname,
		Name:     name,
	}

	return postCustomer(customer)
}

// /customers/update + body(json)
func UpdateCustomer(uuidStr string) (models.Customer, error) {
	scanner := bufio.NewScanner(os.Stdin)

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return models.Customer{}, fmt.Errorf("failed to parse: %s", uuidStr)
	}

	login := service.GetInput(scanner, "Enter login")
	password := service.GetInput(scanner, "Enter password")
	role := service.GetInput(scanner, "Enter role")
	surname := service.GetInput(scanner, "Enter surname")
	name := service.GetInput(scanner, "Enter name")

	customer := models.Customer{
		Id:       id,
		Login:    login,
		Password: password,
		Role:     role,
		Surname:  surname,
		Name:     name,
	}

	return updateCustomer(customer)
}

// /customers/delete?id=...
func DeleteCustomer(id string) (models.Customer, error) {
	req, err := http.NewRequest(http.MethodDelete, config.Backend_url+"/customers/delete?id="+id, nil)
	if err != nil {
		return models.Customer{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: limitTime,
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Customer{}, err
	}
	defer resp.Body.Close()
	
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Customer{}, err
	}

	var customer models.Customer
	err = json.Unmarshal(resp_body, &customer)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func postCustomer(customer models.Customer) (models.Customer, error) {
	bs, err := json.Marshal(customer)
	if err != nil {
		fmt.Println("Error marshaling customer:", err)
		return models.Customer{}, err
	}

	client := &http.Client{
		Timeout: limitTime,
	}

	resp, err := client.Post(config.Backend_url+"/customers/insert", "application/json", bytes.NewBuffer(bs))
	if err != nil {
		fmt.Println("Error sending request")
		return models.Customer{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Customer{}, fmt.Errorf("failed to post customer: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot read resp body")
		return models.Customer{}, err
	}

	var outCustomer models.Customer
	if err = json.Unmarshal(body, &outCustomer); err != nil {
		fmt.Println("Cannot unmarshal response body:", err)
		return models.Customer{}, err
	}

	return outCustomer, nil
}

func updateCustomer(customer models.Customer) (models.Customer, error) {
	bs, err := json.Marshal(customer)
	if err != nil {
		fmt.Println("Error marshaling customer:", err)
		return models.Customer{}, err
	}
	req, err := http.NewRequest(http.MethodPatch, config.Backend_url+"/customers/update", bytes.NewBuffer(bs))
	if err != nil {
		return models.Customer{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: limitTime,
	}
	resp, err := client.Do(req)
	if err != nil {
		return models.Customer{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Customer{}, fmt.Errorf("failed to patch customer: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot read response body:", err)
		return models.Customer{}, err
	}

	var outCustomer models.Customer
	if err = json.Unmarshal(body, &outCustomer); err != nil {
		fmt.Println("Cannot unmarshal response body:", err)
		return models.Customer{}, err
	}

	return outCustomer, nil
}
