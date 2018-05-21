## Summary

OpenSDS dashboard uses the front-end development framework Angular 2 and relies on PrimeNG UI Components. Regardless of deployment or two development, prepare the corresponding environment.

## Environment

1.Install NodeJS
Open the official website of NodeJS and download the latest package for installation.
When the installation is complete, open a terminal/console window, enter and launch：
- node -v
- npm -v
See whether the installation is successful or not.

2.Open the "dashboard" directory，open a terminal/console window in the directory, enter and launch：
- npm install
- npm install -g @angular/cli

3.Start OpenSDS Dashboard
- npm start

4.Build OpenSDS Dashboard
- ng build --prod

## Deplyment

In the dashboard directory, find the directory named dist, copy the contents of the directory to the root directory of the web server.

## Links
[Primeng] https://www.primefaces.org/primeng/
[Primeng Github] https://github.com/primefaces/primeng/
[NodeJS] https://nodejs.org/
[Angular] https://angular.io/
[Angular/cli] https://cli.angular.io/

