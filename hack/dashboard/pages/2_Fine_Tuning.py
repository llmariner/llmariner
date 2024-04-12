import streamlit as st
from openai import OpenAI
import pandas as pd
import datetime

st.markdown(
    """
# Fine Tuning
"""
)

client = OpenAI(
  base_url="http://localhost:80/v1",
  api_key='dummy',
)

models = client.models.list()
modelsData = {
    'ID': [],
    'Created at': []
}
for m in models.data:
    modelsData['ID'].append(m.id)
    modelsData['Created at'].append(datetime.datetime.fromtimestamp(int(m.created)))
st.subheader("Models")
st.table(pd.DataFrame(modelsData))

files = client.files.list()
filesData = {
    'ID': [],
    'Filename': [],
    'Created at': []
}
for f in files.data:
    filesData['ID'].append(f.id)
    filesData['Filename'].append(f.filename)
    filesData['Created at'].append(datetime.datetime.fromtimestamp(int(f.createdAt)))
st.subheader("Files")
st.write(pd.DataFrame(filesData))

jobs = client.fine_tuning.jobs.list()
jobsData = {
    'ID': [],
    'Created at': []
}
for j in jobs.data:
    jobsData['ID'].append(j.id)
    jobsData['Created at'].append(datetime.datetime.fromtimestamp(int(j.createdAt)))
st.subheader("Jobs")
st.write(pd.DataFrame(jobsData))
