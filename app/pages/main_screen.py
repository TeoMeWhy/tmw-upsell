import streamlit as st

def main_screen():
    st.title("Bem-vindo ao Up Sell")
    st.write("Esta é a tela inicial do aplicativo.")

    st.markdown("## Clientes")
    c1, c2 = st.columns(2)
    with c1:
        if st.button("Cadastrar Cliente"):
            st.switch_page("pages/customer_register.py")
    
    with c2:
        if st.button("Buscar Cliente"):
            st.switch_page("pages/search.py")

    st.markdown("---")

    st.markdown("## Produtos")
    if st.button("Cadastrar Produto"):
        pass


if __name__ == "__main__":
    main_screen()