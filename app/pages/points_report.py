import streamlit as st
from utils import utils
import pandas as pd

def points_report():

    data = st.session_state['data']

    utils.show_user(data)

    resp = utils.get_user_transactions(data)
    if resp.status_code != 200:
        return
    
    df = pd.DataFrame(resp.json()["transactions"])

    df = df.rename(columns={
                    "cpf":"CPF",
                    "Name":"Nome",
                    "dtTransaction": "Data",
                    "Points":"Pontos",
                    "product": "Produto",
                    "qtdeProduto": "Quantidade"})
    
    df["Nome"] = df["Nome"].apply(lambda x: str(x).title())
    st.dataframe(df)

    utils.footer_buttons("pages/customer.py")

if __name__ == "__main__":
    points_report()