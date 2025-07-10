import os
import sys

from openai import OpenAI

if len(sys.argv) != 2:
    print("Usage: python run_transcription.py <audio_file>")
    sys.exit(1)

audio_file = sys.argv[1]

client = OpenAI(
    api_key=os.getenv("LLMARINER_API_KEY"),
    base_url="http://localhost:8080/v1",
)

model = "openai-whisper-large-v3-turbo"

response = client.audio.transcriptions.create(
    model=model,
    file=open(audio_file, "rb")
)
print(response)
