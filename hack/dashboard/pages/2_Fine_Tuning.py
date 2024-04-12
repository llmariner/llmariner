import streamlit as st
from openai import OpenAI
import pandas as pd
import datetime

st.markdown(
    """
# Fine Tuning Dashboard
"""
)

# style
th_props = [
  ('font-size', '18px'),
  ('text-align', 'center'),
  ('font-weight', 'bold'),
  ('color', '#6d6d6d'),
]

td_props = [
  ('font-size', '18px')
]

styles = [
  dict(selector="th", props=th_props),
  dict(selector="td", props=td_props)
]

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
st.table(pd.DataFrame(modelsData).style.set_properties(**{'text-align': 'left'}).set_table_styles(styles))

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
st.table(pd.DataFrame(filesData).style.set_properties(**{'text-align': 'left'}).set_table_styles(styles))

jobs = client.fine_tuning.jobs.list()
jobsData = {
    'ID': [],
    'Created at': []
}
for j in jobs.data:
    jobsData['ID'].append(j.id)
    jobsData['Created at'].append(datetime.datetime.fromtimestamp(int(j.createdAt)))
st.subheader("Jobs")
st.table(pd.DataFrame(jobsData).style.set_properties(**{'text-align': 'left'}).set_table_styles(styles))
