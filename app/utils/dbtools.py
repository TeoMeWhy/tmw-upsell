import sqlalchemy
import pandas as pd

def connect_db():
    engine = sqlalchemy.create_engine("sqlite:///data/database.db")
    return engine

def get_products(engine):
    df = pd.read_sql_table("tb_products", engine)
    return df

