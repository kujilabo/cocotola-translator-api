# cocotola-translator-api

```sh
curl -H "Authorization: Basic dXNlcjpwYXNzd29yZA==" localhost:8080/v1/user/dictionary/lookup
```

```sh
grpcurl -plaintext -proto ./proto/translator_user.proto -d '{"fromLang2":"en", "toLang2": "ja", "text": "book"}' -H "Authorization: basic dXNlcjpwYXNzd29yZA==" localhost:50151 proto.TranslatorUser/DictionaryLookup 
```
