// import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener, AfterViewInit } from '@angular/core';
import { Http } from '@angular/http';
import { Router } from '@angular/router';
import { I18NService, Consts, ParamStorService } from 'app/shared/api';
// import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { MenuItem, SelectItem } from './components/common/api';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: []
})
export class AppComponent implements OnInit, AfterViewInit{
    chromeBrowser: boolean = false;

    isLogin: boolean;

    hideLoginForm: boolean = false;

    linkUrl = "";

    username: string;

    password: string;

    dropMenuItems: MenuItem[];

    currentTenant: string="";

    showLoginAnimation: boolean=false;

    showLogoutAnimation: boolean=false;

    menuItems = [];

    menuItems_tenant = [
        {
            "title": "Home",
            "description": "Update 5 minutes ago",
            "routerLink": "/home"
        },
        {
            "title": "Volume",
            "description": "23 volumes",
            "routerLink": "/block"
        }
    ]

    menuItems_admin = [
        {
            "title": "Home",
            "description": "Update 5 minutes ago",
            "routerLink": "/home"
        },
        {
            "title": "Volume",
            "description": "23 volumes",
            "routerLink": "/block"
        },
        // {
        //     "title": "Multi-Cloud Service",
        //     "description": "5 replications, 1 migrations",
        //     "routerLink": "/cloud"
        // },
        {
            "title": "Profile",
            "description": "7 profiles have been created",
            "routerLink": "/profile"
        },
        {
            "title": "Resource",
            "description": "5 storages, 2 availability zone",
            "routerLink": "/resource"
        },
        {
            "title": "Identity",
            "description": "Managing tenants and users",
            "routerLink": "/identity"
        }
    ];

    activeItem: any;

    private msgs: any = [{ severity: 'warn', summary: 'Warn Message', detail: 'There are unsaved changes'}];

    constructor(
        private el: ElementRef,
        private http: Http,
        private router: Router,
        private paramStor: ParamStorService
        // private I18N: I18NService,
        // private viewContainerRef: ViewContainerRef,
        // private appService: AppService,
        // private router: Router
    ){}
    
    ngOnInit() {
        let currentUserInfo = this.paramStor.CURRENT_USER();
        if(currentUserInfo != ""){
            this.username = this.paramStor.CURRENT_USER().split("|")[0];
            this.currentTenant = this.paramStor.CURRENT_TENANT().split("|")[0];

            if(this.username=="admin"){
                this.menuItems = this.menuItems_admin;
            }else{
                this.menuItems = this.menuItems_tenant;
            }

            this.isLogin = true;
            this.hideLoginForm = true;
            // this.router.navigateByUrl("home");
            this.activeItem = this.menuItems[0];
        }else{
            this.isLogin = false;
            this.hideLoginForm = false;
        }

        this.dropMenuItems = [
            { 
                label: "Switch Tenant", 
                items:[
                    {
                        label: "TenantA", command:()=>{}
                    },
                    {
                        label: "TenantB", command:()=>{}
                    }
                ]
            },
            { label: "Logout", command:()=>
                {
                    this.logout();
                }
            }
        ];

        
    }

    ngAfterViewInit(){
        this.loginBgAnimation();
    }

    loginBgAnimation(){
        let obj =this.el.nativeElement.querySelector(".login-bg");
        if(obj){
            let obj_w = obj.clientWidth;
            let obj_h = obj.clientHeight;
            let dis = 50;
            obj.addEventListener("mousemove", (e)=>{
                let MX = e.clientX;
                let MY = e.clientY;
                let offsetX = (obj_w - 2258)*0.5 + (obj_w-MX)*dis / obj_w;
                let offsetY = (obj_h - 1363)*0.5 + (obj_h-MY)*dis / obj_h;
                obj.style.backgroundPositionX = offsetX +"px";
                obj.style.backgroundPositionY = offsetY +"px";
            })
        }
    }
    
    login() {
        let request: any = { auth: {} };
        request.auth = {
            "identity": {
                "methods": [
                    "password"
                ],
                "password":{
                    "user": {
                        "name": this.username,
                        "domain": {
                            "name": "Default"
                        },
                        "password": this.password
                    }
                }
            }
        }

        this.http.post("/v3/auth/tokens", request).subscribe((res)=>{
            this.paramStor.AUTH_TOKEN(res.headers.get('x-subject-token'));
            let userid = res.json().token.user.id;

            // Get user owned tenants
            let reqUser: any = { params:{} };
            this.http.get("/v3/users/"+ userid +"/projects", reqUser).subscribe((objRES) => {
                // Get token authentication with scoped
                let g_token_id = res.headers.get('x-subject-token'); 
                let req: any = { auth: {} };
                req.auth = {
                    "identity": {
                        "methods": [
                            "token"
                        ],
                        "token": {
                            "id": g_token_id
                        }
                    },
                    "scope": {
                    "project": {
                        "name": objRES.json().projects[0].name,
                        "domain": { "id": "default" }
                    }
                    }
                }

                this.http.post("/v3/auth/tokens", req).subscribe((r)=>{
                    this.paramStor.AUTH_TOKEN( r.headers.get('x-subject-token') );
                    if( this.paramStor.AUTH_TOKEN() != '' ){
                        this.paramStor.CURRENT_TENANT(objRES.json().projects[0].name + "|" +objRES.json().projects[0].id);
                        this.paramStor.CURRENT_USER(this.username + "|"+ userid);

                        if(this.username == "admin"){
                            this.menuItems = this.menuItems_admin;
                        }else{
                            this.menuItems = this.menuItems_tenant;
                        }

                        this.isLogin = true;
                        this.currentTenant = objRES.json().projects[0].name;
                        this.router.navigateByUrl("home");
                        this.activeItem = this.menuItems[0];

                        // annimation for after login
                        this.showLoginAnimation = true;
                        setTimeout(() => {
                            this.showLoginAnimation = false;
                            this.hideLoginForm = true;
                        }, 500);
                    }
                })
            })
        });
    }

    logout() {
        this.paramStor.AUTH_TOKEN("");
        this.paramStor.CURRENT_USER("");
        this.paramStor.CURRENT_TENANT("");
        
        // annimation for after login
        this.hideLoginForm = false;
        this.showLogoutAnimation = true;
        setTimeout(() => {
            this.showLogoutAnimation = false;
            this.password = "";
            this.isLogin = false;
        }, 500);

    }

    menuItemClick(event, item) {
        this.activeItem = item;
    }

    supportCurrentBrowser(){
        let ie,
            firefox,
            safari,
            chrome,
            cIE = 11,
            cFirefox = 40,
            cChrome = 40;
        let ua = navigator.userAgent.toLowerCase();
        let isLinux = (ua.indexOf('linux') >= 0);

        if(this.isIE()) {
            if(ua.indexOf('msie') >= 0) {
                ie = this.getSys(ua.match(/msie ([\d]+)/));
            } else {
                ie = this.getSys(ua.match(/trident.*rv:([\d]+)/));
            }
        }else if(navigator.userAgent.indexOf("Firefox") > 0){
            firefox = this.getSys(ua.match(/firefox\/([\d]+)/));
        }else if(ua.indexOf("safari") != -1 && !(ua.indexOf("chrome") != -1)) {
            safari = this.getSys(ua.match(/version\/([\d]+)/));
        }else if(ua.indexOf("chrome") != -1) {
            chrome = this.getSys(ua.match(/chrome\/([\d]+)/));
        }

        if ((firefox) / 1 < cFirefox || (chrome) / 1 < cChrome || (ie) / 1 < cIE) {
            return true;
        }

        return false;
    }

    isIE() {
        return navigator.userAgent.toLowerCase().indexOf('trident') >= 0;
    }

    getSys (browserVersionArr) {
        if( !browserVersionArr) {
            return 0;
        } else if( browserVersionArr.length < 2) {
            return 0;
        } else {
            return browserVersionArr[1];
        }
    }
}
