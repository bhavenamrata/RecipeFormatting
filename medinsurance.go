package main

import (
        "encoding/json"
        "fmt"
        "strconv"
        "strings"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)

// MedInsuranceChaincode example simple Chaincode implementation
type MedInsuranceChaincode struct {
}

type PatientProfile struct {
        PatientID  string    `json:"patientid"`
        PatientName  string `json:"patientname"`
        DateOfBirth string `json:"dob"`
}

type STATUS_OF_CASE int

const (
  CASE_CREATED STATUS_OF_CASE = iota
  APPROVAL_PENDING
  APPROVAL_ACCEPTED
)


type CaseInfo struct {
        CaseId  string    `json:"caseid"`
        PatientInfo  string    `json:"patientinfo"`
        DateOfExamination   string `json:"doe"`
        MedicalReport string `json:"medreport"`
        CostOfService int          `json:"costofservice"`
        Status STATUS_OF_CASE `json:"status"`
}

// ===================================================================================
func main() {
        err := shim.Start(new(MedInsuranceChaincode))
        if err != nil {
                fmt.Printf("Error starting Simple chaincode: %s", err)
        }
}

// Init initializes chaincode
// ===========================
func (t *MedInsuranceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
        return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *MedInsuranceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        function, args := stub.GetFunctionAndParameters()
        fmt.Println("invoke is running " + function)

        // Handle different functions
        if function == "AddPatient" { //create a new patient
                return t.addPatient(stub, args)
        } else if function == "AddCase" { //create a new case
                return t.addCase(stub, args)
        } else if function == "ReadPatient" { //query patient data
                return t.readPatient(stub, args)
        } else if function == "ReadCase" { //Query case data
                return t.readCase(stub, args)
		} else if function == "ClaimMoney" { //Query case data
                return t.claimMoney(stub, args)
		} else if function == "ApproveRequest" { //Query case data
                return t.approveRequest(stub, args)
        }
        fmt.Println("invoke did not find func: " + function) //error
        return shim.Error("Received unknown function invocation")
}

// ============================================================
// addPatient - create a new Patient, store into chaincode state
// ============================================================
func (t *MedInsuranceChaincode) addPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        //   0       1                   2
        // "1", "bob parker", "12Sep1967"
        if len(args) != 3 {
                return shim.Error("Incorrect number of arguments. Expecting 4")
        }

        // ==== Input sanitation ====
        fmt.Println("- start add patient")
        if len(args[0]) <= 0 {
                return shim.Error("1st argument must be a non-empty string")
        }
        if len(args[1]) <= 0 {
                return shim.Error("2nd argument must be a non-empty string")
        }
        if len(args[2]) <= 0 {
                return shim.Error("3rd argument must be a non-empty string")
        }


        patientID := args[0]
        patientName := strings.ToLower(args[1])
        DOB := strings.ToLower(args[2])

        // ==== Check if patient already exists ====
        patientAsBytes, err := stub.GetState(patientName)
        if err != nil {
                return shim.Error("Failed to get patient: " + err.Error())
        } else if patientAsBytes != nil {
                fmt.Println("This patient already exists: " + patientName)
                return shim.Error("This patient already exists: " + patientName)
        }


        // ==== Create patient object and marshal to JSON ====
        patient := &PatientProfile{patientID, patientName, DOB}
        patientJSONasBytes, err := json.Marshal(patient)
        if err != nil {
                return shim.Error(err.Error())
                fmt.Println("Error marshalling patient")
        }

        // === Save patient to state ===
        err = stub.PutState(patientName, patientJSONasBytes)
        if err != nil {
                return shim.Error(err.Error())
        }
        fmt.Println("- end init patient")
        return shim.Success(nil)

}

// ============================================================
// addCase - create a new Patient, store into chaincode state
// ============================================================
func (t *MedInsuranceChaincode) addCase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        //   0       1          2          3                            4
        // "3abc", "1",         "12Sep2017", "medicalReport"   "4000"
        if len(args) != 5 {
                return shim.Error("Incorrect number of arguments. Expecting 5")
        }

        // ==== Input sanitation ====
        fmt.Println("- start add case")
        if len(args[0]) <= 0 {
                return shim.Error("1st argument must be a non-empty string")
        }
        if len(args[1]) <= 0 {
                return shim.Error("2nd argument must be a non-empty string")
        }
        if len(args[2]) <= 0 {
                return shim.Error("3rd argument must be a non-empty string")
        }
        if len(args[3]) <= 0 {
                return shim.Error("4th argument must be a non-empty string")
        }
        if len(args[4]) <= 0 {
                return shim.Error("5th argument must be a non-empty string")
        }

        caseId := args[0]
        patientID := args[1]
        dateOfExamination := strings.ToLower(args[2])
        medicalReport := args[3]
        costOfService, _ := strconv.Atoi(args[4])
        var status STATUS_OF_CASE = CASE_CREATED

        // ==== Create case object and marshal to JSON ====
        caseinfo := &CaseInfo{caseId, patientID, dateOfExamination, medicalReport, costOfService, status}
        caseJSONasBytes, err := json.Marshal(caseinfo)
        if err != nil {
                return shim.Error(err.Error())
                fmt.Println("Error marshalling case")
        }
        // === Save patient to state ===
        err = stub.PutState(caseId, caseJSONasBytes)
        if err != nil {
                return shim.Error(err.Error())
        }
        fmt.Println("- end init patient")
        return shim.Success(nil)
}


// ======================================================
// readPatient - read a patient from chaincode state
// ======================================================
func (t *MedInsuranceChaincode) readPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var patientName, jsonResp string
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting name of the patient to query")
        }

        patientName = args[0]
        valAsbytes, err := stub.GetState(patientName) //get the patient from chaincode state
		fmt.Println(valAsbytes)
        if err != nil {
                jsonResp = "{\"Error\":\"Failed to get state for " + patientName + "\"}"
                fmt.Println("Failed to get state for patient")
                return shim.Error(jsonResp)
        } else if valAsbytes == nil {
                jsonResp = "{\"Error\":\"patient does not exist: " + patientName + "\"}"
                fmt.Println("patient does not exist")
                return shim.Error(jsonResp)
        }

        return shim.Success(valAsbytes)
}

// ======================================================
// readCase - read a case from chaincode state
// ======================================================
func (t *MedInsuranceChaincode) readCase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var caseID, jsonResp string
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting case number to query")
        }

        caseID = args[0]
        valAsbytes, err := stub.GetState(caseID) //get the case from chaincode state
        if err != nil {
                jsonResp = "{\"Error\":\"Failed to get state for " + caseID + "\"}"
                fmt.Println("Failed to get state for case")
                return shim.Error(jsonResp)
        } else if valAsbytes == nil {
                jsonResp = "{\"Error\":\"case does not exist: " + caseID + "\"}"
                fmt.Println("case does not exist")
                return shim.Error(jsonResp)
        }

        return shim.Success(valAsbytes)
}


// ===============================================================
// ClaimMoney - Request for claiming money via passing the caseID
// ===============================================================
func (t *MedInsuranceChaincode) claimMoney(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var caseID, jsonResp string
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting case ID to claim")
        }

        caseID = args[0]
		
        valAsbytes, err := stub.GetState(caseID) //get the case from chaincode state
        if err != nil {
                jsonResp = "{\"Error\":\"Failed to get state for " + caseID + "\"}"
                fmt.Println("Failed to get state for case")
                return shim.Error(jsonResp)
        } else if valAsbytes == nil {
                jsonResp = "{\"Error\":\"case does not exist: " + caseID + "\"}"
                fmt.Println("case does not exist")
                return shim.Error(jsonResp)
        }
		
		var newcase CaseInfo
		err = json.Unmarshal(valAsbytes, &newcase)
	    newcase.Status = APPROVAL_PENDING
		
		updatedcaseAsBytes, _ := json.Marshal(&newcase)
		err = stub.PutState(caseID, updatedcaseAsBytes)
		if err != nil {
			fmt.Println("Error storing Case back")
			return shim.Error("Error storing Case back")
		}	
		fmt.Println("Updated case with Id " + caseID )
		
		// create event to notify 
        eventtonotify := "caseID:" + newcase.CaseId
        stub.SetEvent("provideapproval", []byte(eventtonotify))
        return shim.Success(nil)
		
        return shim.Success(updatedcaseAsBytes)
}


// ===============================================================
// approveRequest - Approve the pending claim request.
// ===============================================================
func (t *MedInsuranceChaincode) approveRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var caseID, jsonResp string
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting case ID to claim")
        }

        caseID = args[0]
		
        valAsbytes, err := stub.GetState(caseID) //get the case from chaincode state
        if err != nil {
                jsonResp = "{\"Error\":\"Failed to get state for " + caseID + "\"}"
                fmt.Println("Failed to get state for case")
                return shim.Error(jsonResp)
        } else if valAsbytes == nil {
                jsonResp = "{\"Error\":\"case does not exist: " + caseID + "\"}"
                fmt.Println("case does not exist")
                return shim.Error(jsonResp)
        }
		
		var newcase CaseInfo
		err = json.Unmarshal(valAsbytes, &newcase)
	    newcase.Status = APPROVAL_ACCEPTED
		
		updatedcaseAsBytes, _ := json.Marshal(&newcase)
		err = stub.PutState(caseID, updatedcaseAsBytes)
		if err != nil {
			fmt.Println("Error storing Case back")
			return shim.Error("Error storing Case back")
		}	
		fmt.Println("Updated case with Id " + caseID )
		
		// create event to notify 
        eventtonotify := "caseID:" + newcase.CaseId
        stub.SetEvent("provideapproval", []byte(eventtonotify))
        return shim.Success(nil)
		
        return shim.Success(updatedcaseAsBytes)
}


