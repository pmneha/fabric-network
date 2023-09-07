package main

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Employee struct {
	Name		string	`json:"name"`
	Email		string	`json:"email"`
	Designation	string	`json:"designation"`
	Mobile		int	`json:"mobile"`
}

func (s *SmartContract) CreateEmployee(ctx contractapi.TransactionContextInterface, employeeData string) (string,
	error) {
		if len(employeeData) == 0 {
			return "", fmt.Errorf("Please pass the correct employee data")
		}

		var employee Employee
		err := json.Unmarshal([]byte(employeeData), &employee)
		if err != nil {
			return "", fmt.Errorf("Failed while unmarshling employee data. %s", err.Error())
		}

		employeeAsBytes, err := json.Marshal(employee)
		if err != nil {
			return "", fmt.Errorf("Failed while marshling employee. %s", err.Error())
		}

		return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(employee.Email, employeeAsBytes)

	}

func (s *SmartContract) UpdateEmployee(ctx contractapi.TransactionContextInterface, email string, name string,
	designation string, mobile int) (string, error) {
		employeeAsBytes, err := ctx.GetStub().GetState(email)
		if err != nil {
			return "", fmt.Errorf("Failed to get employee data. %s", err.Error())
		}

		if employeeAsBytes == nil {
			return "", fmt.Errorf("%s does not exist", email)
		}

		employee := new(Employee)
		_ = json.Unmarshal(employeeAsBytes, employee)

		employee.Name = name
		employee.Designation = designation
		employee.Mobile = mobile

		employeeAsBytes, err = json.Marshal(employee)
		if err != nil {
			return "", fmt.Errorf("Failed while marshling employee. %s", err.Error())
		}

		return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(employee.Email, employeeAsBytes)
	}

func (s *SmartContract) ReadEmployee(ctx contractapi.TransactionContextInterface, email string) (*Employee, error) {
	if len(email) == 0 {
		return nil, fmt.Errorf("Please provide correct email Id")
	}

	employeeAsBytes, err := ctx.GetStub().GetState(email)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if employeeAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", email)
	}

	employee := new(Employee)
	_ = json.Unmarshal(employeeAsBytes, employee)

	return employee, nil
}

func (s *SmartContract) DeleteEmployee(ctx contractapi.TransactionContextInterface, email string) (string, error) {
	if len(email) == 0 {
		return "", fmt.Errorf("Please provide correct email Id")
	}
	return ctx.GetStub().GetTxID(), ctx.GetStub().DelState(email)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		log.Panicf("Error create employee chaincode: %s", err.Error())
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincodes: %s", err.Error())
	}
}
