#!/usr/bin/env python3

import os
import json
import asyncio

async def main():
    output = os.getenv('OUTPUT', '')
    continuation = os.getenv('CONTINUATION', '') == 'true'
    is_chat = os.getenv('CHAT', '') == 'true'

    # only use the part of the output starting with "Retrieved the following"
    if "Retrieved the following" in output:
        output = "Retrieved the following"+output.split("Retrieved the following")[1]

        msg = f"""
Use the content within the following <KNOWLEDGE></KNOWLEDGE> tags as your learned knowledge.
<KNOWLEDGE>
{output}
</KNOWLEDGE>
If this knowledge seems irrelevant to the user query, ignore it.
Avoid mentioning that you retrieved the information from the context or the knowledge tool.
Only provide citations if explicitly asked for it and if the source references are available in the knowledge.
Answer in the language that the user asked the question in.
"""
    else:
        msg = "No data retrieved from knowledge base."

    print(msg)


asyncio.run(main())