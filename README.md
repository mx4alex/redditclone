# redditclone
Клон реддита, написанный на Go

## Использование

### Сервис поддерживает следующие эндпоинты:
- `POST` `/api/register` - регистрация
- `POST` `/api/login` - логин
- `GET` `/api/posts/` - список всех постов
- `POST` `/api/posts/` - добавление поста
- `GET` `/api/posts/{CATEGORY_NAME}` - список постов конкретной категории
- `GET` `/api/post/{POST_ID}` - детали поста с комментариями
- `POST` `/api/post/{POST_ID}` - добавление комментария
- `DELETE` `/api/post/{POST_ID}/{COMMENT_ID}` - удаление комментария
- `GET` `/api/post/{POST_ID}/upvote` - рейтинг поста вверх
- `GET` `/api/post/{POST_ID}/downvote` - рейтинг поста вниз
- `GET` `/api/post/{POST_ID}/unvote` - отмена голоса
- `DELETE` `/api/post/{POST_ID}` - удаление поста
- `GET` `/api/user/{USER_LOGIN}` - получение всех постов конкретного пользователя
