# Music Lirbrary API

- **Listing song data with filtering and pagination:**
    ```http
    GET /song
    ```
     ```http
    GET /song?group=beatles&page=1
    ```
    - queries for filtering:
        - group (group name)
        - name (song name)
    - queries for pagination:
        - page
        - page_size
    - sample output:
    ```json
    "songs": [
        {
            "id": 1,
            "song": "let it be",
            "group": "the beatles",
            "verses": [
                {
                    "id": 1,
                    "verse_number": 1,
                    "name": "When I find myself in times of trouble/n"
                },
                {
                    "id": 2,
                    "verse_number": 2,
                    "name": "Mother Mary comes to me/n"
                }
            ]
        }
    ]
}
- **Listing all verses for a certain song with filtering and pagination:**
- - required parameter: `id`
    ```http
    GET /song/{id}
    ```
     ```http
    GET /song/1?verseNumber=1&page=1
    ```
    - queries for filtering:
        - VerseNumber
    - queries for pagination:
        - page
        - page_size
    - sample output:
    ```json
    "verses": [
        {
            "id": 1,
            "verse_number": 1,
            "name": "When I find myself in times of trouble/n"
        }
    ]

    ```

- **Update song info:**
    - required parameter: `id`
     ```http
    PATCH /song/:id
    ```
    - input body:
    ```json
    {
    "group": "Muse",
    "song": "Supermassive Black Hole"
    }
    ```
- **Delete song :**
    - required parameter: `id`
     ```http
    DELETE /song/:id
    ```
- **Adding new song data**
    ```http
    POST /song
    ```
    - input body:
    ```json
    {
    "group": "Muse",
    "song": "Supermassive Black Hole"
    }
    ```
- **Второй пункт, где надо делать запрос в АПИ, не очень понятно было при добавлении чего надо делать запрос, сделал, что надо передавать в body group и name**
    ```http
    POST /verse
    ```
    - input body:
    ```json
    {
    "group": "Muse",
    "song": "Supermassive Black Hole"
    }
    ```
- **Swagger documentation**
    ```http
    GET /swagger/
    ```
---
### START PROJECT
- **.env file is in config dir and should be in it other way might not launch. 
- **install dependencies:**
```
go mod tidy
```
- **run project:**
```
go run ./cmd/api 
```
- **migrate down (chage -database parametr):**
```
migrate -path ./migrations -database 'postgres://postgres:123123@localhost:5432/postgres?sslmode=disable' down
```