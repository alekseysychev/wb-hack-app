import pandas as pd
import numpy as np
from gensim.models.fasttext import FastText 
import gensim
from numba import jit
from fastapi import FastAPI

app = FastAPI()


qp = pd.read_csv('./qp.csv', index_col = 'Unnamed: 0', engine='python')


model = FastText.load('./fasttext.model')

X_wv = np.load('./X_wv.npy')


def make_trans():
    a = 'a b c d e f g h i j k l m n o p q r s t u v w x y z ё'.split()
    b = 'а б с д е ф г х и ж к л м н о п к р с т у в в х у з е'.split()
    trans_dict = dict(zip(a, b))
    trans_table = ''.join(a).maketrans(trans_dict)
    return trans_table
    
trans_table = make_trans()

@jit(nopython=True)
def cosine_similarity_numba(u, v):
    cos = []
    for j in range(u.shape[0]):
      assert(u[j].shape[0] == v[j].shape[0])
      uv = 0
      uu = 0
      vv = 0
      for i in range(u[j].shape[0]):
          uv += u[j][i]*v[j][i]
          uu += u[j][i]*u[j][i]
          vv += v[j][i]*v[j][i]
      cos_theta = 1
      if uu!=0 and vv!=0:
          cos_theta = uv/np.sqrt(uu*vv)
      cos.append(cos_theta)
    return cos



#топ рекомендаций без нажатия клавиши и только с 1 символом
def q_rec_top(n, symbol=''):
  q_pop = qp.sort_values('query_popularity', ascending=False)['query']
  if n == 0:
    top10_without_sym = q_pop[:10].to_list()
    return top10_without_sym
  if n == 1:
    top10_with_1_sym = q_pop[q_pop.str.find(symbol) == 0][:10].tolist()
    return top10_with_1_sym

## предикт функция
@app.get("/api/query/{input}")
def sBar(input:str):
  text = input
  text = text.lower()
  text = text.translate(trans_table)
  if len(text) <=1:
    top_tags = q_rec_top(len(text), text)
  else:
    pred_vect = model.wv[text]
    pred_t = np.expand_dims(pred_vect, axis=0)
    k = np.repeat(pred_t, X_wv.shape[0], axis=0)
    k = k.astype('float64')
    score = cosine_similarity_numba(X_wv,k)
    ind_pred_10 = np.argsort(score)[-10:][::-1]
    ind_pred_100 = np.argsort(score)[-100:][::-1]
    top_10 = qp.iloc[ind_pred_10].sort_values('query_popularity', ascending=False)
    top_100 = qp.iloc[ind_pred_100].sort_values('query_popularity', ascending=False)
    ten_1=top_100[top_100['query_new'].str.contains(fr"^{text}")]['query_new'][:10]
    ten_2 = pd.Series([0])
    top_ten = ten_1.to_list()
    if len(ten_1) < 10:
      ten_2 = top_100[~top_100['query_new'].str.contains(fr"^{text}") & top_100['query_new'].str.contains(fr"{text}")]['query_new'][:(10 -len(ten_1))]
      top_ten = ten_1.to_list() + ten_2.to_list()
  return top_ten if top_ten else top_100['query_new'][:10].to_list()


@app.get("/api/test/{input}")
def test(input:str):
  print(input)
  return 
