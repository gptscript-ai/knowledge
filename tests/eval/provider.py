# Learn more about building a Python provider: https://promptfoo.dev/docs/providers/python/
import json
import subprocess

def call_api(prompt, options, context):
    # The 'options' parameter contains additional configuration for the API call.
    config = options.get('config', None)
    # additional_option = config.get('additionalOption', None)

    # The 'context' parameter provides info about which vars were used to create the final prompt.
    file = context['vars'].get('file', None)

    # The prompt is the final prompt string after the variables have been processed.
    # Custom logic to process the prompt goes here.
    # For instance, you might call an external API or run some computations.

    print(f">>> Checking on file: {file}")


    # exec a binary
    input = f'{{"prompt": "{prompt}"}}'
    result = subprocess.run(["gptscript", "--workspace=assets/", "--disable-cache", "provider.gpt", input], capture_output=True, text=True)

    if result.returncode != 0:
        print(f"Error: {result.stderr}")
        return {"error": result.stderr}


    # The result should be a dictionary with at least an 'output' field.
    result = {
        "output": result.stdout,
    }



    return result

