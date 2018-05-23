## Summary
OpenSDS dashboard uses the front-end development framework Angular 2 and relies
on PrimeNG UI Components. Regardless of deployment or two development, prepare
the corresponding environment.

## Environment
* NodeJS
Download the latest version of NodeJS package from [NodeJS](https://nodejs.org/)
official website for installation.

* Angular CLI
```shell
npm install -g @angular/cli
```

## Deployment
* Install Angular "node_modules"
```shell
cd dashboard && npm install
```

* Build OpenSDS dashboard
```shell
ng build --prod
```
After the build work finished, the files in the `dist` folder should be copied to the root
directory of the web server.

## Links
* [Primeng](https://www.primefaces.org/primeng/)
* [Primeng Github](https://github.com/primefaces/primeng/)
* [NodeJS](https://nodejs.org/)
* [Angular](https://angular.io/)
* [Angular/cli](https://cli.angular.io/)
