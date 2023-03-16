package eval_test

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/picatz/oaiegen/pkg/eval"
)

func TestReadHCL(t *testing.T) {
	// Create the HCL file.
	hcl := `
eval {
	system = "You are about to be asked a question. Please answer as concisely as possible."
	user   = "OpenAI was founded in 20"
	ideal  = "15"
}

eval {
	system = "You are about to be asked a question. Please answer as concisely as possible."
	user   = "Once upon a "
	ideal  = "time"
}
`

	// Create the file.
	f, err := ioutil.TempFile(t.TempDir(), "test*.hcl")
	if err != nil {
		t.Fatal(err)
	}

	// Write the HCL to the file.
	_, err = f.Write([]byte(hcl))
	if err != nil {
		t.Fatal(err)
	}

	// Close the file.
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Read the HCL file.
	ebm, err := eval.ReadHCL(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	testJSONLFile := filepath.Join(t.TempDir(), "test.jsonl")

	// Write the JSONL file.
	err = ebm.WriteFile(testJSONLFile)
	if err != nil {
		t.Fatal(err)
	}

	// Print the JSONL file.
	fb, err := os.ReadFile(testJSONLFile)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(fb))

	// Verify there are two JSONL objects.
	bs := bufio.NewScanner(bytes.NewReader(fb))
	bs.Split(bufio.ScanLines)
	var count int

	// TODO: check the JSONL objects directly.
	for bs.Scan() {
		count++
	}

	if count != 2 {
		t.Fatalf("expected 2 JSONL objects, got %d", count)
	}
}
