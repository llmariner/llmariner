import streamlit as st
from openai import OpenAI

st.markdown(
    """
# Chat
"""
)

client = OpenAI(
  base_url="http://localhost:80/v1",
  api_key='dummy',
)

if "messages" not in st.session_state:
    st.session_state.messages = []

for message in st.session_state.messages:
    with st.chat_message(message["role"]):
        st.markdown(message["content"])

if prompt := st.chat_input("What is up?"):
    st.session_state.messages.append({"role": "user", "content": prompt})
    with st.chat_message("user"):
        st.markdown(prompt)

    with st.chat_message("assistant"):
        stream = client.chat.completions.create(
            model="gemma:2b",
            messages=[{"role": "user", "content": prompt}],
            stream=True,
        )
        response = st.write_stream(stream)
    st.session_state.messages.append({"role": "assistant", "content": response})
