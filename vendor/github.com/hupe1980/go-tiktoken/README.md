# ✂️ go-tiktoken
![Build Status](https://github.com/hupe1980/go-tiktoken/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/go-tiktoken.svg)](https://pkg.go.dev/github.com/hupe1980/go-tiktoken)
> OpenAI's [tiktoken](https://github.com/openai/tiktoken) tokenizer written in Go. The vocabularies are embedded and do not need to be downloaded at runtime.

## Installation
```
go get github.com/hupe1980/go-tiktoken
```

## How to use
```golang
package main

import (
	"fmt"
	"log"

	"github.com/hupe1980/go-tiktoken"
)

func main() {
	encoding, err := tiktoken.NewEncodingForModel("gpt-3.5-turbo")
	if err != nil {
		log.Fatal(err)
	}

	ids, tokens, err := encoding.Encode("Hello World", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("IDs:", ids)
	fmt.Println("Tokens:", tokens)
}
```
Output:
```text
IDs: [9906 4435]
Tokens: [Hello  World]
```

For more example usage, see [_examples](./_examples).

## Supported Encodings
- ✅ cl100k_base
- ✅ p50k_base
- ✅ p50k_edit
- ✅ r50k_base
- ✅ gpt2 
- ✅ claude 

## License
[MIT](LICENCE)