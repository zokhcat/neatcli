# neatcli

CLI for [Neatlogs](https://neatlogs.com) which does prompt versioning, tool description
management, and trace inspection.

Prompt management endpoints were extracted from the
[Python SDK](https://github.com/neatlogs/neatlogs/blob/4f5daf2270a8ecc208e81b1cbb8cab8cc0412eb6/neatlogs/prompt/client.py#L307).

## Getting started

```
NEATLOGS_API_KEY=nlw_xxx neatcli init --project my-project
```

Creates `~/.neatlogs/config.yaml` and a `.neatlogs/` workspace in the current
directory with `prompts/` and `tools/` subdirectories.

## Build

```
make build
make install
make clean
make run
```

## Prompt management

Works with the Neatlogs management API. Every prompt has a version history and
labels (`production`, `staging`, etc.).

| Command | Description |
|---|---|
| `neatcli prompt list` | List prompts and versions. Filter by `--name` or `--label`. |
| `neatcli prompt get <name>` | Print prompt content and metadata. Optional `--label` or `--version`. |
| `neatcli prompt create <name>` | Create a new prompt. Use `--content`, `--file`, or `--type chat`. |
| `neatcli prompt pull <name>` | Download a prompt to `.neatlogs/prompts/<name>.yaml`. |
| `neatcli prompt push <name>` | Upload local YAML as a new version. |
| `neatcli prompt diff <name> <a> <b>` | Unified diff between two versions (numbers or labels). |
| `neatcli prompt promote <name> <version>` | Move a label to a version (`--label production`). |
| `neatcli prompt rollback <name>` | Revert production label to the previous version. Needs `--force`. |
| `neatcli prompt log <name>` | Version history table. |

```
neatcli prompt list
neatcli prompt list --label production

neatcli prompt get support-prompt --label production
neatcli prompt get support-prompt --version 3

neatcli prompt create support-prompt --file prompt.txt --labels staging
neatcli prompt create greeting-prompt --type chat --labels production

neatcli prompt pull support-prompt
neatcli prompt pull support-prompt --version 2

neatcli prompt push support-prompt --message "Updated tone"

neatcli prompt diff support-prompt 2 3
neatcli prompt diff support-prompt production staging

neatcli prompt promote support-prompt 3 --label production

neatcli prompt rollback support-prompt --force

neatcli prompt log support-prompt
```

## Tool descriptions

Operates on local YAML files in `.neatlogs/tools/`. The server-side tool
management API is not publicly documented, so these are local-only for now.

| Command | Description |
|---|---|
| `neatcli tool list` | List tools from workspace and server. |
| `neatcli tool describe <name>` | Show a tool's description and JSON schema. |
| `neatcli tool update <name>` | Create or update a tool description locally. |

```
neatcli tool list

neatcli tool describe get_account

neatcli tool update get_account --desc "Fetch account by customer ID"
neatcli tool update get_account --schema '{"type":"object","properties":{"id":{"type":"string"}}}'
```

## Traces

Placeholder commands. Traces are viewable in the dashboard at
https://app.neatlogs.com. Neatlogs does not currently expose a read-back
API for traces.

| Command | Description |
|---|---|
| `neatcli trace list` | List recent traces. |
| `neatcli trace get <id>` | View a trace's full span tree. |

## Implementation notes

Prompt commands hit the Neatlogs API directly. The endpoints were mapped from
the Python SDK linked above and should work against any Neatlogs instance.

Tool commands are local-only. The YAML files in `.neatlogs/tools/` are ready to
sync if Neatlogs publishes a tool management API.

Trace commands are stubs. If Neatlogs adds trace querying to their API, the CLI
is ready to wire it in.

Config lives in `~/.neatlogs/config.yaml`. The API key can also be set via
`NEATLOGS_API_KEY`.
