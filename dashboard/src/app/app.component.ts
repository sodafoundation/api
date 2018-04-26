// import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Http } from '@angular/http';
import { Router } from '@angular/router';
import { I18NService } from 'app/shared/api';
// import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { MenuItem, SelectItem } from './components/common/api';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: []
})
export class AppComponent implements OnInit{
    chromeBrowser: boolean = false;

    isLogin: boolean;

    linkUrl = "";

    userRoles: SelectItem[];

    currentRole: string = "admin";

    username: string;

    password: string;

    dropMenuItems: MenuItem[];

    menuItems = [
        {
            "title": "Home",
            "description": "Update 5 minutes ago",
            "routerLink": "/home"
        },
        {
            "title": "Block Service",
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
        private router: Router
        // private I18N: I18NService,
        // private viewContainerRef: ViewContainerRef,
        // private appService: AppService,
        // private router: Router
    ){}
    
    ngOnInit() {
        if( localStorage['x-subject-token'] != '' ){
            this.isLogin = true;
            this.router.navigateByUrl("home");
            this.activeItem = this.menuItems[0];
        }else{
            this.isLogin = false;
        }

        this.userRoles = [
            { label: "Administrator", value: "admin" },
            { label: "Tenant", value: "tenant" }
        ]

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
            let req1: any = { params:{} };
            let userid = res.json().token.user.id;
            console.log(userid);
            this.http.get("/v3/users/"+ userid +"/projects", req1).subscribe((roleRES) => {
                
            })



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
                    "name": "demo",
                    "domain": { "id": "default" }
                  }
                }
            }

            this.http.post("/v3/auth/tokens", req).subscribe((r)=>{
                console.log("正式token", r.headers.get('x-subject-token') );
                localStorage['x-subject-token'] =  r.headers.get('x-subject-token');

                console.log("存储gengxin的token",  localStorage['x-subject-token']);
                if( localStorage['x-subject-token'] != '' ){
                    this.isLogin = true;
                    this.router.navigateByUrl("home");
                    this.activeItem = this.menuItems[0];
                }
            })

        });
    }

    logout() {
        localStorage['x-subject-token'] = "";
        this.isLogin = false;
        this.password = "";
        console.log("xxx",localStorage['x-subject-token'])
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
