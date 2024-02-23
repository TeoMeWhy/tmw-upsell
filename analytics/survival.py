# %%

import pandas as pd
import sqlalchemy

import matplotlib.pyplot as plt

def import_query(path):
    with open(path, "r") as open_file:
        query = open_file.read()
    return query

query = import_query("survival.sql")

# %%
engine = sqlalchemy.create_engine("sqlite:///../data/database.db")

# %%
df = pd.read_sql_query(query, engine)
df
# %%

plt.plot(df['qtDays'], df['pct_acum'])
plt.grid(True)
plt.title("Probabilidade Acumulada - Recorrência Teo Me Why")
plt.xlabel("Dias para recorrência")
plt.ylabel("% Base")
# %%
