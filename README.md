# API-Blockchain
API de conexão com blockchain

1 - Autenticação
    a - Cadastro
    b - Login
    c - Recuperação de senha
    d - JWT e api-key para comunicação entre sistemas
    e - Notificação via email (Sendgrid)
    f - níveis de permissão

2 - Inserção de dados brutos
    a - Upload de arquivos (10mb)
    b - Strucut de informações desses arquivo (e.g nome, criador, data, comentários)
    c - interação com o contrato
    d - Notificação via email (Sendgrid) - para admin e para publisher

3 - Inserção de teses
    a - Upload de arquivos (10mb)
    b - Strucut de informações desses arquivo (e.g nome, criador, data, responsável)
    c - interação com o contrato
    d - Notificação via email (Sendgrid) - para admin e para publisher

4 - Consulta de dados
    a - Validação de tempo de permissionamento no DB
    b - retorno em formato json e struct de dados
    d - Notificação via email (Sendgrid) - para admin e para publisher

5 - Permissionamento
    a - Só adms podem conceder permissionamento
    b - garantir permissões (WRITER, READER, UPDATER, ADMIN)
    c - garantir permissões temporários (Calculo do bloco)
    d - Revogar permissão temporária
    e - extender permissão temporária


Integrações
    1 - Sengrid - cliente email
    2 - MongoDB/Compass - banco de dados
    3 - Integraçaõ com blockchain (https://pkg.go.dev/github.com/ethereum/go-ethereum)
    4 - Token JWT
    5 - Segurança (Helmet e CORS)
    6 - Pinata ou Web3Storage (https://web3.storage/ https://www.pinata.cloud/)

Arquitetura

    1 - Docker
    2 - Arquitetura limpa 
    3 - Padrão de código (BDD, DDD, MVC , etc)
    4 - Testes
    5 - CI/CD (Github Actions yaml) (https://www.redhat.com/pt-br/topics/devops/what-is-ci-cd)



------- How to run --------

1 - Create a .env at root :
API_KEY = "yourApiKey"
API_SECRET "yourApiSecret"
JWT_TOKEN "yourJwtToken"

2 - ro gun cmd/api/main/main.go

3 - Send a request to http://localhost:8080/register to register

Attach in the Body of the request
{"username":"your_username", "password":"your_password", "email": "your_email"}

4 - Send a request to http://localhost:8080/login to login, only with username and password at the body:
{"username":"your_username", "password":"your_password"}

5 - To upload a file to ipfs http://localhost:8080/upload and add a file in the body of the request

6 - To interact with a contract, Add a "Action" header on the request. Also,
the JWT token shoulb be passed as "Authorization" header 

*Remember to put your contract and account address on BlockchainInteracion function
Command to generate Go Bindings
abigen --abi=./cmd/api/contract/contract.abi --bin=./cmd/api/contract/contract.bin --pkg=contract --out=./cmd/api/contract/contract.go