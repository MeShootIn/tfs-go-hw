/*
Package bill provides work with files of reports on financial transactions of companies (.json format).
*/
package bill

import (
	"math"
	"strconv"
	"time"
	"unicode/utf8"
)

// UnsupportedBill error says that this bill cannot be processed and must be skipped.
type UnsupportedBill struct {
	msg string
}

func (e UnsupportedBill) Error() string {
	return e.msg
}

// InvalidBill error says that this bill contains invalid fields, but can be processed.
type InvalidBill struct {
	msg string
}

func (e InvalidBill) Error() string {
	return e.msg
}

// worstError returns the worst error of the two passed in the following priority: UnsupportedBill, InvalidBill, nil.
func worstError(e1, e2 error) error {
	if _, ok := e1.(UnsupportedBill); ok {
		return e1
	}

	if e2 == nil {
		return e1
	}

	return e2
}

type (
	Type      interface{}
	Value     interface{}
	ID        interface{}
	CreatedAt interface{}
	Company   interface{}
)

type Body struct {
	Type  *Type  `json:"type"`
	Value *Value `json:"value"`
	ID    *ID    `json:"id"`
}

type Operation struct {
	*Body
	CreatedAt *CreatedAt `json:"created_at"`
}

type Bill struct {
	Company *Company `json:"company"`
	*Operation
	OperationStruct *Operation `json:"operation"`
}

// isInteger checks if the given float64 number is an integer.
func isInteger(f float64) bool {
	return math.Abs(f-float64(int(f))) <= 1e-16
}

// parseType checks the Type for validity and returns a pointer to a value with an error.
func parseType(opType *Type) (*string, error) {
	if opType == nil {
		return nil, InvalidBill{msg: "operation type was not passed"}
	}

	invalidType := InvalidBill{
		msg: "operation type can only take one of the values: \"income\", \"outcome\", \"+\", \"-\"",
	}

	switch (*opType).(type) {
	case string:
		str := (*opType).(string)

		if !(str == "income" || str == "outcome" || str == "+" || str == "-") {
			return nil, invalidType
		}

		return &str, nil
	default:
		return nil, invalidType
	}
}

// parseValue checks the Value for validity and returns a pointer to a value of one of the possible types with an error.
func parseValue(value *Value) (*float64, *string, error) {
	if value == nil {
		return nil, nil, InvalidBill{msg: "operation value was not passed"}
	}

	invalidValue := InvalidBill{
		msg: "operation value can only be of the following types: int, float (always integer), string (always integer)",
	}

	switch (*value).(type) {
	case float64:
		flt := (*value).(float64)

		if !isInteger(flt) {
			return nil, nil, invalidValue
		}

		return &flt, nil, nil
	case string:
		str := (*value).(string)

		if _, err := strconv.Atoi(str); err != nil {
			return nil, nil, invalidValue
		}

		return nil, &str, nil
	default:
		return nil, nil, invalidValue
	}
}

// parseID checks the ID for validity and returns a pointer to a value of one of the possible types with an error.
func parseID(id *ID) (*string, *float64, error) {
	if id == nil {
		return nil, nil, UnsupportedBill{msg: "operation id was not passed"}
	}

	invalidID := UnsupportedBill{msg: "operation id can only be of type int and string"}

	switch (*id).(type) {
	case string:
		str := (*id).(string)

		if utf8.RuneCountInString(str) == 0 {
			return nil, nil, invalidID
		}

		return &str, nil, nil
	case float64:
		flt := (*id).(float64)

		if !isInteger(flt) {
			return nil, nil, invalidID
		}

		return nil, &flt, nil
	default:
		return nil, nil, invalidID
	}
}

// parseCompany checks the Company for validity and returns a pointer to a value with an error.
func parseCompany(company *Company) (*string, error) {
	if company == nil {
		return nil, UnsupportedBill{msg: "company was not passed"}
	}

	invalidCompany := UnsupportedBill{msg: "company name can only be a non-empty string"}

	switch (*company).(type) {
	case string:
		str := (*company).(string)

		if utf8.RuneCountInString(str) == 0 {
			return nil, invalidCompany
		}

		return &str, nil
	default:
		return nil, invalidCompany
	}
}

// parseCreatedAt checks the CreatedAt for validity and returns a pointer to a value with an error.
func parseCreatedAt(createdAt *CreatedAt) (*string, error) {
	if createdAt == nil {
		return nil, UnsupportedBill{msg: "the \"created_at\" field was not passed"}
	}

	invalidCreatedAt := UnsupportedBill{msg: "time must be a string in \"RFC3339\" format"}

	switch (*createdAt).(type) {
	case string:
		str := (*createdAt).(string)
		str = FormatQuotes(str)

		if _, err := time.Parse(time.RFC3339, str); err != nil {
			return nil, invalidCreatedAt
		}

		return &str, nil
	default:
		return nil, invalidCreatedAt
	}
}

// checkBody checks the Body for validity.
func checkBody(body *Body) error {
	if body == nil {
		return UnsupportedBill{msg: "operation body was not passed"}
	}

	var err error

	_, _, errID := parseID(body.ID)
	_, errType := parseType(body.Type)
	_, _, errValue := parseValue(body.Value)

	for _, e := range [...]error{errID, errType, errValue} {
		err = worstError(err, e)
	}

	return err
}

// checkOperation checks the embedded structure or field of Operation for validity.
func checkOperation(bill Bill) error {
	operationRoot := bill.Operation
	operationStruct := bill.OperationStruct
	var (
		body      *Body
		createdAt *CreatedAt
		err       error
	)

	switch {
	case operationRoot != nil && operationStruct == nil:
		body = operationRoot.Body
		createdAt = operationRoot.CreatedAt
	case operationRoot == nil && operationStruct != nil:
		body = operationStruct.Body
		createdAt = operationStruct.CreatedAt
	case operationRoot != nil && operationStruct != nil:
		switch {
		case operationRoot.Body != nil && operationStruct.Body == nil:
			body = operationRoot.Body

			switch {
			case operationRoot.CreatedAt == nil && operationStruct.CreatedAt != nil:
				createdAt = operationStruct.CreatedAt
			default:
				return UnsupportedBill{msg: "one field \"created_at\" must be set"}
			}
		case operationRoot.Body == nil && operationStruct.Body != nil:
			body = operationStruct.Body

			switch {
			case operationRoot.CreatedAt != nil && operationStruct.CreatedAt == nil:
				createdAt = operationRoot.CreatedAt
			default:
				return UnsupportedBill{msg: "one field \"created_at\" must be set"}
			}
		default:
			return UnsupportedBill{msg: "one operation body must be set"}
		}
	default:
		return UnsupportedBill{msg: "operation was not passed"}
	}

	_, errCreatedAt := parseCreatedAt(createdAt)
	errBody := checkBody(body)

	for _, e := range [...]error{errCreatedAt, errBody} {
		err = worstError(err, e)
	}

	return err
}

// checkBill checks the Bill for validity.
func checkBill(bill Bill) error {
	var err error

	_, errCompany := parseCompany(bill.Company)
	errOperation := checkOperation(bill)

	for _, e := range [...]error{errCompany, errOperation} {
		err = worstError(err, e)
	}

	return err
}
