## Dump server DB

```shell
mongodump -db=store-api-go --uri=storerestapi.com --port=27017 --username=root --password=hello123 --authenticationDatabase=admin --out=/opt/backup/mongodump-2022-11-18

# NOt workable
mongodump -db=klovercloud-cost-management-security-server --host="mongodb://root:KuUQdSchulgdjZWQ@116.202.169.133:27035/?authSource=admin" --out=./backup/mongodump-202-10-24


```
## mongodb Restore to server db

```shell
mongorestore mongodb://root:hello123@3.7.68.106:27017 /opt/backup/mongodump-2022-03-24
```
