# eflows4HPC REST API to trigger workflows

**Status:** Draft

## Abstract

This API is designed to be used by eflows4HPC *end-users*. It should allow to:

* list available workflows
* trigger a workflow execution
* monitor a workflow execution
* cancel a workflow execution

Design of workflows themself is out of the scope of this API and is done by another user called *developer* in the Alien4Cloud application.

### Open questions

* [ ] should we have an endpoint to list workflows executions? (Priority 1)
* [ ] should we have an endpoint to describe a workflow (typically expected inputs / outputs)? (Priority 1)
* [ ] should we have an endpoint to stream logs/events? (Events Priority 3 / Logs Priority 4)
* [ ] Authentication / Authorization (see below)

## API Design

### Authentication / Authorization

While this is identified as a mandatory feature, there is only an optional HTTP Basic authentication implemented in the first MVP
(Minimum Viable Product).

### API Endpoints

#### List available workflows

This API endpoint allows to list workflows that could be triggered by a *end-user*

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
    "workflow1",
    "workflow2"
  ]
}
```

#### Trigger a workflow execution

This API endpoint allows the *end-user* to trigger a workflow execution

***TODO***: should refine inputs

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

***TODO:*** Should work on the response model

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
