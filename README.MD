Запускать в screen

файлы для модели можно скачать ссылкой ниже https://drive.google.com/drive/folders/1EmBRMRyt1Pug6zdIKxC9hLtHOnvtYF0l?usp=sharing
## Запуск локально

```sh
pip install req.txt

uvicorn main:app --reload
```

Nginx раздает фронт и путь /api/ проксирует на uvicorn сервер 
