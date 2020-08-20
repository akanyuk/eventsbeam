# protoc-gen-web

Генерация web-хэндлеров на основе protobuf спецификации для использования совместно с веб-сервером.  

## Required: 

- [protoc](https://github.com/google/protobuf)
- [protoc-gen-go](https://github.com/golang/protobuf)
- [protoc-gen-micro](https://github.com/micro/micro/tree/master/cmd/protoc-gen-micro)

## Опции хэндлера

* Тип операции: GET, PUT, POST, DELETE, PATCH
* Имя и описание для OpenAPI

## Поддерживаемые параметры запроса

* `path` – `integer`
* `body` – `structure` 

Поддержка других типов запросов будет добавляться по мере необходимости.

## Типы генерируемых ответов

### По умолчанию

По умолчанию сгенерированный хэндлер возвращает ответ, указанный в параметрах `rpc` в формате `JSON`

### Хэндлер скачивания файла

Если в качестве ответа rpc-процедуры указать тип `web.FileDownloadResponse`, то будет сгенерирован хэндлер, возвращающий файл для скачивания. Поддерживается *Conditional Get* на основе даты модификации файла, которую возвращает ресивер.

## swagger-tricks

При генерации API формируются комментарии, совместимые с генератором [go-swagger](https://github.com/go-swagger). Для того, чтобы сформировать корректный swagger-файл, необходимо учитывать следующее:  

* Для всех структуры, используемых в параметрах rpc-вызовов необходим комментарий `swagger:model`
* Если в качестве параметров используется как вся структура, так и какое-либо поле этой структуры, то помечать поле комментарием `swagger:ignore`. Пример: хэндлер обновления объекта, принимающий идентификатор объекта в пути запроса. 

## Примеры использования

### Пример proto-файла

```proto
import "web.proto";

service Foo {
    rpc GetMessage (Dummy) returns (Dummy) {
        option (web.handler) = {
            get: "/sounds"
            description: "Список звуков"
        };
        option (web.handler).patameters = {
            path: "ID звукового файла"
            name: "id"
        };
    }
}

message Dummy {}
```

### Вызов генератора
```cmd
protoc --proto_path=. --micro_out=. --pbweb_out=. --go_out=. sounds.proto
```

### Подробнее

Более подробно с настройкой генерации можно ознакомиться на примере [sounds_api](https://git.sedmax.ru/BACKEND/sounds_api). 