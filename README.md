[![GitHub Workflow][1]](https://github.com/missuo/openai-translate/actions)
[![Go Version][2]](https://github.com/missuo/openai-translate/blob/main/go.mod)
[![Go Report][3]](https://goreportcard.com/badge/github.com/missuo/openai-translate)
[![Maintainability][4]](https://codeclimate.com/github/missuo/openai-translate/maintainability)
[![GitHub License][5]](https://github.com/missuo/openai-translate/blob/main/LICENSE)
[![Docker Pulls][6]](https://hub.docker.com/r/missuo/deeplx)
[![Releases][7]](https://github.com/missuo/openai-translate/releases)

[1]: https://img.shields.io/github/actions/workflow/status/missuo/openai-translate/build.yml?logo=github
[2]: https://img.shields.io/github/go-mod/go-version/missuo/openai-translate?logo=go
[3]: https://goreportcard.com/badge/github.com/missuo/openai-translate
[4]: https://api.codeclimate.com/v1/badges/b5b30239174fc6603aca/maintainability
[5]: https://img.shields.io/github/license/missuo/openai-translate
[6]: https://img.shields.io/docker/pulls/missuo/openai-translate?logo=docker
[7]: https://img.shields.io/github/v/release/missuo/openai-translate?logo=smartthings

## Features
- Setting up an API Key once allows it to be called from anywhere.
- Unlimited requests, billed by usage.
- Deployed on overseas servers, it can bypass the Great Firewall (GFW).

## Usage
### Request Parameters
- text: string
- source_lang: string
- target_lang: string

### Response
```json
{
  "code": 200,
  "data": "Hello, Britain!",
  "source_lang": "ZH",
  "target_lang": "EN"
}
```
### Docker Compose
```bash
mkdir openai-translate && cd openai-translate
wget https://raw.githubusercontent.com/missuo/openai-translate/main/compose.yaml
nano compose.yaml # Modify to YOUR_API_KEY
docker compose up -d
```

### Docker
```bash
docker run -itd -p 23333:23333 -e OPENAI_KEY=YOUR_API_KEY ghcr.io/missuo/openai-translate:latest
```
### Setup on [Bob App](https://bobtranslate.com/)
> [!IMPORTANT]  
> **This project is fully compatible with the bob-plugin-deeplx plugin.**
1. Install [bob-plugin-deeplx](https://github.com/missuo/bob-plugin-deeplx) on Bob.

2. Setup the API. (If you use Brew to install locally you can skip this step)
![c5c19dd89df6fae1a256d](https://missuo.ru/file/c5c19dd89df6fae1a256d.png)

### Use in Python
```python
import httpx, json

deeplx_api = "http://127.0.0.1:23333/translate"

data = {
	"text": "Hello World",
	"source_lang": "EN",
	"target_lang": "ZH"
}

post_data = json.dumps(data)
r = httpx.post(url = deeplx_api, data = post_data).text
print(r)
```

