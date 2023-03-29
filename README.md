# go-news-aggregator
Это сервис осуществляет парсинг RSS-лент новостных сайтов в формате `.xml`. 

Он работает в составе приложения со следующими микросервисами:
- [news-gateway]() - входная точка приложения
- [news-comments-service]() - сервис комментариев
- [news-formatter-service]() - сервис проверки комментариев

<br />

---

## API 

* GET `/news/{id}` - возвращает определенную новость.
* GET `/news` - возвращает новость с указанными параметрами запроса `filter` и `page`

---

## RSS парсер
Перечень сайтов указывается в `config.json` <br>
```json
{
    "rss":[
       "https://rssexport.rbc.ru/rbcnews/news/30/full.rss",
       "https://habr.com/ru/rss/hub/go/all/?fl=ru"
    ],
    "request_period": 5
}
```