import requests
import streamlit as st

from utils import dbtools

from pages.main_screen import main_screen

st.session_state.data = {}

def main():
    main_screen()

if __name__ == "__main__":
    main()
