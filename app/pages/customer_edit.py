import requests
import streamlit as st
from utils import utils

def customer_edit():
    data = st.session_state["data"]
    
    placeholder = st.empty()
    
    with placeholder.container():
        data_new, enter = utils.customer_form(data, "Aplicar")

    if enter:
        data.update(data_new)
        url = "http://localhost:8080/customers"
        resp = requests.put(url=url, json=data)

        if resp.status_code == 200:
            st.info("Cliente alterado com sucesso!")
            
            data.update(data_new)
            st.session_state["data"] = data

        else:
            st.error("Erro ao alterar o cliente: ", resp.json()["error"])

    utils.footer_buttons("pages/customer.py")

if __name__ == "__main__":
    customer_edit()