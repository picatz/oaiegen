# oaiegen

WIP tool for generating [OpenAI evals](https://github.com/openai/evals).

## Install

```console
$ go install github.com/picatz/oaiegen
```

## Usage

Use HCL to define your OpenAI evals (currently only `evals.elsuite.basic.match:Match` is really supported):

```hcl
# Example evals.hcl file contents.

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
```

Then run `oaiegen` to generate the evals:

```console
$ oaiegen -workdir="/tmp/oaiegen" -file="evals.hcl"
```

HCL formatted evals are converted to JSONL formatted samples and written to the working directory:

```console
$ ls /tmp/oaiegen
samples.jsonl
$ cat /tmp/oaiegen/samples.jsonl
{"input":[{"role":"system","content":"You are about to be asked a question. Please answer as concisely as possible."},{"role":"user","content":"OpenAI was founded in 20"}],"ideal":"15"}
{"input":[{"role":"system","content":"You are about to be asked a question. Please answer as concisely as possible."},{"role":"user","content":"Once upon a "}],"ideal":"time"}
```

> **Note**: you will still need to manually perform other steps, but this is a start. See https://github.com/openai/evals/pull/1 for more details.