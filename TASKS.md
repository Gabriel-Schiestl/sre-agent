# Stress Test Intelligence Platform — Task Board

> Estratégia: **Frontend first**, completo com chamadas de API prontas para o backend.
> Após frontend finalizado, backend é construído módulo a módulo.

---

## Legenda

- `[ ]` — Pendente
- `[x]` — Concluído
- `[~]` — Em progresso
- `[!]` — Bloqueado
- `[-]` — Suspenso (fora do escopo atual)

---

## Mapa de Rotas do Frontend

```
/suites                        → listagem de todas as suites
/suites/[id]                   → detalhe da suite
                                   ├── info da suite (nome, descrição)
                                   ├── seção de microsserviços (lista + add/edit/delete via modal)
                                   └── listagem de runs com status
/runs/[id]                     → detalhe do run
                                   ├── status (pending / analyzing / done / failed)
                                   ├── diagnóstico (quando done)
                                   └── chat contextual (quando done)
```

---

## FASE 1 — FRONTEND

---

### EPIC F0 — Setup do Projeto

> Fundação. Bloqueante para todos os outros epics de frontend.

- [x] **F0.1** — Criar projeto Next.js com App Router
  ```bash
  npx create-next-app@latest packages/ui --typescript --tailwind --app --src-dir
  ```

- [x] **F0.2** — Instalar e configurar shadcn/ui
  ```bash
  npx shadcn-ui@latest init
  ```
  Componentes necessários: `button`, `input`, `textarea`, `label`, `dialog`, `form`, `badge`, `card`, `separator`, `skeleton`, `sonner`, `table`, `select`

- [x] **F0.3** — Instalar dependências
  ```bash
  npm install @tanstack/react-query axios react-hook-form zod @hookform/resolvers
  npm install @tanstack/react-query-devtools --save-dev
  ```

- [x] **F0.4** — Configurar variável de ambiente
  Criar `frontend/.env.local` e `frontend/.env.example`:
  ```
  NEXT_PUBLIC_API_URL=http://localhost:8080
  ```

- [x] **F0.5** — Criar API client base (`src/lib/api.ts`)
  - Instância `axios` com `baseURL` de `NEXT_PUBLIC_API_URL`
  - Interceptor de erro: extrai mensagem da resposta e relança como `Error`
  - Todas as funções de acesso à API ficam aqui — nenhum componente usa axios diretamente

- [x] **F0.6** — Configurar TanStack Query Provider
  - `src/providers/query-provider.tsx` com `QueryClient` e `QueryClientProvider`
  - Envolver `layout.tsx` raiz com o provider
  - `ReactQueryDevtools` apenas em `development`

- [x] **F0.7** — Criar estrutura de pastas
  ```
  src/
  ├── app/
  │   ├── suites/
  │   │   ├── page.tsx              ← /suites
  │   │   └── [id]/
  │   │       └── page.tsx          ← /suites/[id]
  │   └── runs/
  │       └── [id]/
  │           └── page.tsx          ← /runs/[id]
  ├── components/
  │   ├── ui/                       ← gerado pelo shadcn
  │   └── shared/                   ← componentes do projeto
  ├── hooks/                        ← custom hooks com TanStack Query
  ├── lib/
  │   └── api.ts
  ├── types/                        ← tipos TypeScript dos domínios
  └── providers/
  ```

- [x] **F0.8** — Criar layout base (`src/app/layout.tsx`)
  - Sidebar simples com link para "Suites" e nome da plataforma no topo
  - Área de conteúdo principal com `{children}`
  - Redirecionar `/` para `/suites`

---

### EPIC F1 — Tipos TypeScript

> Define os contratos alinhados com a API do backend. Feito antes de qualquer hook ou componente.

- [x] **F1.1** — `src/types/suite.ts`
  ```ts
  export interface TestSuite {
    id: string
    name: string
    description: string
    createdAt: string
    updatedAt: string
  }

  export interface CreateSuitePayload { name: string; description: string }
  export interface UpdateSuitePayload { name: string; description: string }
  ```

- [x] **F1.2** — `src/types/microservice.ts`
  ```ts
  export interface Microservice {
    id: string
    testSuiteId: string
    name: string
    description: string
    language: string
    mainEndpoints: string[]
    cpuLimit?: string
    memoryLimit?: string
    sloLatencyP99Ms?: number
    sloErrorRatePct?: number
    createdAt: string
  }

  export interface CreateMicroservicePayload {
    name: string
    description: string
    language: string
    mainEndpoints: string[]
    cpuLimit?: string
    memoryLimit?: string
    sloLatencyP99Ms?: number
    sloErrorRatePct?: number
  }
  ```

- [x] **F1.3** — `src/types/run.ts`
  ```ts
  export type RunStatus = 'pending' | 'analyzing' | 'done' | 'failed'

  export interface TestRun {
    id: string
    testSuiteId: string
    name: string
    virtualUsers: number
    durationSeconds: number
    notes: string
    status: RunStatus
    createdAt: string
  }
  ```

- [x] **F1.4** — `src/types/diagnosis.ts`
  ```ts
  export interface ErrorCategory {
    category: string
    description: string
    occurrences: number
    affectedEndpoints: string[]
    severity: 'low' | 'medium' | 'high' | 'critical'
  }

  export interface Hypothesis {
    title: string
    evidence: string
    priority: number
  }

  export interface Bottleneck {
    microservice: string
    confidence: 'low' | 'medium' | 'high'
    hypotheses: Hypothesis[]
  }

  export interface Diagnosis {
    id: string
    testRunId: string
    errorPlan: ErrorCategory[]
    bottlenecks: Bottleneck[]
    nextSteps: string[]
    createdAt: string
  }

  export interface ChatMessage {
    role: 'user' | 'assistant'
    content: string
  }
  ```

---

### EPIC F2 — Camada de API e Hooks

> Toda a lógica de comunicação com o backend centralizada aqui.

- [x] **F2.1** — Funções de API em `src/lib/api.ts`

  **Suites:**
  - `getSuites()` → `GET /suites`
  - `getSuite(id)` → `GET /suites/:id` (retorna suite + microservices)
  - `createSuite(payload)` → `POST /suites`
  - `updateSuite(id, payload)` → `PUT /suites/:id`
  - `deleteSuite(id)` → `DELETE /suites/:id`

  **Microservices:**
  - `createMicroservice(suiteId, payload)` → `POST /suites/:id/microservices`
  - `updateMicroservice(id, payload)` → `PUT /microservices/:id`
  - `deleteMicroservice(id)` → `DELETE /microservices/:id`

  **Runs:**
  - `getRuns(suiteId)` → `GET /suites/:id/runs`
  - `getRun(id)` → `GET /runs/:id`
  - `createRun(suiteId, formData)` → `POST /suites/:id/runs` (multipart/form-data com arquivo `.jtl`)
  - `getDiagnosis(runId)` → `GET /runs/:id/diagnosis`
  - ~~`sendChatMessage(runId, message)` → `POST /runs/:id/chat`~~ (suspenso)

- [x] **F2.2** — Hooks de suites (`src/hooks/use-suites.ts`)
  - `useSuites()` — lista todas, `queryKey: ['suites']`
  - `useSuite(id)` — detalhe com microservices, `queryKey: ['suites', id]`
  - `useCreateSuite()` — mutation, invalida `['suites']` on success
  - `useUpdateSuite()` — mutation, invalida `['suites']` e `['suites', id]`
  - `useDeleteSuite()` — mutation, invalida `['suites']`

- [x] **F2.3** — Hooks de microservices (`src/hooks/use-microservices.ts`)
  - `useCreateMicroservice(suiteId)` — invalida `['suites', suiteId]`
  - `useUpdateMicroservice(suiteId)` — invalida `['suites', suiteId]`
  - `useDeleteMicroservice(suiteId)` — invalida `['suites', suiteId]`

- [x] **F2.4** — Hooks de runs (`src/hooks/use-runs.ts`)
  - `useRuns(suiteId)` — `queryKey: ['runs', suiteId]`
  - `useRun(id)` — `queryKey: ['runs', id]`
  - `useRunStatus(id)` — polling a cada 3s, para quando `status === 'done' || 'failed'`
  - `useCreateRun(suiteId)` — mutation, invalida `['runs', suiteId]`

- [x] **F2.5** — Hooks de diagnóstico (`src/hooks/use-diagnosis.ts`)
  - `useDiagnosis(runId)` — `queryKey: ['diagnosis', runId]`

- [-] **F2.6** — Hook de chat (`useSendChatMessage`) — suspenso

---

### EPIC F3 — Página `/suites`

- [x] **F3.1** — Criar `src/app/suites/page.tsx`
  - Título "Test Suites" + botão "Nova Suite" no canto direito
  - Usa `useSuites()` para carregar dados
  - Estado de loading: skeleton de linhas
  - Estado vazio: mensagem com call-to-action para criar a primeira suite

- [x] **F3.2** — Componente de tabela de suites (`src/components/shared/suites-table.tsx`)
  - Colunas: Nome, Descrição, Criado em, Ações
  - Linha inteira clicável → navega para `/suites/:id`
  - Botão de deletar abre dialog de confirmação antes de executar
  - Loading state no botão de delete durante a mutation

- [x] **F3.3** — Dialog de criação de suite (`src/components/shared/create-suite-dialog.tsx`)
  - Campos: Nome (required, mín. 3 chars) e Descrição (required)
  - Validação com Zod + React Hook Form
  - Submit com loading state, fecha e invalida query on success, toast de erro on failure

---

### EPIC F4 — Página `/suites/[id]`

> Uma única página com três seções: info da suite, microsserviços e runs.

- [x] **F4.1** — Criar `src/app/suites/[id]/page.tsx`
  - Breadcrumb: Suites > [nome da suite]
  - Header com nome, descrição e botão "Editar Suite"
  - Seção de Microsserviços abaixo do header
  - Seção de Test Runs abaixo dos microsserviços

- [x] **F4.2** — Dialog de edição de suite (`src/components/shared/edit-suite-dialog.tsx`)
  - Pré-populado com dados atuais
  - Mesma validação do create

- [x] **F4.3** — Seção de microsserviços (`src/components/shared/microservices-section.tsx`)
  - Header "Microsserviços" + botão "Adicionar"
  - Lista de cards: nome, linguagem (badge), endpoints principais resumidos
  - SLOs exibidos se definidos (ex: "p99 < 200ms")
  - Botão de editar e deletar por card
  - Estado vazio: "Nenhum microsserviço cadastrado"

- [x] **F4.4** — Formulário de microsserviço (`src/components/shared/microservice-form.tsx`)
  - Reutilizado em create e edit
  - Campos obrigatórios: Nome, Descrição, Linguagem (select: Go, Java, Node.js, Python, .NET, Outro)
  - Endpoints: campo de lista dinâmica com add/remove (array de inputs)
  - Seção "Configurações opcionais" collapsível: CPU Limit, Memory Limit, SLO latência p99 (ms), SLO error rate (%)
  - Validação Zod completa

- [x] **F4.5** — Dialog de criação de microsserviço (`src/components/shared/create-microservice-dialog.tsx`)

- [x] **F4.6** — Dialog de edição de microsserviço (`src/components/shared/edit-microservice-dialog.tsx`)

- [x] **F4.7** — Seção de runs (`src/components/shared/runs-section.tsx`)
  - Header "Test Runs" + botão "Novo Run"
  - Tabela: Nome, Usuários Virtuais, Duração, Status (badge), Criado em
  - Badge de status: `pending` (cinza), `analyzing` (amarelo animado), `done` (verde), `failed` (vermelho)
  - Linha clicável → navega para `/runs/:id`
  - Estado vazio: "Nenhum run executado ainda"

- [x] **F4.8** — Dialog de criação de run (`src/components/shared/create-run-dialog.tsx`)
  - Campos: Nome, Usuários Virtuais (number > 0), Duração em segundos (number > 0), Notas (textarea)
  - Upload de arquivo `.jtl`: input file com validação de extensão
  - Preview do arquivo selecionado (nome + tamanho formatado)
  - On submit: envia via multipart/form-data, redireciona para `/runs/:id` on success

---

### EPIC F5 — Página `/runs/[id]`

- [x] **F5.1** — Criar `src/app/runs/[id]/page.tsx`
  - Breadcrumb: Suites > [nome da suite] > [nome do run]
  - Metadados do run: usuários virtuais, duração, notas, data de criação
  - Delega renderização principal para `<RunStatusView>`

- [x] **F5.2** — Componente de status do run (`src/components/shared/run-status-view.tsx`)
  - Consome `useRunStatus(runId)` com polling automático
  - `pending`: ícone de relógio + "Aguardando início do processamento..."
  - `analyzing`: spinner + "Analisando resultados com IA..." + barra de progresso indeterminada
  - `done`: renderiza `<DiagnosisView runId={runId} />`
  - `failed`: mensagem de erro com orientação de próximos passos

---

### EPIC F6 — Diagnóstico

- [x] **F6.1** — Componente container (`src/components/shared/diagnosis-view.tsx`)
  - Consome `useDiagnosis(runId)` para carregar o diagnóstico
  - Renderiza as três seções: categorias de erro, gargalos e próximos passos

- [x] **F6.2** — Seção de categorias de erro (`src/components/shared/error-plan-section.tsx`)
  - Cards por categoria de erro
  - Badge de severidade com cor: `critical` vermelho, `high` laranja, `medium` amarelo, `low` cinza
  - Endpoints afetados como chips
  - Contagem de ocorrências em destaque

- [x] **F6.3** — Seção de gargalos (`src/components/shared/bottlenecks-section.tsx`)
  - Um card por microsserviço suspeito
  - Header: nome do serviço + badge de confiança
  - Lista de hipóteses numeradas por prioridade (título + evidência)

- [x] **F6.4** — Seção de próximos passos (`src/components/shared/next-steps-section.tsx`)
  - Lista numerada de ações sugeridas

- [-] **F6.5** — Painel de chat — suspenso

---

### EPIC F7 — Polimento

- [x] **F7.1** — `Toaster` global (sonner) com feedback de sucesso/erro em todas as mutations

- [x] **F7.2** — Componente `ErrorState` reutilizável: mensagem + botão de retry para todas as queries que podem falhar

- [x] **F7.3** — Skeleton loaders em todas as seções com dados assíncronos

- [x] **F7.4** — Revisão de responsividade: viewport 1280px e 1440px

---

## Dependências entre Epics (Frontend)

```
F0 (setup)
 └── F1 (tipos)
      └── F2 (hooks de API)
           ├── F3 (página /suites)
           ├── F4 (página /suites/[id])       ← depende de F3
           └── F5 (página /runs/[id])
                └── F6 (diagnóstico e chat)   ← depende de F5
                     └── F7 (polimento)
```

---

## FASE 2 — BACKEND

---

### EPIC B0 — Setup do Projeto Backend

- [ ] **B0.1** — Inicializar módulo Go
  ```bash
  go mod init github.com/schiestl/sre-agent
  ```

- [ ] **B0.2** — Criar estrutura de pastas
  ```
  cmd/api/
  internal/registry/
  internal/runner/
  internal/analyst/
  pkg/llm/
  migrations/
  data/uploads/
  ```

- [ ] **B0.3** — Instalar dependências
  ```bash
  go get github.com/go-chi/chi/v5
  go get github.com/mattn/go-sqlite3
  go get github.com/golang-migrate/migrate/v4
  go get github.com/google/uuid
  go get github.com/anthropics/anthropic-sdk-go
  go get github.com/joho/godotenv
  ```

- [ ] **B0.4** — Criar `Makefile` com targets: `dev`, `build`, `test`, `migrate-up`, `migrate-down`

- [ ] **B0.5** — Criar `.env` e `.env.example`
  ```
  DATABASE_PATH=./data/sre-agent.db
  ANTHROPIC_API_KEY=sk-ant-...
  PORT=8080
  FRONTEND_URL=http://localhost:3000
  ```

---

### EPIC B1 — Modelo de Dados e Repositório

- [ ] **B1.1** — Migration `001_create_test_suites.sql`
- [ ] **B1.2** — Migration `002_create_microservices.sql` (com `main_endpoints` como JSON)
- [ ] **B1.3** — Migration `003_create_test_runs.sql` (status enum: pending/analyzing/done/failed)
- [ ] **B1.4** — Migration `004_create_diagnoses.sql` (error_plan, bottlenecks, next_steps como JSON; ~~chat_history~~ suspenso)
- [ ] **B1.5** — `internal/registry/model.go`: structs Go das 4 entidades
- [ ] **B1.6** — `internal/registry/repository.go`: CRUD completo de TestSuite, Microservice, TestRun
- [ ] **B1.7** — `internal/registry/repository.go`: `UpdateStatus`, `SaveDiagnosis`, `GetDiagnosis` (~~`AppendChatMessage`~~ suspenso)

---

### EPIC B2 — Módulo `runner`

- [ ] **B2.1** — `internal/runner/model.go`: contratos `RunPayload` e `AggregatedData`
- [ ] **B2.2** — `internal/runner/parser.go`: leitura do CSV JTL, campos relevantes (seção 4.2 do documento)
- [ ] **B2.3** — `internal/runner/aggregator.go`: p50/p90/p99 por endpoint, error rate, volume
- [ ] **B2.4** — `internal/runner/aggregator.go`: agrupamento de erros por responseCode + failureMessage
- [ ] **B2.5** — `internal/runner/aggregator.go`: timeline em janelas de 30s
- [ ] **B2.6** — `internal/runner/service.go`: orquestra parse → agregação, retorna `AggregatedData`
- [ ] **B2.7** — Testes unitários com fixture `testdata/sample.jtl`

---

### EPIC B3 — Pacote `/pkg/llm`

- [ ] **B3.1** — `pkg/llm/client.go`: struct `Client`, construtor `New(apiKey string)`
- [ ] **B3.2** — Método `Complete(ctx, systemPrompt string, messages []Message) (string, error)` usando `claude-sonnet-4-6`
- [ ] **B3.3** — Tratamento de erros da API: timeout, rate limit, resposta malformada

---

### EPIC B4 — Módulo `analyst`

- [ ] **B4.1** — `internal/analyst/model.go`: contratos `AnalysisPayload`, `Diagnosis` (~~`ChatPayload`~~ suspenso)
- [ ] **B4.2** — `internal/analyst/prompt.go`: builder do prompt completo (contexto de microsserviços + métricas + schema JSON de resposta)
- [ ] **B4.3** — `internal/analyst/service.go`: chama LLM, parseia JSON, valida estrutura, retorna `Diagnosis`
- [-] **B4.4** — `internal/analyst/service.go`: método `Chat` — suspenso

---

### EPIC B5 — Módulo `registry`

- [ ] **B5.1** — `internal/registry/service.go`: CRUD de suites e microsserviços (wraps repository)
- [ ] **B5.2** — `internal/registry/service.go`: `CreateRun` — salva arquivo `.jtl`, persiste run com `pending`, dispara goroutine
- [ ] **B5.3** — Goroutine de análise: `pending → analyzing → runner → analyst → done/failed`
- [-] **B5.4** — `internal/registry/service.go`: `SendChatMessage` — suspenso
- [ ] **B5.5** — `internal/registry/handler.go`: todos os handlers de Suites (GET list, POST, GET id, PUT, DELETE)
- [ ] **B5.6** — `internal/registry/handler.go`: handlers de Microservices (POST, PUT, DELETE)
- [ ] **B5.7** — `internal/registry/handler.go`: handlers de Runs (GET list, POST multipart, GET id, GET diagnosis) (~~POST chat~~ suspenso)

---

### EPIC B6 — Entrypoint

- [ ] **B6.1** — `cmd/api/main.go`: carrega env, inicializa DB, roda migrations, instancia e injeta todos os módulos
- [ ] **B6.2** — Configurar router `chi` com todos os grupos de rotas
- [ ] **B6.3** — Middlewares: CORS (origem do env), logger de requests, recoverer de panics
- [ ] **B6.4** — Criar `data/uploads/` no startup se não existir
- [ ] **B6.5** — Smoke test end-to-end: suite → microsserviço → upload JTL → polling → diagnóstico

---

## Dependências entre Epics (Backend)

```
B0 (setup)
 └── B1 (dados)
      ├── B2 (runner) ─── paralelo ─── B3 (pkg/llm)
      │                                      │
      │                               B4 (analyst)   ← depende de B2 e B3
      │                                      │
      └──────────────────────────────  B5 (registry)  ← depende de B1, B2, B4
                                             │
                                       B6 (entrypoint)
```

---

## Progresso Geral

| Fase | Epics | Tasks | Concluídas |
|------|-------|-------|------------|
| Frontend | F0–F7 | 38 | 36 |
| Backend | B0–B6 | 35 | 0 |
| **Total** | **15** | **73** | **36** |
