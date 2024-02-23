import streamlit as st
import requests

from utils import utils

def customer():

    data = st.session_state["data"]

    placeholder = st.empty()
    with placeholder.container():
        utils.show_user(data)

        c1, c2,_, c3, c4 = st.columns([1,1,2,1,1])

        adicao = c1.button("Pontos")
        extrato = c2.button("Extrato")
        editar = c3.button("Editar cliente")
        exclusao = c4.button("Excluir cliente")

    if adicao:
        st.switch_page("pages/points_mgmt.py")

    if extrato:
        st.switch_page("pages/points_report.py")

    if editar:
        st.switch_page("pages/customer_edit.py")

    if exclusao:
        utils.excluir(data)
        placeholder.empty()
        st.info("Usuário removido com sucesso!")

    utils.footer_buttons("pages/search.py")




if __name__ == "__main__":
    customer()
