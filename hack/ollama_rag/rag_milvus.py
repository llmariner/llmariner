# inspired from https://python.langchain.com/docs/integrations/vectorstores/milvus/

from langchain_community.document_loaders import TextLoader
from langchain_community.vectorstores import Milvus
from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain_text_splitters import CharacterTextSplitter

from langchain_community.chat_models import ChatOllama
from langchain.schema.output_parser import StrOutputParser
from langchain_community.document_loaders import PyPDFLoader
from langchain.schema.runnable import RunnablePassthrough
from langchain.vectorstores.utils import filter_complex_metadata
from langchain import hub


class ChatPDF:
    vector_store = None
    retriever = None
    chain = None

    def __init__(self):
        self.model = ChatOllama(model="gemma:2b")
        self.text_splitter = CharacterTextSplitter(chunk_size=1000, chunk_overlap=0)
        self.prompt = hub.pull("rlm/rag-prompt-llama")

    def ingest(self, pdf_file_path: str):
        docs = PyPDFLoader(file_path=pdf_file_path).load()
        chunks = self.text_splitter.split_documents(docs)
        chunks = filter_complex_metadata(chunks)

        embeddings = HuggingFaceEmbeddings(model_name="all-MiniLM-L6-v2")
        COLLECTION_NAME = 'rag_db'
        URI = 'http://milvus.default.svc.cluster.local:19530'
        connection_args = { 'uri': URI }
        
        vector_store = Milvus(
            embedding_function=embeddings,
            collection_name=COLLECTION_NAME,
            connection_args=connection_args,
            drop_old=True,
        ).from_documents(
            chunks,
            embedding=embeddings,
            collection_name=COLLECTION_NAME,
            connection_args=connection_args,
        )

        self.retriever = vector_store.as_retriever(
           search_type="mmr", 
           search_kwargs={"k": 1}
        )

        self.chain = ({"context": self.retriever, "question": RunnablePassthrough()}
                      | self.prompt
                      | self.model
                      | StrOutputParser())

    def ask(self, query: str):
        if not self.chain:
            return "Please, add a PDF document first."

        return self.chain.invoke(query)

    def clear(self):
        self.vector_store = None
        self.retriever = None
        self.chain = None

if __name__ == "__main__":
    s = ChatPDF()
    s.ingest("/home/src/SchoolProfile2023-2024.pdf")
    output = s.ask("What classes are offered at Gunn")
    print(output)
