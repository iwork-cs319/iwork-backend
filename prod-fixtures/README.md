## Prod-Fixtures
- Run tables.sql  in ../resources/ to clean the DB
- Run following command to upload a base floor plan
```bash
curl --request POST \
  --url https://icbc-go-api.herokuapp.com/floors \
  --header 'content-type: multipart/form-data; \
  --form 'name=West 2nd Avenue' \
  --form image=@Location1.jpg
```
- Use the server to upload users.csv file
``` bash
curl --request POST \
  --url https://icbc-go-api.herokuapp.com/users \
  --header 'content-type: multipart/form-data; \
  --form users=@users.csv
```
- Use the server to upload workspaces.csv file
```bash
curl --request POST \
  --url https://icbc-go-api.herokuapp.com/assignments \
  --header 'content-type: multipart/form-data; \
  --form assignments=@workspaces.csv
```
- Run book_offer.sql in db console