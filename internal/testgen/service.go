package testgen

import (
	"github.com/ak95asb/dsa-dojo/internal/problem"
)

// Service handles test case generation operations
type Service struct {
	interactive *InteractiveInput
	jsonImport  *JSONImporter
	generator   *Generator
}

// NewService creates a new testgen service
func NewService() *Service {
	return &Service{
		interactive: NewInteractiveInput(),
		jsonImport:  NewJSONImporter(),
		generator:   NewGenerator(),
	}
}

// GenerateInteractive generates test cases from interactive user input
func (s *Service) GenerateInteractive(prob *problem.ProblemDetails, append bool) error {
	// Get test cases from interactive input
	testCases, err := s.interactive.Collect()
	if err != nil {
		return err
	}

	// Generate test file
	return s.generator.Generate(prob, testCases, append)
}

// GenerateFromFile generates test cases from JSON file
func (s *Service) GenerateFromFile(prob *problem.ProblemDetails, filePath string, append bool) error {
	// Import test cases from JSON file
	testCases, err := s.jsonImport.ImportFromFile(filePath)
	if err != nil {
		return err
	}

	// Generate test file
	return s.generator.Generate(prob, testCases, append)
}

// TestCase represents a single test case
type TestCase struct {
	Name     string
	Inputs   []interface{}
	Expected interface{}
}
