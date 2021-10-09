package bill

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type InfoInvalid struct {
	CreatedAt int64
	ID        ID
}

type Report struct {
	Company              string      `json:"company"`
	ValidOperationsCount uint        `json:"valid_operations_count"`
	Balance              int         `json:"balance"`
	InvalidOperations    interface{} `json:"invalid_operations,omitempty"`
}

// ProcessBills parses the Bill slice and writes the report to a file (.json format).
func ProcessBills(bills []Bill, outputFileName string) error {
	output, err := os.Create(outputFileName)
	defer DeferClose(output)

	if err != nil {
		return err
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "\t")
	reportMap := map[string]Report{}

	for _, bill := range bills {
		err = checkBill(bill)

		switch err.(type) {
		case UnsupportedBill:
			continue
		case InvalidBill:
			var createdAt *CreatedAt

			if bill.Operation != nil && bill.Operation.CreatedAt != nil {
				createdAt = bill.Operation.CreatedAt
			}

			if bill.OperationStruct != nil && bill.OperationStruct.CreatedAt != nil {
				createdAt = bill.OperationStruct.CreatedAt
			}

			timeStr, _ := parseCreatedAt(createdAt)
			timeTime, _ := time.Parse(time.RFC3339, *timeStr)
			unixTime := timeTime.Unix()

			var id *ID

			if bill.Operation != nil && bill.Operation.Body != nil && bill.Operation.Body.ID != nil {
				id = bill.Operation.Body.ID
			}

			if bill.OperationStruct != nil && bill.OperationStruct.Body != nil && bill.OperationStruct.Body.ID != nil {
				id = bill.OperationStruct.Body.ID
			}

			companyPtr, _ := parseCompany(bill.Company)
			company := *companyPtr
			report, ok := reportMap[company]

			if !ok {
				report = Report{
					Company:           company,
					InvalidOperations: []InfoInvalid{},
				}
			}

			report.InvalidOperations = append((report.InvalidOperations).([]InfoInvalid), InfoInvalid{
				CreatedAt: unixTime,
				ID:        *id,
			})
			reportMap[company] = report
		case nil:
			var (
				tp    *Type
				value *Value
			)

			if bill.Operation != nil && bill.Operation.Body != nil {
				tp = bill.Operation.Body.Type
				value = bill.Operation.Body.Value
			}

			if bill.OperationStruct != nil && bill.OperationStruct.Body != nil && bill.OperationStruct.Body.Value != nil {
				tp = bill.OperationStruct.Body.Type
				value = bill.OperationStruct.Body.Value
			}

			var profit int
			valueFlt, _, _ := parseValue(value)

			if valueFlt != nil {
				valueStr := fmt.Sprintf("%.0f", *valueFlt)
				profit, _ = strconv.Atoi(valueStr)
			} else {
				_, valueStr, _ := parseValue(value)
				profit, _ = strconv.Atoi(*valueStr)
			}

			operationTypePtr, _ := parseType(tp)
			operationType := *operationTypePtr

			if operationType == "outcome" || operationType == "-" {
				profit *= -1
			}

			companyPtr, _ := parseCompany(bill.Company)
			company := *companyPtr
			report, ok := reportMap[company]

			if !ok {
				report = Report{
					Company:           company,
					InvalidOperations: []InfoInvalid{},
				}
			}

			report.ValidOperationsCount++
			report.Balance += profit
			reportMap[company] = report
		}
	}

	for company := range reportMap {
		if reportMap[company].InvalidOperations == nil {
			continue
		}

		sort.Slice(reportMap[company].InvalidOperations, func(i, j int) bool {
			invalidOperations := (reportMap[company].InvalidOperations).([]InfoInvalid)

			return invalidOperations[i].CreatedAt < invalidOperations[j].CreatedAt
		})

		invalidOperations := (reportMap[company].InvalidOperations).([]InfoInvalid)
		ids := make([]ID, 0, len(invalidOperations))

		for _, invalidOperation := range invalidOperations {
			ids = append(ids, invalidOperation.ID)
		}

		reportMap[company] = Report{
			Company:              company,
			ValidOperationsCount: reportMap[company].ValidOperationsCount,
			Balance:              reportMap[company].Balance,
			InvalidOperations:    ids,
		}
	}

	reports := make([]Report, 0, len(reportMap))

	for _, report := range reportMap {
		reports = append(reports, report)
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Company < reports[j].Company
	})

	return encoder.Encode(reports)
}
