package main

import (
	"errors"
	"fmt"
    "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) ([]byte, error) {
	_, args := stub.GetFunctionAndParameters()
    
    if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("init_parameter", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is ur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) ([]byte, error) {
	
    function, args := stub.GetFunctionAndParameters()
    
    fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub)
	} else if function == "campaign" {
		return t.campaign(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	}else if function == "CreateUser" {
	  return t.CreateUser(stub,args)
    }
	
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) ([]byte, error) {
	function, args := stub.GetFunctionAndParameters()
    
    fmt.Println("query is running " + function)

	// Handle different functions
	if function == "get" { //read a variable
		return t.get(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}





// CreateUser - invoke function to register new user
func (t *SimpleChaincode) CreateUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, ifcampaign string
	var err error
	var value string
	fmt.Println("running CreateUser()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //email for user, campaignname for the campaign
	ifcampaign = args[1]  // true if it is a campaign
	
	if ifcampaign == "true"{
		value = strconv.Itoa(0)
	}else{
		value = strconv.Itoa(500)
	}
	
	err = stub.PutState(key, []byte(value)) //associate the balance with user/campaign
	if err != nil {
		return nil, err
	}
	return nil, nil
}


// campaign - invoke function to create new campaign
func (t *SimpleChaincode) campaign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key  string
	var err error
	var value string
	fmt.Println("running campaign()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //email for user, campaignname for the campaign
	value = args[1]  // true if it is a campaign
	
	err = stub.PutState(key, []byte(value)) //associate the balance with user/campaign
	if err != nil {
		return nil, err
	}
	return nil, nil
}


// update - invoke function to update balance of campaign/user after investment
func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, action string
	var err error
	var temp, amount int
	fmt.Println("running update()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //email for user, campaignname for the campaign
	action = args[1]  // possible values : add/subtract
	amount, err = strconv.Atoi(args[2])  // amount to add/subtract, converted to int from string
	
	getvalue, err := stub.GetState(key)
	temp, err = strconv.Atoi(string(getvalue))   // converted to int from string
	if err != nil {
		return nil, err
	}
    
	if action == "add"{
		temp = temp + amount  //add the balance in campaign
	}else if action == "subtract"{
		temp = temp - amount //subtract the balance of user after investment
	}
	
	setvalue := strconv.Itoa(temp)  // converted to string form int
	err = stub.PutState(key, []byte(setvalue))
	if err != nil {
		return nil, err
	}
	return nil, nil
}


// get - query function to read key/value pair
func (t *SimpleChaincode) get(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
