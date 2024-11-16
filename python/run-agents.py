# Must precede any llm module imports

from langtrace_python_sdk import langtrace

import os

langtrace.init(api_key = os.getenv('LANGTRACE_API_KEY'))

print("Hello from CrewAI and LangTrace!")
