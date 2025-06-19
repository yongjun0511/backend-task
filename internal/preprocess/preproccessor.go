package preprocess

import (
	"bufio"
	"log"
	"os"

	"banksalad-backend-task/internal/domain"
	"banksalad-backend-task/internal/filter"
	"banksalad-backend-task/internal/parser"
	"banksalad-backend-task/internal/validator"
)

type Preprocessor struct {
	path      string
	parser    parser.Parser
	validator validator.Validator
	filter    filter.ContactFilter
}

func NewPreprocessor(path string, p parser.Parser, v validator.Validator, f filter.ContactFilter) *Preprocessor {
	return &Preprocessor{
		path:      path,
		parser:    p,
		validator: v,
		filter:    f,
	}
}

func (pp *Preprocessor) Run() (map[string]struct{}, map[string]struct{}, error) {
	f, err := os.Open(pp.path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var records []domain.UserRecord
	sc := bufio.NewScanner(f)
	lineNum := 1

	for sc.Scan() {
		line := sc.Text()

		if ok, reason := pp.validator.ValidateLine(line); !ok {
			log.Printf("[Validation Error] line %d: %s\nReason: %s\n", lineNum, line, reason)
			lineNum++
			continue
		}

		record := pp.parser.ParseLine(line)
		records = append(records, record)
		lineNum++
	}

	if err := sc.Err(); err != nil {
		return nil, nil, err
	}

	emailSet, phoneSet := pp.filter.Extract(records)
	return emailSet, phoneSet, nil
}
