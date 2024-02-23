import streamlit as st
from utils import utils
from utils import dbtools


def points_mgmt():

    engine = dbtools.connect_db()
    products = dbtools.get_products(engine)
    list_products = products["descProducts"].unique().tolist()

    data = st.session_state["data"]
    utils.show_user(data)

    with st.form("produtos"):
        columns = st.columns([3, 1])
        columns[0].subheader("Produto")
        columns[1].subheader("Quantidade")

        columns_new = st.columns([3, 1])
        produto = columns_new[0].selectbox(
            label=f"produto", options=list_products, label_visibility="hidden"
        )
        qtde = columns_new[1].slider(f"qtde", 0, 50, 0, label_visibility="hidden")
        enter = st.form_submit_button("Aplicar")

    if enter:

        points = int(
            products[products["descProducts"] == produto]["PtProduct"].values[0] * qtde
        )
        
        if qtde == 0:
            st.warning("Entre com uma quantidade mínima do produto")
        
        else:
            data_product = {produto: qtde}

            resp = utils.add_points(points, data_product)
            if resp.status_code == 200:
                data = utils.search_cpf(data["CPF"]).json()
                st.session_state["data"] = data
                st.switch_page("pages/points_mgmt.py")

            else:
                st.error(f"Erro ao atribuir pontos: {resp.json()}")

    utils.footer_buttons("pages/customer.py")


if __name__ == "__main__":
    points_mgmt()
