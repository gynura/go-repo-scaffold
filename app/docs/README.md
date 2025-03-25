## Docs

### Requirements

- [swagger-cli](https://www.npmjs.com/package/swagger-cli)
- [spectral](https://meta.stoplight.io/docs/spectral/)
- [swagger-ui-watcher](https://www.npmjs.com/package/swagger-ui-watcher)

### How to bundle the API spec

From **users-microservice/app** repository, run:
```bash
make bundle_api
```

### How to lint the API spec

From **users-microservice/app** repository, run:
```bash
make lint_api
```

### How to build and serve the documentation

From **users-microservice/app** repository, run:
```bash
make serve_api
```
