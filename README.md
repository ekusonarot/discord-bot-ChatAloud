***---incomplete---***

***ChatAloud***

The TTS function is implemented using the docomo speech synthesis API.
Please apply from the link below and get a Token.
https://dev.smt.docomo.ne.jp/?p=docs.api.page&api_name=text_to_speech&p_name=api_reference

Create an app on Discord and get a Token.

**Usage**

Create a file in the following format in setting.json.

```
{
    "DISCORD_TOKEN": "{Your Discord Token}",
    "CLIENT_ID": "{Your Client ID}",
    "DOCOMO_API_TOKEN": "{Your docomo api Token}",
    "URL": "https://api.apigw.smt.docomo.ne.jp/futureVoiceCrayon/v1/textToSpeech"
}

```

Execute the following command

```
docker-compose up -d
```
