# Rate Limiter API

## Propósito

Esta API implementa um rate limiter (limitador de requisições) utilizando Redis como mecanismo de persistência. O objetivo é proteger sistemas de abusos, limitando o número de requisições por IP ou por token de acesso (API Key) em um determinado intervalo de tempo. O limite por token, se existir, sempre se sobrepõe ao limite por IP.

## Funcionalidades
- Limite de requisições por IP e/ou por API Key.
- Janela de tempo configurável (ex: 1 minuto).
- Tempo de bloqueio configurável após exceder o limite.
- Persistência eficiente usando Redis.
- Middleware pronto para uso em aplicações Go (chi).

## Estrutura do Projeto
```
├── cmd/server/           # Código principal da API
├── configs/              # Configuração e leitura de variáveis de ambiente
├── internal/entity/      # Entidades e estratégias de persistência
│   ├── client_entity/    # Interface e implementação do client Redis
│   └── rate_limiter/     # Lógica do rate limiter e middleware
├── Dockerfile            # Build da aplicação Go
├── docker-compose.yml    # Sobe o Redis via Docker
├── .env                  # Variáveis de ambiente
```

## Pré-requisitos
- Go 1.21+
- Docker e Docker Compose

## Configuração

1. **Clone o repositório:**
   ```sh
   git clone https://github.com/JoaoPedroVicentin/rate-limiter
   cd rate-limiter
   ```

2. **Configure o arquivo `.env` na raiz do projeto:**
   Exemplo:
   ```env
   API_KEY=12345
   WEB_SERVER_PORT=8080
   REDIS_ADDRESS=localhost:6379
   MAX_REQUESTS_PER_IP=5
   MAX_REQUESTS_PER_API_KEY=10
   TIME_INTERVAL_MINUTES=1
   BLOCK_DURATION_MINUTES=1
   ```

3. **Suba o Redis com Docker Compose:**
   ```sh
   docker compose up -d
   ```

4. **Rode a API localmente:**
   ```sh
   cd cmd/server
   go run main.go
   ```

## Testando a API

1. **Faça requisições HTTP:**
   ```sh
   curl -H "api-key: 12345" http://localhost:8080/
   ```
   - Troque o valor do header `api-key` conforme sua configuração.
   - Se exceder o limite, receberá status 429 (Too Many Requests).

2. **Testes automatizados:**
   ```sh
   go test github.com/JoaoPedroVicentin/rate-limiter/internal/entity/rate_limiter
   ```

<div align="center">
<h3>👨‍💻</h3>
    <h3> Criado por João Pedro Vicentin!</h3>
    <div>
        <h3>
            <a href="https://www.linkedin.com/in/joaopedrovicentin/" target="_blank">Linkedin</a>
            <a href='https://github.com/JoaoPedroVicentin' target='_blank'>Github</a>
            <a href="https://contate.me/joao-pedro-lopes-vicentin" target="_blank">Whatsapp</a>
        </h3>
    </div>
</div>