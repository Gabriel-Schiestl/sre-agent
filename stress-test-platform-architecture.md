# Stress Test Intelligence Platform — Documento de Arquitetura e Planejamento

> Versão 1.0 — PoC  
> Status: Em planejamento

---

## 1. Visão Geral

### 1.1 O Problema

O processo atual de stress testing em ambientes com múltiplos microsserviços é lento, manual e altamente dependente da experiência do desenvolvedor. O fluxo típico envolve:

1. Executar o teste no **JMeter** e aguardar a conclusão
2. Analisar manualmente os erros gerados no relatório `.jtl`
3. Abrir o **Grafana** para acompanhar consumo de CPU, memória e latência
4. Consultar o **Kubernetes** para verificar restarts de pods e eventos
5. Inferir qual microsserviço está com gargalo com base na experiência pessoal
6. Adicionar logs à aplicação, fazer build, deploy e reiniciar o teste
7. Repetir o ciclo até o problema ser identificado e corrigido

Esse loop é **moroso, repetitivo e pouco sistemático**. O diagnóstico é subjetivo, o conhecimento gerado se perde e cada desenvolvedor reinventa a roda a cada novo teste.

### 1.2 A Solução

Uma plataforma interna de inteligência para stress testing que **automatiza a coleta de evidências, correlaciona dados de múltiplas fontes e utiliza IA para gerar diagnósticos estruturados e acionáveis**.

O desenvolvedor passa de executor manual para supervisor — ele configura o contexto uma vez, faz o upload do resultado do teste e recebe um plano de diagnóstico com hipóteses priorizadas, evidências e próximos passos sugeridos.

### 1.3 Proposta de Valor

| Antes | Depois |
|---|---|
| Análise manual de centenas de erros | Diagnóstico estruturado gerado automaticamente |
| Conhecimento informal e subjetivo | Hipóteses priorizadas com evidências |
| Loop de build/deploy para adicionar logs | Contexto de logs e métricas já agregados |
| Sem histórico entre execuções | Histórico de runs por suite consultável |
| Depende da experiência individual | IA com contexto completo dos microsserviços |

---

## 2. O Sistema

### 2.1 Conceitos Fundamentais

**Test Suite**
Agrupador permanente que representa um cenário de stress testing recorrente. Uma suite contém os microsserviços envolvidos naquele cenário e acumula todas as execuções históricas. É configurada uma vez e reutilizada em todos os testes subsequentes.

**Microservice**
Unidade cadastrada dentro de uma suite. Representa um microsserviço real do ambiente, com seus metadados relevantes: linguagem, endpoints principais, limites de recursos e SLOs definidos. Esse cadastro é o contexto que a IA utiliza para gerar diagnósticos precisos.

**Test Run**
Uma execução específica de uma suite. Criado a cada upload de arquivo `.jtl`, carrega metadados do teste (usuários virtuais, duração, notas) e passa pelo pipeline de análise. Possui um ciclo de vida próprio com status rastreável.

**Diagnosis**
O produto final do sistema. Um diagnóstico estruturado gerado pela IA com base nos dados agregados do run e no contexto da suite. Contém categorização de erros, identificação de microsserviços suspeitos, hipóteses priorizadas e próximos passos sugeridos.

### 2.2 Fluxo Principal

```
┌─────────────────────────────────────────────────────────────┐
│  1. Dev cadastra uma Test Suite                             │
│     └── Adiciona os Microservices envolvidos               │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│  2. Dev entra na suite e cria um Test Run                   │
│     └── Informa: usuários virtuais, duração, notas         │
│     └── Faz upload do arquivo .jtl                         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│  3. Sistema parseia e agrega o .jtl (runner)               │
│     └── Extrai métricas por endpoint                       │
│     └── Calcula p50/p90/p99, error rate, timeline          │
│     └── Classifica erros por tipo e frequência             │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│  4. IA analisa os dados com contexto completo (analyst)    │
│     └── Recebe: métricas + dados da suite + microservices  │
│     └── Gera: categorização, suspeitos, hipóteses          │
│     └── Prioriza por impacto e evidência                   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│  5. Dev visualiza o diagnóstico e interage via chat        │
│     └── Relatório estruturado com evidências               │
│     └── Chat contextual sobre o run específico             │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 Processamento Assíncrono

O processamento após o upload é **assíncrono**. O motivo é objetivo: a combinação de parse de arquivo potencialmente grande, agregação de métricas e chamada à IA pode levar entre 15 e 60 segundos. Processar de forma síncrona degradaria a experiência e criaria risco de timeout.

O fluxo técnico:

1. Upload recebido → `TestRun` criado com `status: pending`
2. Goroutine disparada em background
3. Frontend faz polling a cada 3 segundos no endpoint de status
4. Quando `status: done`, frontend carrega o diagnóstico automaticamente

Essa abordagem não requer infraestrutura de fila no MVP. Na evolução do sistema, a goroutine é substituída por um worker com fila (NATS, RabbitMQ) sem alteração na API.

---

## 3. Arquitetura

### 3.1 Visão Macro

```
┌─────────────────────────────────────────────┐
│              Frontend (Next.js)             │
│  Suites │ Microservices │ Runs │ Diagnóstico │
└──────────────────┬──────────────────────────┘
                   │ HTTP / REST
┌──────────────────▼──────────────────────────┐
│              Backend (Go)                   │
│         Monolito Modular                    │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │            registry                  │   │
│  │  Ponto de entrada e orquestrador     │   │
│  │  TestSuite │ Microservice │ TestRun  │   │
│  │  Todos os HTTP handlers              │   │
│  └────────┬─────────────────┬───────────┘   │
│           │                 │               │
│  ┌────────▼──────┐ ┌────────▼──────────┐    │
│  │    runner     │ │      analyst      │    │
│  │  Motor de     │ │  Motor de IA      │    │
│  │  execução     │ │  Prompt + Claude  │    │
│  │  Parse JTL    │ │  Diagnóstico      │    │
│  │  Agregação    │ │  Chat             │    │
│  └───────────────┘ └───────────────────┘    │
│                                             │
│  /pkg/llm        → Client Anthropic         │
│  /cmd/api        → Main, injeção, rotas     │
└─────────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│              Banco de Dados                 │
│         SQLite (PoC) → PostgreSQL           │
└─────────────────────────────────────────────┘
```

### 3.2 Princípio dos Módulos

O backend adota **monolito modular** — um único binário Go com fronteiras internas bem definidas entre módulos. Cada módulo tem responsabilidade única e comunica-se com os outros exclusivamente por **interfaces e contratos explícitos de dados**, nunca por import direto de tipos internos.

A regra de ouro:

> `registry` é o único orquestrador. `runner` e `analyst` são motores puros — recebem dados ricos, processam e devolvem resultados. Não acessam o banco. Não conhecem HTTP.

### 3.3 Módulos

#### `registry` — Orquestrador
- Expõe todos os HTTP handlers do sistema
- Gerencia o CRUD de `TestSuite`, `Microservice` e `TestRun`
- Orquestra o fluxo de análise: chama `runner`, recebe resultado, chama `analyst`, persiste diagnóstico
- Único módulo com acesso ao banco de dados
- Atualiza o status do `TestRun` ao longo do ciclo de vida

#### `runner` — Motor de Execução
- Recebe `RunPayload` completo do registry (sem lookups)
- Parseia o arquivo `.jtl` em formato CSV
- Agrega métricas por endpoint: p50, p90, p99, error rate, volume
- Classifica erros por tipo (`responseCode`, `failureMessage`)
- Constrói timeline de degradação por janelas de tempo
- Devolve `AggregatedData` estruturado ao registry
- Sem handlers HTTP. Sem acesso ao banco.

#### `analyst` — Motor de IA
- Recebe `AnalysisPayload` completo do registry
- Constrói o prompt com contexto dos microsserviços e métricas agregadas
- Chama a API do Claude via `/pkg/llm`
- Parseia e valida a resposta estruturada em JSON
- Gerencia o histórico de chat contextual por run
- Devolve `Diagnosis` estruturado ao registry
- Sem handlers HTTP. Sem acesso ao banco.

#### `/pkg/llm` — Client Anthropic
- Wrapper sobre a API da Anthropic
- Sem lógica de negócio
- Reutilizável por qualquer módulo que precise chamar a IA

#### `/cmd/api` — Entrypoint
- Instancia todos os módulos
- Injeta dependências (banco, llm client)
- Registra as rotas de cada módulo no router
- Configura middleware (CORS, logging, recovery)
- Nenhuma lógica de negócio

### 3.4 Estrutura de Pastas

```
/
├── cmd/
│   └── api/
│       └── main.go              ← entrypoint, injeção, rotas
├── internal/
│   ├── registry/
│   │   ├── handler.go           ← HTTP handlers
│   │   ├── service.go           ← orquestração e regras de negócio
│   │   ├── repository.go        ← acesso ao banco
│   │   └── model.go             ← TestSuite, Microservice, TestRun, Diagnosis
│   ├── runner/
│   │   ├── service.go           ← orquestra parse e agregação
│   │   ├── parser.go            ← leitura do CSV JTL
│   │   ├── aggregator.go        ← cálculo de métricas
│   │   └── model.go             ← RunPayload, AggregatedData
│   └── analyst/
│       ├── service.go           ← orquestra prompt, LLM e chat
│       ├── prompt.go            ← construção do prompt
│       └── model.go             ← AnalysisPayload, Diagnosis
└── pkg/
    └── llm/
        └── client.go            ← wrapper Anthropic SDK
```

### 3.5 Contratos entre Módulos

```go
// registry → runner
type RunPayload struct {
    Run           TestRun
    Suite         TestSuite
    Microservices []Microservice
    JTLContent    []byte
}

// runner → registry
type AggregatedData struct {
    TotalRequests   int
    ErrorRate       float64
    LatencyP50Ms    float64
    LatencyP90Ms    float64
    LatencyP99Ms    float64
    ErrorsByType    []ErrorGroup
    EndpointMetrics []EndpointMetric
    Timeline        []TimelinePoint
}

// registry → analyst
type AnalysisPayload struct {
    Run           TestRun
    Suite         TestSuite
    Microservices []Microservice
    Data          AggregatedData
}

// analyst → registry
type Diagnosis struct {
    ErrorPlan   []ErrorCategory
    Bottlenecks []Bottleneck
    NextSteps   []string
    RawResponse string
}
```

---

## 4. Modelo de Dados

### 4.1 Entidades

```
TestSuite
├── id            UUID
├── name          string
├── description   string
├── createdAt     timestamp
└── updatedAt     timestamp

Microservice
├── id              UUID
├── testSuiteId     UUID (FK)
├── name            string
├── description     string
├── language        string         ("Go", "Java", "Node.js"...)
├── mainEndpoints   []string
├── cpuLimit        string         ("500m") — opcional
├── memoryLimit     string         ("512Mi") — opcional
├── sloLatencyP99Ms int            — opcional
├── sloErrorRatePct float64        — opcional
└── createdAt       timestamp

TestRun
├── id              UUID
├── testSuiteId     UUID (FK)
├── name            string         ("Black Friday Simulation v3")
├── virtualUsers    int
├── durationSeconds int
├── notes           string
├── status          enum           (pending|analyzing|done|failed)
├── jtlFilePath     string
└── createdAt       timestamp

Diagnosis
├── id           UUID
├── testRunId    UUID (FK, 1:1)
├── errorPlan    JSON
├── bottlenecks  JSON
├── nextSteps    JSON
├── chatHistory  JSON
├── rawResponse  string
└── createdAt    timestamp
```

### 4.2 Campos do JTL Utilizados

O arquivo `.jtl` em formato CSV contém muitos campos. O sistema extrai apenas os relevantes para análise:

| Campo | Uso |
|---|---|
| `timeStamp` | Correlação temporal, construção da timeline |
| `elapsed` | Latência em ms — base para p50/p90/p99 |
| `label` | Nome do endpoint/request |
| `responseCode` | HTTP status code — classificação de erros |
| `success` | Filtro primário de falhas |
| `failureMessage` | Mensagem de erro detalhada |
| `allThreads` | Usuários ativos no momento — correlação de carga |
| `URL` | Endpoint exato chamado |

---

## 5. Inteligência Artificial

### 5.1 Estratégia de Prompt

O prompt é a peça mais crítica do sistema. A qualidade do diagnóstico depende diretamente da qualidade do contexto fornecido à IA. A estrutura adotada:

```
[SISTEMA]
Você é um especialista em performance e SRE. Analisa resultados
de stress tests e identifica gargalos com base em evidências concretas.
Responda sempre em JSON estruturado conforme o schema fornecido.
Seja específico — hipóteses genéricas não têm valor.

[CONTEXTO DOS MICROSSERVIÇOS]
Para cada microsserviço:
- Nome, linguagem, descrição
- Endpoints principais
- SLOs definidos (latência p99, error rate aceitável)
- Limites de recursos (CPU, memória)

[PERFIL DO TESTE]
- Usuários virtuais
- Duração total
- Notas do desenvolvedor

[DADOS AGREGADOS DO JMETER]
- Total de requests e error rate geral
- Métricas por endpoint: p50/p90/p99, volume, error rate
- Top erros por frequência: responseCode + mensagem
- Timeline de degradação: janelas de 30s com spike de erro

[TAREFA]
1. Categorize cada tipo de erro (timeout, 5xx, connection refused, etc)
2. Identifique os microsserviços com maior probabilidade de ser gargalo
3. Para cada suspeito, liste hipóteses priorizadas com evidências
4. Aponte os endpoints mais críticos com dados que sustentam
5. Sugira próximos passos concretos de investigação

[SCHEMA DE RESPOSTA]
{ schema JSON aqui }
```

### 5.2 Schema de Resposta da IA

```json
{
  "errorPlan": [
    {
      "category": "connection_timeout",
      "description": "Requisições excedendo tempo limite de conexão",
      "occurrences": 847,
      "affectedEndpoints": ["POST /payments", "GET /orders/{id}"],
      "severity": "critical"
    }
  ],
  "bottlenecks": [
    {
      "microservice": "payment-service",
      "confidence": "high",
      "hypotheses": [
        {
          "title": "Esgotamento de pool de conexões com banco",
          "evidence": "p99 de 4200ms em POST /payments com error rate de 23% a partir de 150 usuários simultâneos",
          "priority": 1
        }
      ]
    }
  ],
  "nextSteps": [
    "Verificar configuração do pool de conexões do payment-service",
    "Analisar slow queries no banco durante o período de degradação"
  ]
}
```

### 5.3 Chat Contextual

Após o diagnóstico, o desenvolvedor pode fazer perguntas sobre o run específico. Cada turno do chat recebe:

- O diagnóstico já gerado como contexto base
- Os dados agregados do run
- O histórico de mensagens da sessão

Isso permite perguntas como:
- *"Quais endpoints foram mais afetados?"*
- *"A partir de quantos usuários o sistema começou a degradar?"*
- *"O checkout-service tem SLO de 200ms — ele foi respeitado?"*

---

## 6. Stack Tecnológico

### 6.1 Backend

| Tecnologia | Justificativa |
|---|---|
| **Go** | Performance, tipagem forte, goroutines nativas para processamento assíncrono |
| **SQLite** (PoC) | Zero configuração, fácil migração para PostgreSQL |
| **Anthropic SDK** | Direto, sem abstrações desnecessárias no MVP |
| **csv-parse** (Go stdlib) | Parsing nativo de CSV sem dependências externas |

### 6.2 Frontend

| Tecnologia | Justificativa |
|---|---|
| **Next.js** | Ecosystem consolidado, roteamento, SSR quando necessário |
| **Tailwind + shadcn/ui** | Visual profissional com velocidade de desenvolvimento |
| **TanStack Query** | Polling de status do run, cache de estado de servidor |
| **React Hook Form + Zod** | Formulários de cadastro com validação tipada |

---

## 7. API — Contrato HTTP

### Suites

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/suites` | Lista todas as suites |
| `POST` | `/suites` | Cria uma nova suite |
| `GET` | `/suites/:id` | Detalhe de uma suite com seus microsserviços |
| `PUT` | `/suites/:id` | Atualiza uma suite |
| `DELETE` | `/suites/:id` | Remove uma suite |

### Microsserviços

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/suites/:id/microservices` | Adiciona microsserviço à suite |
| `PUT` | `/microservices/:id` | Atualiza um microsserviço |
| `DELETE` | `/microservices/:id` | Remove um microsserviço |

### Runs

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/suites/:id/runs` | Lista runs de uma suite |
| `POST` | `/suites/:id/runs` | Cria run e faz upload do .jtl |
| `GET` | `/runs/:id` | Detalhe do run + status |
| `GET` | `/runs/:id/diagnosis` | Diagnóstico gerado |
| `POST` | `/runs/:id/chat` | Envia mensagem no chat do run |

---

## 8. Roadmap

### PoC — Fase 1 (objetivo imediato)

Fluxo completo de ponta a ponta com análise baseada apenas no `.jtl` e no contexto cadastrado.

- [ ] CRUD de Test Suites e Microservices
- [ ] Upload e parse de arquivo `.jtl` (CSV)
- [ ] Agregação de métricas por endpoint
- [ ] Integração com Claude para diagnóstico estruturado
- [ ] Relatório de diagnóstico no frontend
- [ ] Chat contextual sobre o run
- [ ] Histórico de runs por suite
- [ ] Processamento assíncrono com polling de status

### Fase 2 — Correlação Real

Integração com fontes externas de dados para diagnósticos com evidências mais profundas.

- [ ] Coleta de métricas Prometheus por janela de tempo do teste
- [ ] Eventos do Kubernetes: restarts, OOMKilled, CrashLoopBackOff
- [ ] Agregação de logs das aplicações durante o teste
- [ ] Correlation engine: cruzamento de métricas com timeline do JTL
- [ ] Detecção de baseline (métricas em condição normal vs sob carga)

### Fase 3 — Memória e Comparação

Inteligência acumulada entre execuções.

- [ ] Comparação automática entre runs da mesma suite
- [ ] Detecção de regressão: "a latência do checkout piorou 40% vs run anterior"
- [ ] Histórico de gargalos identificados por microsserviço
- [ ] Memória institucional: a IA sabe que o payment-service já teve problema de pool antes

### Fase 4 — Análise de Código

O nível mais profundo de diagnóstico.

- [ ] RAG sobre repositórios dos microsserviços
- [ ] Integração com GitHub/GitLab
- [ ] IA identifica trechos de código relacionados aos gargalos encontrados
- [ ] Sugestões de refatoração com contexto de código + métricas

### Fase 5 — Automação e Observação Passiva

- [ ] Sugestão de remediações automáticas (com aprovação humana)
- [ ] Modo de observação em produção: detecta padrões similares aos stress tests
- [ ] Runbooks automatizados para problemas recorrentes
- [ ] Alertas proativos baseados em padrões históricos

---

## 9. Decisões de Arquitetura e Justificativas

| Decisão | Alternativa Considerada | Justificativa |
|---|---|---|
| Monolito modular | Microsserviços | PoC não justifica overhead de infra distribuída. Módulos bem definidos facilitam extração futura |
| `registry` como orquestrador | Módulos peer-to-peer | Ponto único de controle simplifica rastreabilidade e evita dependências circulares |
| `runner` e `analyst` sem acesso ao banco | Acesso direto | Motores puros são testáveis isoladamente e sem efeitos colaterais |
| Dados passados por valor entre módulos | Busca por ID | Elimina lookups desnecessários e desacopla módulos do banco |
| Assíncrono com goroutine + polling | Síncrono / WebSocket | Simples no MVP, evolutivo para fila sem mudança de contrato |
| SQLite no MVP | PostgreSQL | Zero configuração, suficiente para PoC, migração trivial |
| Anthropic direto | LangChain / LangGraph | Reduz complexidade no MVP. Frameworks de orquestração entram nas fases 2-3 |

---

## 10. Gaps Conhecidos e Riscos

### Gaps a resolver nas próximas fases

- **Baseline de métricas**: sem dados em condição normal, é impossível saber se 70% de CPU é problemático para um serviço específico
- **Correlação temporal precisa**: o JTL e o Prometheus precisam ter timestamps sincronizados para correlação confiável
- **Logs não estruturados**: a maioria dos microsserviços não produz logs estruturados, o que complica extração automatizada
- **Metadados de fase do teste**: ramp-up, plateau e ramp-down têm pesos diferentes na análise — o JTL não distingue isso automaticamente

### Riscos da análise de código (Fase 4)

- **Falsos positivos**: a IA pode apontar trechos de código como gargalo quando o problema real é de infraestrutura
- **Contexto insuficiente**: diagnóstico de código precisa de código + configuração + queries + índices do banco para ser preciso
- **Acesso seguro a repositórios**: requer integração cuidadosa com controles de acesso do GitHub/GitLab

### Riscos de adoção

- Se o sistema exigir configuração excessiva, o desenvolvedor continuará usando o JMeter por inércia
- O cadastro de microsserviços precisa ser rápido e com valor imediato perceptível desde o primeiro uso
- A qualidade do diagnóstico na PoC vai determinar a aprovação — um prompt genérico compromete todo o projeto
