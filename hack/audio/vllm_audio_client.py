# Ref: https://docs.vllm.ai/en/v0.6.0/getting_started/examples/openai_audio_api_client.html

import base64

from openai import OpenAI

# Modify OpenAI's API key and API base to use vLLM's API server.
openai_api_key = "openai"
openai_api_base = "http://localhost:8080/v1"

client = OpenAI(
    # defaults to os.environ.get("OPENAI_API_KEY")
    api_key=openai_api_key,
    base_url=openai_api_base,
)

model = "fixie-ai-ultravox-v0_3"

# Any format supported by librosa is supported
audio_url = "https://vllm-public-assets.s3.us-west-2.amazonaws.com/multimodal_asset/winning_call.ogg"

# Use audio url in the payload
chat_completion_from_url = client.chat.completions.create(
    messages=[{
        "role": "user",
        "content": [
            {
                "type": "text",
                "text": "summarize the audio?"
            },
            {
                "type": "audio_url",
                "audio_url": {
                    "url": audio_url
                },
            },
        ],
    }],
    model=model,
    max_tokens=1024,
)

result = chat_completion_from_url.choices[0].message.content
print(f"Chat completion from url output:{result}")

print("\n\n")

file_path = "./Sports.wav"
# Use base64 encoded audio in the payload
def encode_audio_base64_from_file(file_path: str) -> str:
    """Encode an audio retrieved from a remote url to base64 format."""

    with open(file_path, 'rb') as f:
        result = base64.b64encode(f.read()).decode('utf-8')

    return result

audio_base64 = encode_audio_base64_from_file(file_path=file_path)
chat_completion_from_base64 = client.chat.completions.create(
    messages=[{
        "role":
        "user",
        "content": [
            {
                "type": "text",
                "text": "summarize the audio"
            },
            {
                "type": "audio_url",
                "audio_url": {
                    "url": f"data:audio/wav;base64,{audio_base64}"
                },
            },
        ],
    }],
    model=model,
    max_tokens=1024,
)

result = chat_completion_from_base64.choices[0].message.content
print(f"Chat completion from file output:{result}")

