package eval

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/picatz/openai"
)

// BasicMatch contains the data for an OpenAI
// "evals.elsuite.basic.match:Match" evaluation.
type BasicMatch struct {
	// The system message to provide context.
	//
	// Example: "You are about to be asked a question. Please answer as concisely as possible."
	SystemContent string `hcl:"system,attr"`

	// The user message containing the instruction.
	//
	// Example: "OpenAI was founded in 20"
	UserContent string `hcl:"user,attr"`

	// The ideal response.
	//
	// Example: "15"
	Ideal string `hcl:"ideal,attr"`
}

// BasicMatches is a slice of EvalBasicMatch that can be marshaled to JSONL.
type BasicMatches []*BasicMatch

// MarshalJSON is a custom JSON marshaler for the BasicMatch struct.
func (e *BasicMatch) MarshalJSON() ([]byte, error) {
	// Using the evaluationTemplate struct to create the JSON
	et := DataTemplate{
		Ideal: e.Ideal,
		Input: []*openai.ChatMessage{
			{
				Role:    "system",
				Content: e.SystemContent,
			},
			{
				Role:    "user",
				Content: e.UserContent,
			},
		},
	}

	return json.Marshal(et)
}

// MarshalJSON is a custom JSON marshaler for the EvalBasicMatches slice,
// printing each evaluation on a new line.
func (bms BasicMatches) WriteFile(path string) error {
	// Open file for writing, truncating the file if it already exists.
	fh, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer fh.Close()

	// Print each evaluation JSON on a new line.
	for _, bm := range bms {
		b, err := bm.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal evaluation: %w", err)
		}

		_, err = fh.Write(b)
		if err != nil {
			return fmt.Errorf("failed to write evaluation: %w", err)
		}

		_, err = fh.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	return nil
}

// ReadHCL reads the HCL file at the given path and returns the
// BasicMatches slice.
func ReadHCL(path string) (BasicMatches, error) {
	// Return value if successful.
	bm := BasicMatches{}

	// Read the file data.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening HCL file: %w", err)
	}

	// Parse the HCL configuration.
	f, diags := hclsyntax.ParseConfig(data, path, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing HCL file: %w", diags)
	}

	// Get file body content.
	c, diags := f.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "eval",
				// TODO: consider adding labels to the eval blocks to group them?
				// LabelNames: []string{
				// 	"identifier",
				// },
			},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing HCL file: %s", err.Error())
	}

	// Iterate over the eval blocks.
	for _, block := range c.Blocks {
		if block.Type != "eval" {
			// TODO: consider erroring here?
			continue
		}

		// Create the evaluation.
		eval := &BasicMatch{}

		// Get block body content.
		bc, diags := block.Body.Content(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: "system",
				},
				{
					Name: "user",
				},
				{
					Name: "ideal",
				},
			},
		})
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing HCL file eval block: %w", diags)
		}

		// Get the system attribute.
		system, diags := bc.Attributes["system"].Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing HCL file eval block: %w", diags)
		}

		// Get the user attribute.
		user, diags := bc.Attributes["user"].Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing HCL file eval block: %w", diags)
		}

		// Get the ideal attribute.
		ideal, diags := bc.Attributes["ideal"].Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing HCL file eval block: %w", diags)
		}

		// Set the evaluation fields.
		eval.SystemContent = system.AsString()
		eval.UserContent = user.AsString()
		eval.Ideal = ideal.AsString()

		// Add the evaluation to the slice.
		bm = append(bm, eval)
	}

	// Successfully parsed the HCL file.
	return bm, nil
}
