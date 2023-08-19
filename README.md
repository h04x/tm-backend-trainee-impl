# tm-backend-trainee-impl
Решение [тестового задания](https://github.com/avito-tech/tm-backend-trainee)

## Используемые библиотеки
Роутер Gin  
Postgres pgxpool

## Запуск
Подготовить файл ```.env```  
```docker-compose up --build -d```  
Сделать миграцию вручную из файлика ```migration/up.sql```

## Запросы
Сохранение статистики
```console
curl -v POST http://localhost:12121/save_stat \
-H "Content-Type: application/json" \
-d '{ "date": "2013-12-30", "views": 9, "clicks": 1, "cost": "0.01" }'
```


Получить статистику (Order параметр необязательный)
```console
curl -v --request  POST http://localhost:12121/get_stat \
-H "Content-Type: application/json" \
-d '{ "from": "1234-12-12", "to": "3234-12-12", "Order": "Cost" }'
```
Сброс статистики
```console
curl -v --request DELETE http://localhost:12121/clear_stat
```