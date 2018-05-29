# Summary
OpenSDS dashboard uses the front-end development framework Angular5 (https://angular.io/)
and relies on PrimeNG UI Components (https://www.primefaces.org/primeng/). Regardless of 
deployment or two development, prepare the corresponding environment.

# Prerequisite 

### 1. Ubuntu
* Version information
```shell
root@proxy:~# cat /etc/issue
Ubuntu 16.04.2 LTS \n \l
```

### 2. Nginx installation
```shell
sudo apt-get install nginx
```

### 3. NodeJS installation, NPM will be installed with nodejs.
```shell
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### 4. Angular CLI installation
Specify the version[1.7.4] of angular5 suitable for installation.
```shell
sudo npm install -g @angular/cli@1.7.4
```


# Build & Start
### 1. Git clone dashboard code.
```shell
git clone https://github.com/opensds/opensds.git
```

### 2. Build opensds dashboard.
After the build work finished, the files in the `dist` folder should be copied to the folder ` /var/www/html/`.
```shell
cd opensds/dashboard
sudo npm install
sudo ng build --prod
```

```shell
cp -R opensds/dashboard/dist/* /var/www/html/
```

### 3. Set nginx default config.
```shell
vi /etc/nginx/sites-available/default 
```
Configure proxy, points to the resource server and the authentication server respectively.
Such as: 
* Keystone server `http://1.1.1.0:5000`
* Resource server `http://1.1.1.0:50040`
```shell
location /v3/ {
    proxy_pass http://1.1.1.0:5000/v3/;
}

location /v1beta/ {
    proxy_pass http://1.1.1.0:50040/v1beta/;
}
```

### 4. Restart nginx
```shell
service nginx restart 
```

### 5. Access dashboard in browser.
```shell
http://localhost/
```
