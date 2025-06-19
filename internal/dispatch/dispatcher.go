package dispatch

import (
	"fmt"
	"path/filepath"
	"strings"

	"banksalad-backend-task/internal/parser"
	"banksalad-backend-task/internal/preprocess"
	"banksalad-backend-task/internal/validator"
)

func Resolve(path string) (*preprocess.Preprocessor, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".txt":
		return preprocess.NewPreprocessor(
			path,
			&parser.TxtParser{},
			&validator.TxtValidator{},
		), nil
	default:
		return nil, fmt.Errorf("지원하지 않는 확장자입니다: %s", path)
	}
}
