# Desafio Clean Architecture

## Descrição Desafio

Criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

## Pré-requisitos (testado no linux)

- Make versão 4.3
- Go versão 1.22.2
- Docker versão 24.0.7
- Docker Compose versão v2.3.3

## Execução

1. Iniciar Api:
```bash
make run
```

## Endpoints:
- Rest localhost:8000/order
  - [create_order](./api/create_order.http)
  - [Listar ordens](./api/list_order.http)
- gRPC server localhost:50051
  - Utilizar o [Evans](https://github.com/ktr0731/evans) para as requisições
- GraphQL server localhost:8080
  - Criar ordem  
    ```graphql
    mutation CreateOrder {
        createOrder(input: {id: "someIDD", Price: 44, Tax: 9}) {
            id
            Price
            Tax
            FinalPrice
        }
    }
    ```
  - Listar ordens  
    ```graphql
    query List {
        listOrders {
            id
            Price
            Tax
            FinalPrice
        }
    }
    ```
