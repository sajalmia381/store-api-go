# Docker container Deployment

> Steps to deploy a docker container in ec2 virtual machine.

## 1. Connect or transfer ec2 machine

### - Connect EC2 Machine
```bash
# ARM VPC
sudo ssh -i ~/Desktop/pem/storeApi.pem ubuntu@ec2-52-66-223-120.ap-south-1.compute.amazonaws.com
# NEW ARM VPC
sudo ssh -i ~/Desktop/pem/storeApi.pem ubuntu@ec2-3-7-68-106.ap-south-1.compute.amazonaws.com
```

### - Transfert Data To EC2 Machine
```bash
# Create folder and give permission
sudo mkdir /opt/front-end
sudo chown ubuntu:ubuntu /opt/front-end
sudo chown ubuntu:ubuntu /usr/share/nginx/html


# Transfer data to machine 
sudo scp -i <path-to-key-file> -r <path-to-local-dist-folder>/* ubuntu@<domain name>:/opt/front-end
sudo scp -i ~/Desktop/pem/storeApi.pem -r ./dump/* ubuntu@ec2-3-7-68-106.ap-south-1.compute.amazonaws.com:/opt/local-data
sudo scp -i ~/Desktop/pem/storeApi.pem -r ./dist/store-admin/* ubuntu@ec2-3-7-68-106.ap-south-1.compute.amazonaws.com
:/usr/share/nginx/html
```

## 2. Install Docker 
> Docker installing instruction main source [https://docs.docker.com/engine/install/ubuntu](https://docs.docker.com/engine/install/ubuntu)

### - Uninstall old version
```bash
sudo apt-get remove docker docker-engine docker.io containerd runc
```
### A. Setup docker repository
```bash
sudo apt-get update

sudo apt-get install \
  ca-certificates \
  curl \
  gnupg \
  lsb-release
```
### B. Add Docker’s official GPG key
```bash
sudo mkdir -p /etc/apt/keyrings

curl -fsSL https://download.docker.com/linux/ubuntu/gpg |sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
```
### C. Use the following command to set up the repository
```bash
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```
### D. Update the apt package index, and install the latest version of Docker Engine, containerd, and Docker Compose, or go to the next step to install a specific version:
```bash
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

### E. Verify Docker
```bash
sudo service docker start
sudo docker run hello-world
```

## 3. Install docker-compose
> Docker Compose installing instruction main source [https://docs.docker.com/compose/install/linux](https://docs.docker.com/compose/install/linux)
```bash
 sudo apt-get update
 sudo apt-get install docker-compose-plugin
```


## 4. Setup ufw firewall

```
sudo ufw enable
sudo ufw status
sudo ufw allow ssh (Port 22)
sudo ufw allow http (Port 80)
sudo ufw allow https (Port 443)
```

## 5. Install NGINX and configure

### - Install nginx *
```bash
# Install nginx 
sudo apt install nginx

# Check NGINX config
sudo nginx -t

# Restart NGINX
sudo service nginx restart
```

### - Deploy store-client in <USERNAME>
```bash
sudo nano /etc/nginx/sites-available/default
```
```
server {
  listen 80; # managed by Certbot
  server_name storerestapi.com;
  sendfile on;
  default_type application/octet-stream;
  gzip on;
  gzip_http_version 1.1;

  gzip_disable      "MSIE [1-6]\.";

  gzip_min_length   256;

  gzip_vary         on;

  gzip_proxied      expired no-cache no-store private auth;

  gzip_types        text/plain text/css application/json application/javascript application/x-javascript text/xml application/xml application/xml+rss text/javascript;

  gzip_comp_level   9;

  <USERNAME> /usr/share/nginx/html;

  location / {
    try_files $uri $uri/ /index.html =404;
  }
  location /health {
    return 200 'I am live :)';
  }
}
```

### - Deploy node-api
```bash
sudo nano /etc/nginx/sites-available/storeApiNode
sudo ln -s /etc/nginx/sites-available/storeApiNode /etc/nginx/sites-enabled/storeApiNode
```
```
server {
  charset utf-8;
  listen 80 default_server;
	listen [::]:80 default_server;

  server_name api.storerestapi.com;

  # Node api reverse proxy
  location / {
      proxy_pass http://localhost:8000;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection 'upgrade';
      proxy_set_header Host $host;
      proxy_cache_bypass $http_upgrade;
  }

  // Check health
  location /health {
    return 200 'I am live :)';
  }
}
```


## 6. Add SSL with LetsEncrypt

```
sudo add-apt-repository ppa:certbot/certbot
sudo apt-get update
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d storerestapi.com -d api.storerestapi.com

# Only valid for 90 days, test the renewal process with
certbot renew --dry-run
```

```bash
sudo apt install certbot python3-certbot-nginx
```

Now visit https://yourdomain.com and you should see your Node app

## -
# Use full command
```bash
mongorestore -d store-api /dump/store-api
```

<!-- After Deploy -->
```bash
sudo docker-compose -f docker-compose.prod.yml up --force-recreate --no-deps certbot
```