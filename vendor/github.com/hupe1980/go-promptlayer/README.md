# 🍰 go-promptlayer
![Build Status](https://github.com/hupe1980/go-promptlayer/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/go-promptlayer.svg)](https://pkg.go.dev/github.com/hupe1980/go-promptlayer)
> The Go PromptLayer API client enables seamless integration of the PromptLayer platform in your Go projects. With this client, you can effortlessly incorporate PromptLayer's features and streamline your prompt engineering workflow.

## Installation
```
go get github.com/hupe1980/go-promptlayer
```

## Example Usage
Here's a quick example of how you can use the Go PromptLayer API client:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hupe1980/go-promptlayer"
)

func main() {
	client := promptlayer.NewClient(os.Getenv("PROMPTLAYER_API_KEY"))

	startTime := time.Now()
	endTime := startTime.Add(3 * time.Second)

	output, err := client.TrackRequest(context.Background(), &promptlayer.TrackRequestInput{
		FunctionName: "openai.Completion.create",
		// kwargs will need messages if using chat-based completion
		Kwargs: map[string]any{
			"engine": "text-ada-001",
			"prompt": "My name is",
		},
		Tags: []string{"hello", "world"},
		RequestResponse: map[string]any{
			"id":      "cmpl-6TEeJCRVlqQSQqhD8CYKd1HdCcFxM",
			"object":  "text_completion",
			"created": 1672425843,
			"model":   "text-ada-001",
			"choices": []map[string]any{
				{
					"text":          " advocacy\"\n\nMy name is advocacy.",
					"index":         0,
					"logprobs":      nil,
					"finish_reason": "stop",
				},
			},
		},
		RequestStartTime: startTime,
		RequestEndTime:   endTime,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ID:", output.RequestID)
}
```
Output:
```
ID: 6368262
```
For more example usage, see [_examples](./_examples).

## License
[MIT](LICENCE)
