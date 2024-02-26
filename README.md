# go-http-streaming

Simple Go client for streaming data to an HTTP server via a POST multipart/form-data request

## Features

* Using I/O streams to avoid storing entire objects in memory
* Sending multiple files within one HTTP connection using the `SendUnclose()` method
* Zero-dependencies project

## Installation

```
go get github.com/kozld/go-http-streaming
```

## Example

The example below initializes the client and sends content from a file

```
c := client.New(http.MethodPost, endpoint)

file, _ := os.Open(filepath)

resp, err := c.Send(file)
if err != nil {
    log.Fatalln(err)
}

io.Copy(os.Stdout, resp.Body)
```

The following example shows sending multiple files within a single *HTTP* connection.  
**Note:** we need to close the connection ourselves by calling the `c.Close()` method to return `http.Response` or an `error`

```
var err error

c := client.New(http.MethodPost, endpoint)

file1, _ := os.Open(filepath1)
file2, _ := os.Open(filepath2)

err = c.SendUnclose(file1)
if err != nil {
    log.Fatalln(err)
}

err = c.SendUnclose(file2)
if err != nil {
    log.Fatalln(err)
}

resp, err := c.Close()
if err != nil {
    log.Fatalln(err)
}

io.Copy(os.Stdout, resp.Body)
```


# Комментарии к тестовому заданию

### Requirement

* Установленные docker и docker-compose
* Установленный Go

## Как проверить что все работает?

1. Вначале необходимо запустить тестовый web-сервер, на который будет осуществляться передача файлов через данный клиент. Запуск осуществляется следующей командой:
    ```
    go run server/web.go
    ```

    Сервер запустится на порту `3000`.

2. Внутри директории example/ находится файл main.go с тестовым сценарием по загрузке файлов на сервер.
    Он служит примером практического использования модуля `github.com/kozld/go-http-streaming`.

    Необходимо, чтобы до запуска данного сценария web-сервер был запущен.
    По-умолчанию сценарий ожидает web-сервер по адресу `http://localhost:3000`. Это значение можно переопределить, задав переменную окружения `UPLOAD_HOST`.

    Тестовый сценарий при запуске ожидает переданный флаг `-filepath`, содержащий путь до загружаемого файла.  
    Таким образом, запустить его можно так:

    ```
    go run example/main.go -filepath путь_к_загружаемому_файлу
    ```

    Внутри директории payload/ содержится файл book.txt, его можно использовать в целях тестирования сценария.


3. Чтобы проверить, что размер передаваемых данных может превышать размер доступной для нашего процесса оперативной памяти, предлагается запуск тестового сценария      внутри docker-контейнера с установленным лимитом памяти.
    
    Лимит памяти задается внутри файла docker-compose.yml в корне проекта. Сейчас там установлено ограничение в 20Mb (можно менять).
    
    Для комфортного запуска всего окружения внутри Makefile подготовлена команда `docker`. Она собирает образ на основе Dockerfile и запускает контейнер с тестовым сценарием. 
    
    Файл с полезной нагрузкой передается внутрь контейнера через docker-volume указанный внутри файла docker-compose.yml в следующей секции:
    
    ``` 
    ...
    
    volumes:
          - ${PWD}/payload/:/tmp/payload
    ...      
    ```
    
    Таким образом, внутри нашего контейнера, как минимум уже будет находиться файл `/tmp/paylod/book.txt`.
    
    Изменить значение флага `-filepath` можно также внутри файла docker-compose.yml в следующей секции:
    ```
    ...
    command:
          - "-filepath"
          - "/tmp/payload/10GB.bin"
    ...
    ```
    
    Примечание: тестовый файл payload/book.txt слишком мал для качественной проверки, рекомендуется пробросить внутрь контейнера что-то более тяжелое, я тестировал с файлом 10GB, в случае чего, могу им поделиться. =)
