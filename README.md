# Bubble Files
## Назначение
Микросервис для сохранения файлов на определнный срок и на все время для проекта [Bubble](https://github.com/mongolianpes/Bubble-Messenger-Server)

## Функционал
- сохранения аватарок пользователей на все время, пока пользователь сам не удалит аватарку
- сохранения файлов, которыми обмениваются пользователя на определнное время

## RPC запросы
```
service FilesService {
  rpc Save(SaveFileRequest) returns (SaveFileResponse);
  rpc Del(DelFileRequest) returns (DelFileResponse);
}
```
Подробнее в файле files/proto/files.proto

## Порты
Для корректной работы необходимо открыть порт 8080 для HTTP запросов клиентов и отправки им файлов, а также 8086 для взаимодействия между сервисами
