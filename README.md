# xk6-beanstalkd

Beanstalkd Extension for K6, allows you to write k6 tests against your beanstalk tubes and more.

It makes use of [xk6](https://github.com/grafana/xk6) providing a client for interacting with Beanstalkd, a simple and fast work queue.

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go](https://go.dev/) 1.22.6 or later
- Git

Then:

1. Install `xk6`:
  ```bash
  go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```bash
  xk6 build --with github.com/dnlowman/xk6-beanstalkd=.
  ```

## Usage

Run the binary with your script:

```bash
./k6 run script.js
```

## API

This extension implements the following API:

### `newClient(address)`

Creates a new Beanstalkd client.

- `address`: string, the address of the Beanstalkd server (e.g., "localhost:11300")

### `client.put(data, priority, delay, ttr)`

Puts a job into the currently used tube.

- `data`: string, the job body
- `priority`: number, job priority
- `delay`: number, delay in seconds before the job becomes ready
- `ttr`: number, time to run in seconds

Returns: job ID (number)

### `client.reserve(timeout)`

Reserves a job from a watched tube.

- `timeout`: number, timeout in seconds

Returns: [job ID (number), job body (string)]

### `client.delete(id)`

Deletes a job.

- `id`: number, job ID

### `client.release(id, priority, delay)`

Releases a reserved job back to the ready queue.

- `id`: number, job ID
- `priority`: number, new priority
- `delay`: number, delay in seconds

### `client.bury(id, priority)`

Buries a job.

- `id`: number, job ID
- `priority`: number, new priority

### `client.kick(bound)`

Kicks buried or delayed jobs into the ready queue.

- `bound`: number, upper bound on the number of jobs to kick

Returns: number of jobs actually kicked

### `client.use(tube)`

Changes the tube used for `put` commands.

- `tube`: string, tube name

### `client.watch(tube)`

Adds a tube to the watch list for `reserve` commands.

- `tube`: string, tube name

### `client.ignore(tube)`

Removes a tube from the watch list for `reserve` commands.

- `tube`: string, tube name

### `client.peek(id)`

Peeks at a job.

- `id`: number, job ID

Returns: job body (string)

### `client.stats()`

Returns a map of server statistics.

### `client.statsTube(tube)`

Returns a map of tube statistics.

- `tube`: string, tube name

### `client.listTubes()`

Returns a list of all existing tubes.

## Example

```javascript
import beanstalkd from 'k6/x/beanstalkd';
import { check } from 'k6';

export default function () {
  const client = beanstalkd.newClient('localhost:11300');

  const jobId = client.put('Hello, Beanstalkd!', 1, 0, 60);
  console.log(`Put job with ID: ${jobId}`);

  const [reservedId, jobBody] = client.reserve(5);
  console.log(`Reserved job ${reservedId}: ${jobBody}`);

  check(jobBody, {
    'job content is correct': (body) => body === 'Hello, Beanstalkd!',
  });

  client.delete(reservedId);
  console.log('Job deleted');

  client.close();
}
```

## Notes

- Make sure you have a Beanstalkd server running and accessible at the address you specify when creating a new client.
- Remember to close the client when you're done to release resources.
- This extension is designed for testing purposes and may not include all features of a full Beanstalkd client library.
- This is still in BETA, but should work ok!