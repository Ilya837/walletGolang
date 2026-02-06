# Суть проекта:
Приложение представляет из себя сервер, написанный на golang, который принимает REST запросы и взаимодействует с базой данных. И сервер и база данных запускается в docker контейнере при помощи docker compose. Проект в общем представляет из себя кошелёк, который можно или пополнить или списать с него сумму. 

# Перед запуском:

перед запуском в основной папке проекта (рядом с main.go) нужно создать файл config.env. В нём нужно прописать следующие переменные:

- POSTGRES_DB
- POSTGRES_USER 
- POSTGRES_PASSWORD
- POSTGRES_TABLE
- SERVER_PORT

Пример:

```
POSTGRES_DB=wallets \
POSTGRES_USER=user \
POSTGRES_PASSWORD=1234 \
POSTGRES_TABLE=wallets \
SERVER_PORT=:80 
```

# Запуск:

## При помощи make:
```
<<<<<<< HEAD
make docker-compose-run
make docker-compose-build
=======
make docker-compose-build
make docker-compose-run
>>>>>>> 7633285c290da79ef7996dc43c4a8e1cfd852d30
```
## Без make:
```
docker compose --env-file config.env -p walletapp build
docker compose --env-file config.env -p walletapp up -d
```

# Остановка:

## При помощи make:
```
make docker-compose-stop
```
## Без make:
```
docker compose --env-file config.env -p walletapp down
```
# Принимаемые запросы:

- GET api/v1/wallets/{WALLET_UUID}

        выдаёт балланс на кошельке с соответствующим id

- POST api/v1/wallet 
{
walletId: UUID,
operationType: DEPOSIT or WITHDRAW,
amount: 1000
}

        увеличивает/уменьшает баланс кошелька

- POST api/v1/wallets/wallet/create
{
walletId: UUID
}

        Создаёт кошелёк с соответствующим id (если такого ещё нет)

