package markdown_basic

import (
	"fmt"
	lcgosplitter "github.com/tmc/langchaingo/textsplitter"
	"reflect"
	"strings"
	"unicode/utf8"

	"gitlab.com/golang-commonmark/markdown"
)

// NewMarkdownTextSplitter creates a new Markdown text splitter.
func NewMarkdownTextSplitter(opts ...Option) *MarkdownTextSplitter {
	options := DefaultOptions()

	sp := &MarkdownTextSplitter{
		ChunkSize:      options.ChunkSize,
		ChunkOverlap:   options.ChunkOverlap,
		SecondSplitter: options.SecondSplitter,
	}

	if sp.SecondSplitter == nil {
		sp.SecondSplitter = lcgosplitter.NewRecursiveCharacter(
			lcgosplitter.WithChunkSize(options.ChunkSize),
			lcgosplitter.WithChunkOverlap(options.ChunkOverlap),
			lcgosplitter.WithSeparators([]string{
				"\n\n", // new line
				"\n",   // new line
				" ",    // space
			}),
		)
	}

	return sp
}

// MarkdownTextSplitter markdown header text splitter.
//
// If your origin document is HTML, you purify and convert to markdown,
// then split it.
type MarkdownTextSplitter struct {
	ChunkSize    int
	ChunkOverlap int
	// SecondSplitter splits paragraphs
	SecondSplitter   lcgosplitter.TextSplitter
	CodeBlocks       bool
	ReferenceLinks   bool
	HeadingHierarchy bool

	MaxHeadingLevel int

	SplitOnListItems    bool
	SplitOnHeadingsOnly bool
}

// SplitText splits a text into multiple text.
func (sp MarkdownTextSplitter) SplitText(text string) ([]string, error) {
	mdParser := markdown.New(markdown.XHTMLOutput(true))
	tokens := mdParser.Parse([]byte(text))

	mc := &markdownContext{
		startAt:        0,
		endAt:          len(tokens),
		tokens:         tokens,
		chunkSize:      sp.ChunkSize,
		chunkOverlap:   sp.ChunkOverlap,
		secondSplitter: sp.SecondSplitter,
		hTitleStack:    []string{},
	}

	chunks := mc.splitText()

	return chunks, nil
}

// markdownContext the helper.
type markdownContext struct {
	// startAt represents the start position of the cursor in tokens
	startAt int
	// endAt represents the end position of the cursor in tokens
	endAt int
	// tokens represents the markdown tokens
	tokens []markdown.Token

	// hTitle represents the current header(H1、H2 etc.) content
	hTitle string
	// hTitleStack represents the hierarchy of headers
	hTitleStack []string
	// hTitlePrepended represents whether hTitle has been appended to chunks
	hTitlePrepended bool

	// chunks represents the final chunks
	chunks []string
	// curSnippet represents the current short markdown-format chunk
	curSnippet string
	// chunkSize represents the max chunk size, when exceeds, it will be split again
	chunkSize int
	// chunkOverlap represents the overlap size for each chunk
	chunkOverlap int

	// secondSplitter re-split markdown single long paragraph into chunks
	secondSplitter lcgosplitter.TextSplitter
}

// splitText splits Markdown text.
//
//nolint:cyclop
func (mc *markdownContext) splitText() []string {
	for idx := mc.startAt; idx < mc.endAt; {
		token := mc.tokens[idx]
		switch token.(type) {
		case *markdown.HeadingOpen:
			mc.onMDHeader()
		default:
			mc.startAt = indexOfCloseTag(mc.tokens, idx) + 1
		}

		idx = mc.startAt
	}

	// apply the last chunk
	mc.applyToChunks()

	return mc.chunks
}

// clone clones the markdownContext with sub tokens.
func (mc *markdownContext) clone(startAt, endAt int) *markdownContext {
	subTokens := mc.tokens[startAt : endAt+1]
	return &markdownContext{
		endAt:  len(subTokens),
		tokens: subTokens,

		hTitle:          mc.hTitle,
		hTitleStack:     mc.hTitleStack,
		hTitlePrepended: mc.hTitlePrepended,

		chunkSize:      mc.chunkSize,
		chunkOverlap:   mc.chunkOverlap,
		secondSplitter: mc.secondSplitter,
	}
}

// onMDHeader splits H1/H2/.../H6
//
// format: HeadingOpen/Inline/HeadingClose
func (mc *markdownContext) onMDHeader() {
	endAt := indexOfCloseTag(mc.tokens, mc.startAt)
	defer func() {
		mc.startAt = endAt + 1
	}()

	header, ok := mc.tokens[mc.startAt].(*markdown.HeadingOpen)
	if !ok {
		return
	}

	// check next token is Inline
	inline, ok := mc.tokens[mc.startAt+1].(*markdown.Inline)
	if !ok {
		return
	}

	title := fmt.Sprintf("%s %s", repeatString(header.HLevel, "#"), inline.Content)

	mc.applyToChunks() // change header, apply to chunks

	mc.hTitle = title

	// fill titlestack with empty strings up to the current level
	for len(mc.hTitleStack) < header.HLevel {
		mc.hTitleStack = append(mc.hTitleStack, "")
	}

	// Build the new title from the title stack, joined by newlines, while ignoring empty entries
	mc.hTitleStack = append(mc.hTitleStack[:header.HLevel-1], mc.hTitle)
	mc.hTitle = ""
	for _, t := range mc.hTitleStack {
		if t != "" {
			mc.hTitle = strings.Join([]string{mc.hTitle, t}, "\n")
		}
	}
	mc.hTitle = strings.TrimLeft(mc.hTitle, "\n")

	mc.hTitlePrepended = false
}

// joinSnippet join sub snippet to current total snippet.
func (mc *markdownContext) joinSnippet(snippet string) {
	if mc.curSnippet == "" {
		mc.curSnippet = snippet
		return
	}

	// check whether current chunk exceeds chunk size, if so, apply to chunks
	if utf8.RuneCountInString(mc.curSnippet)+utf8.RuneCountInString(snippet) >= mc.chunkSize {
		mc.applyToChunks()
		mc.curSnippet = snippet
	} else {
		mc.curSnippet = fmt.Sprintf("%s\n%s", mc.curSnippet, snippet)
	}
}

// applyToChunks applies current snippet to chunks.
func (mc *markdownContext) applyToChunks() {
	defer func() {
		mc.curSnippet = ""
	}()

	var chunks []string
	if mc.curSnippet != "" {
		// check whether current chunk is over ChunkSize，if so, re-split current chunk
		if utf8.RuneCountInString(mc.curSnippet) <= mc.chunkSize+mc.chunkOverlap {
			chunks = []string{mc.curSnippet}
		} else {
			// split current snippet to chunks
			chunks, _ = mc.secondSplitter.SplitText(mc.curSnippet)
		}
	}

	// if there is only H1/H2 and so on, just apply the `Header Title` to chunks
	if len(chunks) == 0 && mc.hTitle != "" && !mc.hTitlePrepended {
		mc.chunks = append(mc.chunks, mc.hTitle)
		mc.hTitlePrepended = true
		return
	}

	for _, chunk := range chunks {
		if chunk == "" {
			continue
		}

		mc.hTitlePrepended = true
		if mc.hTitle != "" && !strings.Contains(mc.curSnippet, mc.hTitle) {
			// prepend `Header Title` to chunk
			chunk = fmt.Sprintf("%s\n%s", mc.hTitle, chunk)
		}
		mc.chunks = append(mc.chunks, chunk)
	}
}

// closeTypes represents the close operation type for each open operation type.
var closeTypes = map[reflect.Type]reflect.Type{ //nolint:gochecknoglobals
	reflect.TypeOf(&markdown.HeadingOpen{}):     reflect.TypeOf(&markdown.HeadingClose{}),
	reflect.TypeOf(&markdown.BulletListOpen{}):  reflect.TypeOf(&markdown.BulletListClose{}),
	reflect.TypeOf(&markdown.OrderedListOpen{}): reflect.TypeOf(&markdown.OrderedListClose{}),
	reflect.TypeOf(&markdown.ParagraphOpen{}):   reflect.TypeOf(&markdown.ParagraphClose{}),
	reflect.TypeOf(&markdown.BlockquoteOpen{}):  reflect.TypeOf(&markdown.BlockquoteClose{}),
	reflect.TypeOf(&markdown.ListItemOpen{}):    reflect.TypeOf(&markdown.ListItemClose{}),
	reflect.TypeOf(&markdown.TableOpen{}):       reflect.TypeOf(&markdown.TableClose{}),
	reflect.TypeOf(&markdown.TheadOpen{}):       reflect.TypeOf(&markdown.TheadClose{}),
	reflect.TypeOf(&markdown.TbodyOpen{}):       reflect.TypeOf(&markdown.TbodyClose{}),
}

// indexOfCloseTag returns the index of the close tag for the open tag at startAt.
func indexOfCloseTag(tokens []markdown.Token, startAt int) int {
	sameCount := 0
	openType := reflect.ValueOf(tokens[startAt]).Type()
	closeType := closeTypes[openType]

	// some tokens (like Hr or Fence) are singular, i.e. they don't have a close type.
	if closeType == nil {
		return startAt
	}

	idx := startAt + 1
	for ; idx < len(tokens); idx++ {
		cur := reflect.ValueOf(tokens[idx]).Type()

		if openType == cur {
			sameCount++
		}

		if closeType == cur {
			if sameCount == 0 {
				break
			}
			sameCount--
		}
	}

	return idx
}

// repeatString repeats the initChar for count times.
func repeatString(count int, initChar string) string {
	var s string
	for i := 0; i < count; i++ {
		s += initChar
	}
	return s
}
