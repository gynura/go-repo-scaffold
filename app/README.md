## App

### Requirements

- [swagger-cli](https://www.npmjs.com/package/swagger-cli)
- [spectral](https://meta.stoplight.io/docs/spectral/)
- [swagger-ui-watcher](https://www.npmjs.com/package/swagger-ui-watcher)

### How to maintain infra and API aligned

The idea is that everytime something changes inside the API Spec, we run from **users-microservice/app**: 
```bash
make 
```

Which will bundle, lint and generate the http infra go code. 

### How to bundle the API spec

From **users-microservice/app**, run:
```bash
make bundle_api
```

### How to lint the API spec

From **users-microservice/app**, run:
```bash
make lint_api
```

### How to build and serve the documentation

From **users-microservice/app**, run:
```bash
make serve_api
```

### How to generate the http infra 

From **users-microservice/app**, run:
```bash
make generate_http_infra
```

> **_NOTE:_**  This command will install the ogen generator for you

### How to generate mocks

From **users-microservice/app**, run:
```bash
make generate_mocks
```

> **_NOTE:_**  This command will install the mockgen generator for you
