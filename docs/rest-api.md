# eflows4HPC REST API to trigger workflows

## Abstract

This API is designed to be used by *end-users* of the workflows deployed using the eFlows4HPC HPCWaaS interface. It allows to:

* list available workflows
* trigger a workflow execution
* monitor a workflow execution
* cancel a workflow execution

The design of workflows themself is out of the scope of this API and is done by another user called *developer* in the Alien4Cloud application.

## API

### Authentication / Authorization

HPCWaaS handles authentication with the OAuth 2 protocol. Identities are managed by the Unity identity provider. Authentication in HPCWaaS is a two-ste process:

* Step 1: Retrieve an **access token** by visiting the `/auth/login` endpoint in your browser.
* Step 2: Use the token for sending requests to the API.

#### Authenticating with the REST API

For accessing the REST API with a general utility like `curl`, you need to pass the token in the header, e.g.  
`curl -H "Authorization: Bearer <access_token>" ...`  

#### Authenticating with the CLI utility

For the `waas` CLI utility, you can pass the token in three different places:

* In the WaaS config file with the `access_token` key, e.g.  
  `access_token: <access_token>`
* In the `HW_ACCESS_TOKEN` environment variable, e.g.  
  `export HW_ACCESS_TOKEN=<access_token>`
* In the command-line options, e.g.  
  `waas workflows list -t=<access_token>`  
  or  
  `waas workflows list --access_token=<access_token>` 

The parameters take precendence in the following order: command-line option > environment variable > config file.

### API Endpoints

#### Request authentication token

This API endpoint is to be used in a browser. It allows, after logging in to a Unity server, to retrieve an **access token** that is needed to authenticate when accessing other endpoints.

##### Endpoint

`/auth/login`

#### List available workflows

This API endpoint provides the workflows that could be triggered by an *end-user*

***TODO***: should probably handle pagination

##### Request

`GET /workflows`

##### Response

```
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "workflows": [
    {
      "id": "workflow-unique-id",
      "name": "workflow-name",
      "application_id": "app_id",
      "environment_id": "env_id",
      "environment_name": "env_name"
    },
    {
      "id": "workflow2-unique-id",
      "name": "workflow2-name",
      "application_id": "app2_id",
      "environment_id": "env_id",
      "environment_name": "env_name"
    }
  ]
}
```

#### Trigger a workflow execution

This API endpoint allows the *end-user* to trigger a workflow execution

##### Request

`POST /workflows/<workflow_name>`

```json
{
  "inputs": {
    "input1": "",
    "input2": ""
  }
}
```

##### Response

```
HTTP/1.1 201 Created
Location: /executions/<execution_id>
Content-Length: 0
```

#### Monitor a workflow execution

This API endpoint allows the *end-user* to monitor a workflow execution

##### Request

`GET /executions/<execution_id>`

##### Response

```
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "id": "<execution_id>",
  "status": "RUNNING/SUCCESS/FAILED/CANCELLING/CANCELLED",
  "outputs": {
    "output1": "",
    "output2": ""
  }
}
```

#### Cancel a workflow execution

This API endpoint allows the *end-user* to cancel a workflow execution


##### Request

`DELETE /executions/<execution_id>`

##### Response

```
HTTP/1.1 202 Accepted
Content-Length: 0
```

#### Get logs of a workflow execution

This API endpoint allows the *end-user* to get logs of an execution.

This endpoint implements a long polling mechanism, meaning that the response will be returned only if
new logs are founds or if a given timeout (see query parameters bellow) is reached.
This means that an empty array for logs could be returned in a `200 OK` response if no new logs are found and
the timeout is reached.

This endpoint supports pagination of results using the `from` and `size` parameters described below.

##### Request

`GET /executions/<execution_id>/logs?from=0&timeout=1m&size=-1&levels=0`

Query parameters:

* `from` (default `0`): Used for pagination, get logs from `from` index
* `size` (default `-1`): Used for pagination, get up to `size` log entries
* `timeout` (default `1m`): Long polling maximum duration
* `levels` (default `0` meaning `INFO` + `WARN` + `ERROR`): Bitmask enumeration with `DEBUG=1`, `INFO=2`, `WARN=4`, `ERROR=8` (Ex: `15` means all logs, `9` means `DEBUG` + `ERROR`, ...)

##### Response

```
HTTP/1.1 200 Accepted
Content-Type: application/json
```

```json
{
  "logs": [
    {
      "level": "INFO",
      "timestamp": "2022-10-20T12:32:29.774Z",
      "content": "Workflow \"exec_job\" ended successfully"
    },
    {
      "level": "INFO",
      "timestamp": "2022-10-20T12:32:29.797Z",
      "content": "Status for workflow \"exec_job\" changed to \"done\""
    }
  ],
  "total_results": 205,
  "from": 203
}
```

#### Create an SSH Key Pair

This API endpoint allows the *end-user* to create an SSH Key Pair and to optionally attach metadata to it.
The private key is stored in HashiCorp Vault and the public key is returned to the user along with a randomly
generated identifier for this key that should be used to retrieve the private key.
The public key and the identifier can not be seen again, it should be written down carefully and kept in safe place.

##### Request

`POST /ssh_keys`

```json
{
  "metadata": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

The request body is fully optional.
`privateKey` and `publicKey` are reserved keywords and should not be used as metadata keys.
If provided they will be silently ignored.

##### Response

```
HTTP/1.1 201 Created
Content-Type: application/json
```

```json
{
  "id": "<key_id>",
  "public_key": "<public_key>"
}
```
