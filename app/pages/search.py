import streamlit as st
import requests
import pandas as pd
from utils import utils

def search():

    st.title("Buscar cliente")
    st.session_state['data'] = {}
    with st.form(key='cadastro_cliente'):
        cpf = st.text_input('CPF do Cliente')
        enter = st.form_submit_button('Buscar')

    if enter:
        resp = utils.search_cpf(cpf)
        
        if resp.status_code == 200 :
            data = resp.json()
            st.session_state['data'] = data
            st.switch_page("pages/customer.py")

        elif resp.status_code == 404:
            st.session_state['data'] = {}
            st.warning("Cliente não encontrado")

    if st.button("Voltar"):
        st.switch_page("pages/main_screen.py")


if __name__ == "__main__":
    search()