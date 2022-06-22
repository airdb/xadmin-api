# xadmin-api

## Quickstart

```sh
git clone https://github.com/airdb/xadmin-api.git
brew install bufbuild/buf/buf
make plugins
make buf
make run
```

## Github Actions

### Environments

| Key                        | Describe                     | Default |
| -------------------------- | ---------------------------- | ------- |
| SCF_NAMESPACE              | serverless namespace         | default |
| SERVERLESS_PLATFORM_VENDOR | serverless provider name     | tencent |
| TENCENT_SECRET_ID          | tencent cloud secret id      |         |
| TENCENT_SECRET_KEY         | tencent cloud secret key     |         |
| XADMIN_CONFIG              | the config used to overwrite |         |
| XADMIN_JWT_KEY             | the key used by jtw          |         |


![Alt](https://repobeats.axiom.co/api/embed/3e87e8fc8ca29c0a9aaff199afb8a06e589fbdca.svg "Repobeats analytics image")
