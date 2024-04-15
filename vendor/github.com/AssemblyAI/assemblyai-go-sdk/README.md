<img src="https://github.com/AssemblyAI/assemblyai-go-sdk/blob/main/assemblyai.png?raw=true" width="500"/>

---

[![CI Passing](https://github.com/AssemblyAI/assemblyai-go-sdk/actions/workflows/go.yml/badge.svg)](https://github.com/AssemblyAI/assemblyai-go-sdk/actions/workflows/go.yml)
[![GitHub License](https://img.shields.io/github/license/AssemblyAI/assemblyai-go-sdk)](https://github.com/AssemblyAI/assemblyai-go-sdk/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/AssemblyAI/assemblyai-go-sdk.svg)](https://pkg.go.dev/github.com/AssemblyAI/assemblyai-go-sdk)
[![AssemblyAI Twitter](https://img.shields.io/twitter/follow/AssemblyAI?label=%40AssemblyAI&style=social)](https://twitter.com/AssemblyAI)
[![AssemblyAI YouTube](https://img.shields.io/youtube/channel/subscribers/UCtatfZMf-8EkIwASXM4ts0A)](https://www.youtube.com/@AssemblyAI)
[![Discord](https://img.shields.io/discord/875120158014853141?logo=discord&label=Discord&link=https%3A%2F%2Fdiscord.com%2Fchannels%2F875120158014853141&style=social)
](https://assemblyai.com/discord)

# AssemblyAI Go SDK

A Go client library for accessing [AssemblyAI](https://assemblyai.com).

## Overview

- [AssemblyAI Go SDK](#assemblyai-go-sdk)
  - [Overview](#overview)
  - [Documentation](#documentation)
  - [Quickstart](#quickstart)
    - [Installation](#installation)
    - [Examples](#examples)
      - [Core Transcription](#core-transcription)
      - [Audio Intelligence](#audio-intelligence)
      - [Real-Time Transcription](#real-time-transcription)
  - [Playgrounds](#playgrounds)

## Documentation

Visit our [AssemblyAI API Documentation](https://www.assemblyai.com/docs) to get an overview of our models!

See the reference docs at [pkg.go.dev](https://pkg.go.dev/github.com/AssemblyAI/assemblyai-go-sdk).

## Quickstart

### Installation

```bash
go get github.com/AssemblyAI/assemblyai-go-sdk
```

### Examples

Before you begin, you need to have your API key. If you don't have one yet, [**sign up for one**](https://www.assemblyai.com/dashboard/signup)!

#### Core Transcription

<details>
    <summary>Transcribe an audio file from URL</summary>

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/AssemblyAI/assemblyai-go-sdk"
)

func main() {
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")

	ctx := context.Background()

	audioURL := "https://example.org/audio.mp3"

	client := assemblyai.NewClient(apiKey)

	transcript, err := client.Transcripts.TranscribeFromURL(ctx, audioURL, nil)
	if err != nil {
		log.Fatal("Something bad happened:", err)
	}

	log.Println(*transcript.Text)
}
```

</details>
<details>
    <summary>Transcribe a local audio file</summary>

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/AssemblyAI/assemblyai-go-sdk"
)

func main() {
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")

	ctx := context.Background()

	client := assemblyai.NewClient(apiKey)

	f, err := os.Open("./my-local-audio-file.wav")
	if err != nil {
		log.Fatal("Couldn't open audio file:", err)
	}
	defer f.Close()

	transcript, err := client.Transcripts.TranscribeFromReader(ctx, f, nil)
	if err != nil {
		log.Fatal("Something bad happened:", err)
	}

	log.Println(*transcript.Text)
}
```

</details>

#### Audio Intelligence

<details>
    <summary>Identify entities in a transcript</summary>

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/AssemblyAI/assemblyai-go-sdk"
)

func main() {
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")

	ctx := context.Background()

	audioURL := "https://example.org/audio.mp3"

	client := assemblyai.NewClient(apiKey)

	opts := &assemblyai.TranscriptParams{
		EntityDetection: assemblyai.Bool(true),
	}

	transcript, err := client.Transcripts.TranscribeFromURL(ctx, audioURL, opts)
	if err != nil {
		log.Fatal("Something bad happened:", err)
	}

	for _, entity := range transcript.Entities {
		log.Println(*entity.Text)
		log.Println(entity.EntityType)
		log.Printf("Timestamp: %v - %v", *entity.Start, *entity.End)
	}
}
```

</details>

#### Real-Time Transcription

Check out the [realtime](./examples/realtime) example.

## Playgrounds

Visit one of our Playgrounds:

- [LeMUR Playground](https://www.assemblyai.com/playground/v2/source)
- [Transcription Playground](https://www.assemblyai.com/playground)
