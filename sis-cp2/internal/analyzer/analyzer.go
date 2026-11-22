package analyzer

import (
	"strings"
	"time"
	"sis-cp2/internal/domain"
)

type Processor interface {
	Process(input string) (domain.TaskResult, error)
}

type TextProcessor struct{}

func (p *TextProcessor) Process(input string) (domain.TaskResult, error) {
	if strings.TrimSpace(input) == "" {
		return domain.TaskResult{}, domain.ErrInvalidInput
	}
	
	// Simulate some processing time for concurrency demonstration
	time.Sleep(10 * time.Millisecond)
	
	words := len(strings.Fields(input))
	chars := len(input)
	return domain.TaskResult{WordCount: words, CharCount: chars}, nil
}
