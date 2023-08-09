package videoeditor

import "fmt"

type GenerateError struct {
	file   string
	step   string
	output string
	err    error
}

func (e *GenerateError) Error() string {
	return fmt.Sprintf("Error while %s: %s", e.step, e.err.Error())
}

func (e *GenerateError) Output() string {
	return e.output
}

func NewGenerateError(step string, err error, output string) *GenerateError {
	return &GenerateError{step: step, err: err, output: output}
}
