## How to backup database
`heroku pg:backups:capture --app icbc-go-api`
This will create a backup and return an id
```bash
Backing up DATABASE to b001... done
```

## List backups
`heroku pg:backups --app sushi`

## Restore database
Use id generated from before
`heroku pg:backups:restore b001 DATABASE_URL --app icbc-go-api`
