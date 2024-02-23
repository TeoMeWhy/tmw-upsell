import requests
import streamlit as st


def footer_buttons(backpage):
    st.markdown("---")
    c1, c2, c3 = st.columns([1, 3, 1])
    if c1.button("Voltar"):
        st.switch_page(backpage)

    if c3.button("Início"):
        st.switch_page("pages/main_screen.py")


def show_user(data):
    st.markdown(f"## {data['Name'].title()}")

    c1, c2, c3 = st.columns(3)
    c1.markdown(f"**Email**: {data['Email'].lower()}")
    c2.markdown(f"**Tel. Res.**: {data['TelResidencial']}")
    c3.markdown(f"**Tel. Comer.**: {data['TelComercial']}")

    st.markdown(f"**Pontos**: {data['Points']}")


def customer_form(data={}, enter_buttom_name="Cadastrar"):

    data_new = {}

    if data == {}:
        keys = ["CPF","Name","Email","TelResidencial","TelComercial","Instagram"]
        data = {i:"" for i in keys}

    with st.form(key="cadastro_cliente"):
        data_new["CPF"] = st.text_input("CPF do Cliente", value=data["CPF"])
        data_new["Name"] = st.text_input("Nome do Cliente", value=data["Name"])
        data_new["Email"] = st.text_input("Email", value=data["Email"])
        data_new["TelResidencial"] = st.text_input("Telefone Residencial", value=data["TelResidencial"])
        data_new["TelComercial"] = st.text_input("Telefone Comercial", value=data["TelComercial"])
        data_new["Instagram"] = st.text_input("Instagram", value=data["Instagram"])
        enter = st.form_submit_button(enter_buttom_name)

    return data_new, enter


def reward_points(points):
    data = st.session_state["data"]
    url = "http://localhost:8080/addPoints"
    data["Points"] = (-1) * points
    data["Products"] = {"Resgate de Pontos": 1}
    resp = requests.put(url, json=data)
    return resp


def add_points(points, products):
    data = st.session_state["data"]
    url = "http://localhost:8080/addPoints"
    data["Points"] = points
    data["Products"] = products
    resp = requests.put(url, json=data)
    return resp


def search_cpf(cpf):
    url = "http://localhost:8080/customers"
    params = {"cpf": cpf}
    resp = requests.get(url=url, params=params)
    return resp


def excluir(data):
    url = f"http://localhost:8080/customers?id={data['UUID']}"
    resp = requests.delete(url=url)
    return resp

def get_user_transactions(data):
    url = f"http://localhost:8080/transactions?id={data['UUID']}"
    resp = requests.get(url)
    return resp