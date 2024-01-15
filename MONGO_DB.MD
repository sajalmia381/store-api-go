## Dump server DB

```shell
mongodump -db=store-api-go --uri=storerestapi.com --port=27017 --username=<USERNAME> --password=<PASSWORD> --authenticationDatabase=admin --out=/opt/backup/mongodump-2022-11-18

# NOt workable
mongodump -db=<K8s-DEPLOYMENT-SERVICE> --host="mongodb://<USERNAME>:<PASSWORD>@<SERVER-IP>:<DB-EXPOSED-PORT>/?authSource=admin" --out=./backup/mongodump-202-10-24


```
## mongodb Restore to server db

```shell
mongorestore mongodb://<USERNAME>:<PASSWORD>@3.7.68.106:27017 /opt/backup/mongodump-2022-03-24
```

## Access docker mongo container
```shell
docker exec -it store-api bash/sh
```
## Access mongosh from docker
```shell
mongosh -u <USERNAME> -p hello123
```