import streamlit as st
import requests
from utils import utils

def customer_register():

    placeholder = st.empty()

    st.title("Cadastro de cliente")
    st.write("Cadastre aqui seus clientes")

    placeholder = st.empty()

    with placeholder.container():
        data, enter = utils.customer_form()

    if enter:
        url = "http://localhost:8080/customers"
        resp = requests.post(url=url, json=data)

        if resp.status_code == 201:
            placeholder.empty()
            st.info("Cliente cadastrado com sucesso!")
        else:
            body = resp.json()
            if body["error"] == "Usuário já existente":
                st.error("Usuário já existente! Insira um novo CPF!")
            else:
                st.error("Não foi possível cadastrar o cliente.")
                st.info(f"Erro interno. {body}")

    if st.button("Voltar"):
        st.switch_page("pages/main_screen.py")


if __name__ == "__main__":
    customer_register()
