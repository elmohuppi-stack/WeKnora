# Built-in Models Management Guide

## Overview

Built-in models are system-level model configurations visible to all tenants. Sensitive information is hidden, and they cannot be edited or deleted. Built-in models are typically used to provide default model configurations, ensuring all tenants can use a unified model service.

## Built-in Model Features

- **Visible to all tenants**: Built-in models are visible to every tenant without individual configuration
- **Security protection**: Sensitive information (API Key, Base URL) of built-in models is hidden and cannot be viewed
- **Read-only protection**: Built-in models cannot be edited or deleted; they can only be set as the default model
- **Centralized management**: Maintained uniformly by the system administrator to ensure configuration consistency and security

## How to Add Built-in Models

WeKnora supports two ways to add built-in models: the **recommended** approach is YAML declarative configuration (with automatic idempotent application), while direct SQL insertion is kept as a compatibility path.

### Method 1 (Recommended): YAML Configuration File

#### File Location

The default path is `config/builtin_models.yaml` (in the same directory as `config.yaml` and `builtin_agents.yaml`). To mount it elsewhere, set the environment variable `BUILTIN_MODELS_CONFIG=/absolute/path/builtin_models.yaml` to override.

If the file does not exist at startup, it is silently skipped; parsing failures only log a Warning and do not affect the main process. On each application restart, the file is re-read and UPSERTed into the `models` table by the `id` field (preserving `created_at`, refreshing all other fields).

#### Schema

```yaml
builtin_models:
  - id: <required, stable, UPSERT key>
    tenant_id: <int, default 10000> # Aligned with tenants_id_seq start
    name: <string>
    type: KnowledgeQA | Embedding | Rerank | VLLM | ASR
    source: <string, default "remote"> # local | remote | aliyun | ...
    description: <string, optional>
    is_default: <bool, default false>
    status: <string, default "active">
    parameters:
      base_url: <string>
      api_key: <string, supports ${ENV_VAR}>
      provider: <string> # openai | generic | moonshot | ...
      embedding_parameters: # Embedding type only
        dimension: <int>
        truncate_prompt_tokens: <int>
```

#### Complete Example

```yaml
builtin_models:
  - id: builtin-openai-chat
    name: gpt-4o-mini
    type: KnowledgeQA
    source: remote
    description: OpenAI default chat model
    is_default: false
    parameters:
      base_url: https://api.openai.com/v1
      api_key: ${OPENAI_API_KEY}
      provider: openai

  - id: builtin-openai-embeddings
    name: text-embedding-3-small
    type: Embedding
    source: remote
    parameters:
      base_url: https://api.openai.com/v1
      api_key: ${OPENAI_API_KEY}
      provider: openai
      embedding_parameters:
        dimension: 1536
        truncate_prompt_tokens: 0

  - id: builtin-rerank
    name: bge-reranker-v2-m3
    type: Rerank
    source: remote
    parameters:
      base_url: ${RERANK_BASE_URL}
      api_key: ${RERANK_API_KEY}
      provider: generic
```

#### `${ENV}` Interpolation

Any **string** field (`api_key` / `base_url` / `name`, etc.) can reference environment variables: `${OPENAI_API_KEY}` is replaced at startup by the corresponding `os.Getenv("OPENAI_API_KEY")` value.

- Environment variable exists → replaced with its actual value
- Environment variable does not exist → **the literal `${OPENAI_API_KEY}` string is kept** (so a 401 error directly shows that the env variable was not set, making it easier to debug)
- Shell expansion like `${VAR:-default}` is **not supported**; behavior is consistent with the existing `config.yaml` interpolation implementation
- **Non-string fields cannot use env interpolation** (e.g. `type`, `dimension`, `is_default`), as they must be parsed according to the YAML target type

#### How Env Variables Reach the Container

The `app` service in `docker-compose.yml` already includes:

```yaml
env_file:
  - path: .env
    required: false
```

This means you write the variable values into the `.env` file at the project root, and they are automatically forwarded to the container at startup. There is **no need** to individually pass them through in the `environment:` block. `required: false` ensures the container can still start if `.env` does not exist (adapting to upstream fresh clone scenarios).

The repository's `.env.example` has a **Built-in Models** comment section at the top listing reference variable names for LLM / Embedding / Rerank as a starting point. Copy `.env.example` to `.env`, uncomment the relevant lines, and fill in your values. Variable names are determined by your YAML; the reference section is merely common boilerplate, not reserved words.

Complete end-to-end example:

`.env`

```bash
LLM_MODEL_NAME=gpt-4o-mini
LLM_BASE_URL=https://api.openai.com/v1
LLM_API_KEY=sk-...
LLM_PROVIDER=openai
```

`config/builtin_models.yaml`

```yaml
builtin_models:
  - id: builtin-llm-default
    type: KnowledgeQA
    is_default: true
    name: ${LLM_MODEL_NAME}
    parameters:
      base_url: ${LLM_BASE_URL}
      api_key: ${LLM_API_KEY}
      provider: ${LLM_PROVIDER}
```

Start:

```bash
docker compose up -d
```

#### Post-Startup Verification

```bash
docker compose logs app | grep -E 'Built-in models? config'
```

You will see something like:

```
Built-in model upserted: id=builtin-openai-chat name=gpt-4o-mini type=KnowledgeQA
Built-in model upserted: id=builtin-openai-embeddings name=text-embedding-3-small type=Embedding
Built-in models config applied: 2 entries from /app/config/builtin_models.yaml.
```

#### Docker Deployment

Mount the file in the `volumes` section of the `app` service in `docker-compose.yml`:

```yaml
services:
  app:
    volumes:
      - ./config/builtin_models.yaml:/app/config/builtin_models.yaml:ro
```

The repository provides `config/builtin_models.yaml.example` as a starting point. Copy it to `config/builtin_models.yaml` and modify as needed.

### Method 2: Direct SQL Insertion

Supported providers: `generic` (custom), `openai`, `aliyun`, `zhipu`, `volcengine`, `hunyuan`, `deepseek`, `minimax`, `mimo`, `siliconflow`, `jina`, `openrouter`, `gemini`, `modelscope`, `moonshot`, `qianfan`, `qiniu`, `longcat`, `gpustack`

```sql
-- Example: LLM built-in model
INSERT INTO models (
    id, tenant_id, name, type, source, description,
    parameters, is_default, status, is_builtin
) VALUES (
    'builtin-llm-001',
    10000,
    'gpt-4o-mini',
    'KnowledgeQA',
    'remote',
    'System built-in LLM model',
    '{"base_url": "https://api.openai.com/v1", "api_key": "sk-xxx", "provider": "openai"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- Embedding
INSERT INTO models (
    id, tenant_id, name, type, source, description,
    parameters, is_default, status, is_builtin
) VALUES (
    'builtin-embedding-001',
    10000,
    'text-embedding-3-small',
    'Embedding',
    'remote',
    'System built-in Embedding model',
    '{"base_url": "https://api.openai.com/v1", "api_key": "sk-xxx", "provider": "openai", "embedding_parameters": {"dimension": 1536, "truncate_prompt_tokens": 0}}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- Rerank
INSERT INTO models (
    id, tenant_id, name, type, source, description,
    parameters, is_default, status, is_builtin
) VALUES (
    'builtin-rerank-001',
    10000,
    'bge-reranker-v2-m3',
    'Rerank',
    'remote',
    'System built-in Rerank model',
    '{"base_url": "https://api.jina.ai/v1", "api_key": "jina-xxx", "provider": "jina"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;
```

### Verify Insertion Results

```sql
SELECT id, name, type, is_builtin, status
FROM models
WHERE is_builtin = true
ORDER BY type, created_at;
```

## Promoting an Existing Model to Built-in

If you have manually created a regular model and want to upgrade it to a built-in model:

```sql
UPDATE models
SET is_builtin = true
WHERE id = 'MODEL_ID';
```

## Removing Built-in Models

**Simply delete the entry from the YAML file.** On application restart, the loader automatically soft-deletes rows in the `models` table that were YAML-managed but are no longer declared in the YAML — you no longer need to run SQL manually.

How it works: Every row written by YAML is marked with `managed_by = 'yaml'`. On restart, the loader performs two steps:

1. UPSERT all entries from the current YAML (idempotent by `id`, including resetting `deleted_at` to NULL for previously soft-deleted rows — meaning removing an entry from YAML and adding it back effectively "resurrects" it)
2. Soft-delete rows where `is_builtin = true AND managed_by = 'yaml' AND id NOT IN (the set of ids in the current YAML)`

**Built-in rows inserted manually via SQL (`managed_by = ''`) are never touched by the loader**, remaining fully isolated from YAML management.

### Manual Path Supplement

If you manage built-in models via SQL (`managed_by = ''`), follow the old method for deletion:

```sql
-- Remove built-in flag, revert to a regular model
UPDATE models SET is_builtin = false WHERE id = 'MODEL_ID';

-- Or delete directly
DELETE FROM models WHERE id = 'MODEL_ID';
```

### Emergency Disable of YAML Management

If you accidentally modified the YAML and want to immediately disable YAML management without clearing the file, the fastest approach is: point the `BUILTIN_MODELS_CONFIG` environment variable to a non-existent path and restart — the loader will no-op when it sees the file is missing, **including skipping the drift sweep**, and already-written YAML-managed rows will remain intact.

## Important Notes

1. **ID naming convention**: Use the format `builtin-{type}-{slug}`, e.g. `builtin-openai-chat`, `builtin-rerank`
2. **Tenant ID**: Built-in models can belong to any tenant; the default is `10000` (aligned with the `tenants_id_seq` start value)
3. **YAML and SQL coexistence**: Both methods can be used simultaneously; the loader only touches rows with `managed_by='yaml'`; built-in rows inserted via SQL are completely invisible to the loader
4. **`is_default` uniqueness guarantee**: When an entry is marked `is_default: true` in YAML, the loader first sets all other default models under the same `(tenant_id, type)` to `false`, preventing the API-expected "one default model per type" semantics from being violated
5. **Restart to apply**: After modifying the YAML, run `docker compose restart app` for the new configuration to take effect
6. **Encryption**: API Keys are stored in encrypted form inside the `parameters` JSONB (if `SYSTEM_AES_KEY` is configured); otherwise, they fall back to plaintext for compatibility
7. **Security**: The frontend automatically hides the API Key and Base URL of built-in models, but the raw data still exists in the database — ensure proper database access controls
8. **Parse error self-protection**: If YAML parsing fails, the loader only logs a warning and skips reconciliation — it will **not** execute a drift sweep, ensuring a typo in the YAML does not mass-soft-delete existing built-in models
