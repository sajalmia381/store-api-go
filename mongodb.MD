## Dump server DB

```shell
mongodump --host=storerestapi.com --port=27017 --username=<USERNAME> --authenticationDatabase=admin --out=/opt/backup/mongodump-2022-03-24

# NOt workable
mongodump -db=<K8s-DEPLOYMENT-SERVICE> --host="mongodb://<USERNAME>:<PASSWORD>@<SERVER-IP>:<DB-EXPOSED-PORT>/?authSource=admin" --out=./backup/mongodump-202-10-24


```
## mongodb Restore to server db

```shell
mongorestore mongodb://<USERNAME>:hello123@3.7.68.106:27017 /opt/backup/mongodump-2022-03-24
```
