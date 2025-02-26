# pipe-node

Pipe for installing node.js dependencies and building node.js applications on CI/CD.

`pipe-node [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

## Commands

### `login`

Login to the given NPM registries.

`pipe-node login [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Login

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | npm registries to login to. | `String`<br/>`json([]struct{ username: string, password: string, registry?: string, useHttps?: boolean })` | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | ".npmrc" |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |

##### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String`<br/>`enum("npm", "yarn", "pnpm")` | `false` | yarn |

### `install`

Install node.js dependencies with the given package manager.

`pipe-node install [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Install

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_INSTALL_CWD` | Install CWD for the package manager. | `String` | `false` | . |
| `$NODE_INSTALL_USE_LOCK_FILE` | Use the lockfile while installing the packages. | `Bool` | `false` | false |
| `$NODE_INSTALL_ARGS` | Arguments to append to install command. | `String` | `false` |  |
| `$NODE_CACHE_ENABLE` | Enable caching for the package manager. | `Bool` | `false` | false |

##### Login

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | npm registries to login to. | `String`<br/>`json([]struct{ username: string, password: string, registry?: string, useHttps?: boolean })` | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | ".npmrc" |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |

##### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String`<br/>`enum("npm", "yarn", "pnpm")` | `false` | yarn |

### `build`

`pipe-node build [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Build

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_BUILD_SCRIPT` | package.json script for building operation. | `String`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` | `false` | build |
| `$NODE_BUILD_SCRIPT_ARGS` | package.json script arguments for building operation. | `String`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` | `false` |  |
| `$NODE_BUILD_CWD` | Working directory for build operation. | `String` | `false` | . |

##### Environment

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `Bool` | `false` | false |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `String`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | [<br />  { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },<br />  { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },<br />  { "match" :"^heads/main$", "environment": "develop" },<br />  { "match": "^heads/master$", "environment": "develop" }<br />] |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `Bool` | `false` | false |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `Bool` | `false` | false |

##### GIT

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |

##### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String`<br/>`enum("npm", "yarn", "pnpm")` | `false` | yarn |

### `run`

`pipe-node run [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Command

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_COMMAND_SCRIPT` | package.json script for given command operation. | `String`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` | `false` |  |
| `$NODE_COMMAND_CWD` | Working directory for the given command operation. | `String` | `false` | . |

##### Environment

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `Bool` | `false` | false |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `String`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | [<br />  { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },<br />  { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },<br />  { "match" :"^heads/main$", "environment": "develop" },<br />  { "match": "^heads/master$", "environment": "develop" }<br />] |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `Bool` | `false` | false |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `Bool` | `false` | false |

##### GIT

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |

##### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String`<br/>`enum("npm", "yarn", "pnpm")` | `false` | yarn |
