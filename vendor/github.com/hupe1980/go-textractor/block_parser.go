package textractor

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// blockParser is responsible for parsing and processing Textract blocks.
type blockParser struct {
	idTypeMap  map[string]types.BlockType
	idBlockMap map[string]types.Block
	typeIDMap  map[types.BlockType][]string
}

// newBlockParser creates a new blockParser instance based on the provided Textract blocks.
func newBlockParser(blocks []types.Block) *blockParser {
	idTypeMap := make(map[string]types.BlockType)
	idBlockMap := make(map[string]types.Block)
	typeIDMap := make(map[types.BlockType][]string)

	for _, b := range blocks {
		id := aws.ToString(b.Id)
		idTypeMap[id] = b.BlockType
		idBlockMap[id] = b

		if strings.HasPrefix(string(b.BlockType), "LAYOUT") {
			typeIDMap[types.BlockType("LAYOUT")] = append(typeIDMap["LAYOUT"], id)
		} else {
			typeIDMap[b.BlockType] = append(typeIDMap[b.BlockType], id)
		}
	}

	return &blockParser{
		idTypeMap:  idTypeMap,
		idBlockMap: idBlockMap,
		typeIDMap:  typeIDMap,
	}
}

// createDocument processes the Textract blocks and creates a structured Document.
func (bp *blockParser) createDocument() *Document {
	ids := bp.blockTypeIDs(types.BlockTypePage)
	pages := make([]*Page, len(ids))

	for i, id := range ids {
		b := bp.blockByID(id)

		page := &Page{
			id:       aws.ToString(b.Id),
			number:   int(aws.ToInt32(b.Page)),
			width:    float64(b.Geometry.BoundingBox.Width),
			height:   float64(b.Geometry.BoundingBox.Height),
			childIDs: filterRelationshipIDsByType(b, types.RelationshipTypeChild),
		}

		pageParser := newPageParser(bp, page)
		pageParser.addPageElements()

		pages[i] = page
	}

	return &Document{
		pages: pages,
	}
}

// blockTypeIDs returns the block IDs of a specific block type.
func (bp *blockParser) blockTypeIDs(blockType types.BlockType) []string {
	return bp.typeIDMap[blockType]
}

// blockByID returns the Textract block with the specified ID.
func (bp *blockParser) blockByID(id string) types.Block {
	return bp.idBlockMap[id]
}

// filterRelationshipIDsByType filters relationship IDs in a block based on the specified relationship type.
func filterRelationshipIDsByType(b types.Block, relationshipType types.RelationshipType) []string {
	var ids []string

	// Iterate through each relationship in the block
	for _, r := range b.Relationships {
		// Check if the relationship type matches the specified type
		if r.Type == relationshipType {
			// Append the IDs associated with the matching type to the result slice
			ids = append(ids, r.Ids...)
		}
	}

	return ids
}
