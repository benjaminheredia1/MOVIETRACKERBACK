# 🎬 MovieTracker — Backend

Tracker personal de películas y series construido con **Go**, **PostgreSQL**, **Redis** y **OpenAI**. Permite buscar contenido en TMDB, gestionar una watchlist personal con calificaciones, y chatear con un agente de IA que conoce tu historial mediante **function calling**.

---

## 📋 Requisitos Previos

- **Docker** y **Docker Compose** (v2+) instalados
- **API Key de TMDB** — Registro gratuito en [themoviedb.org](https://www.themoviedb.org)
- **API Key de OpenAI** — Proporcionada por la empresa

---

## ⚙️ Variables de Entorno

Copiar el archivo de ejemplo y configurar con valores reales:

```bash
cp .env.example .env
```

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `host` | Host de PostgreSQL | `postgres` (Docker) / `localhost` (local) |
| `port` | Puerto de PostgreSQL | `5432` |
| `user` | Usuario de PostgreSQL | `sistema` |
| `password` | Contraseña de PostgreSQL | `sistema` |
| `dbname` | Nombre de la base de datos | `sistema` |
| `app_port` | Puerto del servidor HTTP | `8080` |
| `TMDB_API_KEY` | Bearer token de TMDB (v4 auth) | `eyJhbGci...` |
| `TMDB_BASE_URL` | URL base de la API TMDB | `https://api.themoviedb.org/3` |
| `REDIS_HOST` | Host de Redis | `redis` (Docker) / `localhost` (local) |
| `REDIS_PORT` | Puerto de Redis | `6379` |
| `REDIS_PASSWORD` | Contraseña de Redis | (vacío por defecto) |
| `OPENAI_API_KEY` | API Key de OpenAI | `sk-proj-...` |
| `OPENAI_MODEL` | Modelo de OpenAI | `gpt-4o-mini` |

> **⚠️ Importante:** El archivo `.env` con valores reales **NO debe subirse** al repositorio. Solo se versiona `.env.example`.

---

## 🚀 Levantar el Proyecto

### Con Docker Compose (recomendado)

```bash
# 1. Configurar variables de entorno
cp MovieTrackerBack/.env.example MovieTrackerBack/.env
# Editar MovieTrackerBack/.env con las API keys reales

# 2. Levantar todos los servicios
cd MovieTrackerFront
docker-compose up --build
```

Esto levanta **4 servicios**:

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| `postgres` | 5432 | PostgreSQL 16 Alpine |
| `redis` | 6379 | Redis 7 Alpine |
| `backend` | 8080 | API REST en Go (Gin) |
| `frontend` | 3000 | React + Nginx (proxy reverso) |

Acceder a la aplicación: **http://localhost:3000**

Las migraciones de base de datos se ejecutan **automáticamente** al iniciar el backend usando `golang-migrate`.

### Desarrollo local (sin Docker)

Requiere PostgreSQL y Redis corriendo localmente:

```bash
cd MovieTrackerBack
go run cmd/main.go
```

---

## 🗄️ Esquema de Base de Datos

El esquema consta de **3 tablas** con relaciones entre ellas, gestionadas mediante migraciones SQL con `golang-migrate`. No se usa ORM — todo el SQL está escrito manualmente.

```
┌──────────────────────────┐     ┌───────────────────┐     ┌─────────────────────┐
│          ITEMS           │     │    LIST_ITEMS      │     │        LISTA        │
├──────────────────────────┤     ├───────────────────┤     ├─────────────────────┤
│ id SERIAL (PK)           │◄────│ item_id INT (FK)  │     │ id SERIAL (PK)      │
│ tmdb_id INT NOT NULL     │     │ list_id INT (FK)  │────►│ name VARCHAR(255)   │
│ adult BOOLEAN            │     │ added_at TIMESTAMP│     │ description TEXT     │
│ backdrop_path TEXT       │     │ id SERIAL (PK)    │     │ created_at TIMESTAMP│
│ name VARCHAR(255)        │     └───────────────────┘     └─────────────────────┘
│ original_name VARCHAR    │
│ overview TEXT             │
│ poster_path TEXT         │
│ media_type VARCHAR(10)   │     Valores: 'movie' | 'tv'
│ original_language VARCHAR│
│ popularity DECIMAL(10,4) │
│ first_air_date DATE      │
│ softcore BOOLEAN         │
│ genre_ids TEXT            │     IDs de género como cadena
│ origin_country TEXT      │     Países como cadena
│ vote_average DECIMAL(3,1)│
│ vote_count INT           │
│ list_id INT (FK, NULL)   │     Referencia opcional a LISTA
│ status VARCHAR(10)       │     'pending' | 'watched'
│ comentary_user TEXT      │     Comentario al marcar como visto
│ calification_user FLOAT  │     Nota del 1 al 10
│ watched_at TIMESTAMP     │     Fecha en que se marcó como visto
│ added_at TIMESTAMP       │     DEFAULT NOW()
└──────────────────────────┘
```

### Relaciones

- `LIST_ITEMS.item_id` → `ITEMS.id` (`ON DELETE CASCADE`)
- `LIST_ITEMS.list_id` → `LISTA.id` (`ON DELETE CASCADE`)

La tabla `LIST_ITEMS` actúa como tabla pivote para la relación muchos-a-muchos entre items y listas. Los `CASCADE` garantizan integridad referencial: al eliminar una lista o un item, las asociaciones se eliminan automáticamente.

### Migraciones

Se utilizan 3 archivos de migración ejecutados en orden con `golang-migrate`:

| Archivo | Descripción |
|---------|-------------|
| `000001_init_schema.up.sql` | Crea la tabla `ITEMS` con todos los campos de TMDB y seguimiento |
| `000002_add_lista.up.sql` | Crea la tabla `LISTA` para listas personalizadas |
| `000003_add_list_items.up.sql` | Crea la tabla `LIST_ITEMS` (pivote) con foreign keys |

---

## 📡 Endpoints de la API

Base URL: `/api`

### TMDB — Búsqueda y Tendencias

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/api/search?query={texto}` | Buscar películas y series en TMDB (multi search). Resultados cacheados en Redis |
| `GET` | `/api/recomendations` | Obtener personas populares de TMDB |

### Items — Watchlist

| Método | Endpoint | Descripción | Request Body |
|--------|----------|-------------|-------------|
| `GET` | `/api/items` | Listar items con filtros opcionales | — |
| `GET` | `/api/items/:id` | Obtener un item por ID | — |
| `POST` | `/api/items` | Agregar item a la watchlist | `{ tmdb_id, name, media_type, ... }` |
| `PATCH` | `/api/items/:id/watched` | Marcar como visto con calificación | `{ rating: 1-10, commentary: "..." }` |
| `DELETE` | `/api/items/:id` | Eliminar item de la watchlist | — |

**Filtros disponibles en `GET /api/items`** (query params):

| Parámetro | Valores posibles | Ejemplo |
|-----------|-----------------|---------|
| `status` | `pending`, `watched` | `?status=pending` |
| `media_type` | `movie`, `tv` | `?media_type=movie` |
| `order_by` | `added_at`, `name`, `vote_average`, `popularity` | `?order_by=name` |
| `order_dir` | `ASC`, `DESC` | `?order_dir=ASC` |

**Validaciones:**

- `POST /api/items` → `tmdb_id` es requerido, `media_type` debe ser `"movie"` o `"tv"`
- `PATCH /api/items/:id/watched` → `rating` debe estar entre 1 y 10. Se registra `watched_at = NOW()` automáticamente
- Todos los endpoints con `:id` validan que sea un entero válido

### Listas Personalizadas

| Método | Endpoint | Descripción | Request Body |
|--------|----------|-------------|-------------|
| `GET` | `/api/lists` | Listar todas las listas | — |
| `GET` | `/api/lists/:id` | Obtener una lista por ID | — |
| `POST` | `/api/lists` | Crear nueva lista | `{ name: "...", description: "..." }` |
| `PUT` | `/api/lists/:id` | Actualizar una lista | `{ name: "...", description: "..." }` |
| `DELETE` | `/api/lists/:id` | Eliminar una lista (cascade) | — |

### Chat — Agente IA

| Método | Endpoint | Descripción | Request Body |
|--------|----------|-------------|-------------|
| `POST` | `/api/chat` | Enviar mensaje al agente | `{ session_id: "...", role: "user", content: "..." }` |

**Response:** `{ "reply": "Respuesta del agente en Markdown..." }`

---

## 🤖 Agente de Chat con OpenAI

### Descripción

El agente es un asistente personal experto en cine y series, integrado con la base de datos del usuario mediante **function calling (tools)** del SDK oficial `openai-go`. El agente **NO** recibe toda la watchlist en el system prompt — usa tools para consultar la base de datos solo cuando lo necesita, cumpliendo el requisito técnico obligatorio.

### System Prompt

```
Eres un asistente personal experto en películas y series, integrado en una aplicación
estilo Letterboxd/Trakt. Tienes acceso directo a la base de datos personal del usuario
mediante herramientas (tools). Tu trabajo es responder a las consultas del usuario
usando esta información y tu propio conocimiento sobre cine y televisión.

Reglas:
- Si te preguntan por las películas o series que el usuario ha guardado, usa "obtener_items".
- Si te preguntan por las listas personalizadas, usa "obtener_listas".
- Siempre responde de forma amigable, útil, y en formato Markdown.
- Da recomendaciones basadas en los géneros o películas que ya han visto.
```

### Tools Definidas

| Tool | Tipo | Parámetros | Descripción |
|------|------|-----------|-------------|
| `obtener_items` | `function` | Ninguno | Consulta la tabla `ITEMS` vía `ItemRepository.GetAll()` y retorna **todas** las películas y series del usuario, incluyendo estado (`pending`/`watched`), calificaciones (`calification_user`) y comentarios (`comentary_user`) |
| `obtener_listas` | `function` | Ninguno | Consulta la tabla `LISTA` vía `ListRepository.GetAll()` y retorna **todas** las listas personalizadas del usuario con nombres, descripciones y fechas de creación |

**Cuándo las usa el agente:**

- `obtener_items` → Cuando el usuario pregunta sobre su watchlist ("¿qué tengo pendiente?"), pide recomendaciones basadas en gustos ("recomiéndame algo"), consulta estadísticas ("¿cuántas películas he visto?"), o pregunta por calificaciones ("¿cuáles son mis favoritas?")
- `obtener_listas` → Cuando el usuario pregunta por sus listas ("¿qué listas tengo?") o quiere organizar contenido

### Flujo de Funcionamiento

```
Usuario                     Backend                      OpenAI                    PostgreSQL
  │                           │                            │                          │
  │─── POST /api/chat ───────►│                            │                          │
  │    {session_id, content}  │                            │                          │
  │                           │── Leer historial Redis ──► │                          │
  │                           │◄── historial previo ───── │                          │
  │                           │                            │                          │
  │                           │── messages + tools ──────► │                          │
  │                           │                            │                          │
  │                           │◄── tool_call: obtener_items│                          │
  │                           │                            │                          │
  │                           │──── SELECT * FROM items ──────────────────────────► │
  │                           │◄─── items[] (JSON) ───────────────────────────────── │
  │                           │                            │                          │
  │                           │── tool_result (items) ───► │                          │
  │                           │                            │                          │
  │                           │◄── respuesta final ─────── │                          │
  │                           │                            │                          │
  │                           │── Guardar historial Redis ►│                          │
  │◄── {reply: "..."} ────── │                            │                          │
```

**Loop de tools:** Si OpenAI solicita ejecutar una tool, el backend la ejecuta, envía el resultado, y espera una nueva respuesta. Este proceso se repite en un `for` loop hasta que el agente genera una respuesta final sin más tool calls.

### Historial de Conversación

- Se almacena en **Redis** con la clave `chat_{session_id}` y **TTL de 2 horas**
- El `session_id` es generado por el frontend al iniciar una nueva conversación
- Al iniciar cada request, el historial se recupera de Redis y se reconstruyen los mensajes para OpenAI
- Al finalizar, se actualiza el historial con el mensaje del usuario y la respuesta del agente
- **No persiste entre sesiones del navegador** — suficiente para una sesión de uso activo

### Manejo del Tool Call ID

El SDK de OpenAI puede generar `tool_call_id` que excedan los 64 caracteres permitidos por la API. Se implementó una **truncación automática** a 64 caracteres para evitar errores `400 Bad Request`.

---

## 🏗️ Decisiones de Arquitectura

### 1. Arquitectura Hexagonal (Ports & Adapters)

**Decisión:** Organizar el código en capas `domain`, `application` e `infrastructure`.

```
internal/
├── domain/              → Modelos y puertos (interfaces)
│   ├── itemslist.go     → ITEM, LISTA, LISTA_ITEM, Filters
│   ├── chat.go          → ChatMessage
│   ├── tmdb.go          → MediaResult (respuesta TMDB)
│   └── ports.go         → Interfaces: ITEMREPOSITORY, LISTAREPOSITORY, ChatRepository, MediaRepository
│
├── application/         → Servicios (lógica de negocio)
│   ├── items_service.go → Validaciones (tmdb_id, media_type, rating 1-10)
│   ├── list_service.go  → CRUD de listas
│   └── chat_service.go  → Orquestación del chat
│
└── infrastructure/      → Adaptadores (implementaciones concretas)
    ├── http/            → Router Gin + Handlers REST
    ├── postgres/        → Repositorios SQL, conexión, migraciones
    ├── redis/           → Cliente Redis + CacheRepository
    ├── tmdb/            → Cliente HTTP para TMDB (go-resty)
    └── openai/          → Cliente OpenAI con function calling
```

**Por qué:** La capa de dominio define interfaces (ports) que las implementaciones concretas (adapters) satisfacen. Los servicios de aplicación dependen de interfaces, no de PostgreSQL o Redis directamente. Esto permite cambiar cualquier infraestructura sin afectar la lógica de negocio, y facilita el testing con mocks.

### 2. SQL Manual sin ORM

**Decisión:** Todo el acceso a datos se realiza con `database/sql` de la biblioteca estándar de Go.

**Por qué:**
- **Requisito explícito** de la prueba: "no se permite el uso de ORMs"
- **Control total sobre las queries:** los filtros y ordenamiento se manejan en SQL, no filtrando en memoria
- **Performance predecible:** sin capas de abstracción que generen queries subóptimas
- **Transparencia:** cada query es visible y auditable en el repositorio correspondiente

### 3. Inyección de Dependencias Manual en `main.go`

**Decisión:** Construir todo el grafo de dependencias explícitamente en el entrypoint.

```go
// cmd/main.go
itemRepo := postgres.NewItemRepository(db)
itemService := application.NewItemsService(itemRepo, tmdbClient)
itemsHandler := handlers.NewItemsHandler(itemService)
```

**Por qué:** El grafo es lo suficientemente simple para no justificar un framework DI. Todo el wiring es visible en un solo archivo, sin "magia" ni reflexión.

### 4. Estrategia de Caché con Redis

**Decisión:** Redis como caché exclusiva con TTLs diferenciados:

| Dato | TTL | Justificación |
|------|-----|---------------|
| Búsquedas TMDB | 24 horas | Los resultados de búsqueda de TMDB cambian con muy poca frecuencia. Un TTL de 24h reduce significativamente las llamadas a la API externa sin sacrificar frescura. La query se normaliza a minúsculas para evitar duplicados en caché |
| Historial de chat | 2 horas | Suficiente para una sesión de uso activo. Después de 2h sin actividad, es razonable asumir que la sesión terminó y la memoria puede liberarse |

**Por qué Redis no es base de datos:** Redis es volátil por naturaleza. Si se reinicia, se pierden datos. La watchlist y las listas del usuario se persisten en PostgreSQL, que garantiza durabilidad.

### 5. Frontend como Servicio Separado vía Nginx

**Decisión:** El frontend React se sirve como un **contenedor independiente con Nginx**, no embebido en el proceso de Go.

**Por qué:**
- **Separación de responsabilidades:** Go se enfoca en la API; Nginx sirve archivos estáticos de forma óptima
- **Proxy reverso integrado:** Nginx redirige `/api/*` al backend, eliminando CORS sin middleware adicional
- **Rendimiento:** Nginx maneja compresión, cache headers y keep-alive de forma nativa
- **Timeout para OpenAI:** `proxy_read_timeout` de 120s para peticiones al chat que involucran function calling (pueden tardar varios segundos)
- **Escalabilidad:** en producción, frontend y backend pueden escalar independientemente

### 6. Gin como Framework HTTP

**Decisión:** Usar `gin-gonic/gin` v1.12 como router HTTP.

**Por qué:**
- Framework HTTP más adoptado del ecosistema Go, con documentación extensa
- Binding automático de JSON (`ShouldBindJSON`), manejo de parámetros de ruta/query
- Grupos de rutas (`r.Group("/api")`) para organización limpia
- Rendimiento probado a escala

---

## 🛠️ Stack Técnico

| Tecnología | Versión | Rol |
|-----------|---------|-----|
| Go | 1.25 | Lenguaje del backend |
| Gin | 1.12 | Framework HTTP / Router |
| PostgreSQL | 16 | Base de datos principal (SQL manual, sin ORM) |
| Redis | 7 | Caché (búsquedas TMDB, historial chat) |
| OpenAI SDK (`openai-go`) | 1.12 | Agente de chat con function calling |
| go-resty | 2.17 | Cliente HTTP para la API de TMDB |
| golang-migrate | 4.19 | Migraciones de base de datos |
| godotenv | 1.5 | Carga de variables de entorno desde `.env` |
| Docker + Compose | — | Contenedores y orquestación |
| React + Vite | 19 + 8 | Frontend SPA |
| Nginx | Alpine | Proxy reverso + servidor de estáticos |
