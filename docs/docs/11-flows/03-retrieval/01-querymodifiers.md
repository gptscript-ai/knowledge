# QueryModifiers

## Available Query Modifiers

### spellcheck

Fix typos and grammar by asking the LLM to correct the input query.

**Options**

- `Model`

### enhance

Ask the LLM to enhance the input query in a way that it's optimized for vector similarity search.
  
**Options**
- `Model`

Example: [examples/querymodifier_with_model.yaml](https://github.com/gptscript-ai/knowledge/blob/main/examples/querymodifier_with_model.yaml)

### generic 
  
Provide a custom prompt for the LLM to modify the input query.
It may yield multiple subqueries.

**Options**
- `Model`
- `Prompt`

Example: [examples/querymodifier_with_model.yaml](https://github.com/gptscript-ai/knowledge/blob/main/examples/querymodifier_with_model.yaml)