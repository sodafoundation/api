import { Component, OnInit, AfterContentChecked, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener, ViewChildren } from '@angular/core';
import { Http } from '@angular/http';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem, ConfirmationService } from '../../components/common/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl} from '@angular/forms';

let _ = require("underscore");

@Component({
    selector: 'user-list',
    templateUrl: 'userList.html',
    styleUrls: [
        'dialogcss.css'
    ],
    providers: [ConfirmationService],
    animations: []
})
export class UserListComponent implements OnInit, AfterContentChecked{
    tenantUsers = [];
    tenantLists = [];
    createUserDisplay = false;
    isUserDetailFinished = false;
    isEditUser = false;
    userfilter = "";
    myFormGroup;

    selectedUsers = [];

    username: string;
    userRole: string;
    currentUser;

    detailUserInfo: string;
    popTitle: String;

    sortField: string;

    constructor(
        private http: Http,
        private confirmationService: ConfirmationService,
        // private I18N: I18NService,
        // private router: Router,
        private fb: FormBuilder
    ) { 
        this.myFormGroup = this.fb.group({
            "form_username": ["", Validators.required ],
            "form_description":["", Validators.maxLength(200) ],
            "form_tenant": [""],
            "form_isModifyPsw": [""],
            "form_psw": ["", Validators.required ],
            "form_pswConfirm": ["", Validators.required ]
        })
    }

    errorMessage = {
        "form_username": { required: "Username is required."},
        "form_description": { maxlength: "Max. length is 200."},
        "form_tenant": { required: "Tenant is required."},
        "form_psw": { required: "Password is required." },
        "form_pswConfirm": { required: "Password is required." }
    };
    
    label:object = {
        userNameLabel:'Username',
        passwordLabel:'Password',
        descriptionLabel:'Description',
        confirmPasswordLabel:'Confirm Password',
        roleLabel:'Role',
        tenantLabel:'Tenant'
    }
    
    showUserForm(user?): void{
        if(user){
            this.isEditUser = true;
            this.popTitle = "Modify";

            this.username = user.username;
            this.currentUser = user;

            this.createUserDisplay = true;
            this.getTenants();
            console.log(user)
            this.myFormGroup = this.fb.group({
                "form_username": [user.username, Validators.required ],
                "form_description":[user.description, Validators.maxLength(200) ],
                "form_tenant": [user.defaultTenant, Validators.required],
                "form_isModifyPsw": [false]
            })
            

        }else{
            this.isEditUser = false;
            this.popTitle = "Create";

            this.createUserDisplay = true;
            this.getTenants();
            this.getRoles();

            this.myFormGroup = this.fb.group({
                "form_username": ["", Validators.required ],
                "form_description":["", Validators.maxLength(200) ],
                "form_tenant": ["",Validators.required],
                "form_isModifyPsw": [true],
                "form_psw": ["", Validators.required ],
                "form_pswConfirm": ["", Validators.required ]
            })
        }
    }

    createUser(){
        let request: any = { user:{} };
        request.user = {
            "default_project_id": this.myFormGroup.value.form_tenant,
            "domain_id": "default",
            "name": this.myFormGroup.value.form_username,
            "description": this.myFormGroup.value.form_description,
            "password": this.myFormGroup.value.form_psw
        }
        
        if(this.myFormGroup.status == "VALID"){
            this.http.post("/v3/users", request).subscribe((res) => {
                let userInfo = res.json().user;
                let request: any = {};
                this.http.put("/v3/projects/"+ userInfo.default_project_id +"/users/"+ userInfo.id +"/roles/"+ this.userRole, request).subscribe((r) => {
                    this.createUserDisplay = false;
                    this.listUsers();
                })
            });
        }else{

        }
    }

    updateUser(){
        let request: any = { user:{} };
        request.user = {
            "description": this.myFormGroup.value.form_description,
            "default_project_id": this.myFormGroup.value.form_tenant
        }
        if(this.myFormGroup.value.form_isModifyPsw==true){
            request.user["password"] = this.myFormGroup.value.form_psw;
        }
        
        console.log(request);
        if(this.myFormGroup.status == "VALID"){
            this.http.patch("/v3/users/"+ this.currentUser.userid, request).subscribe((res) => {
                this.createUserDisplay = false;
                this.listUsers();
            });
        }
    }
    
    getRoles(){
        let request: any = { params:{} };
        this.http.get("/v3/roles", request).subscribe((res) => {
            res.json().roles.forEach((item, index) => {
                if(item.name == "Member"){
                    this.userRole = item.id;
                }
            })
        });
    }
    getTenants(){
        this.tenantLists = [];

        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        this.http.get("/v3/projects", request).subscribe((res) => {
            res.json().projects.map((item, index) => {
                let tenant = {};
                tenant["label"] = item.name;
                tenant["value"] = item.id;
                this.tenantLists.push(tenant);
            });
        });
    }


    ngOnInit() {
        this.listUsers();
        
    }

    ngAfterContentChecked(){
    //     if(this.isEditUser == true && this.myFormGroup.value.form_isModifyPsw == true){
    //         this.myFormGroup = this.fb.group({
    //             "form_username": [this.currentUser.username, Validators.required ],
    //             "form_description":[this.currentUser.description, Validators.maxLength(200) ],
    //             "form_tenant": [this.currentUser.defaultTenant, Validators.required],
    //             "form_isModifyPsw": [true],
    //             "form_psw": ["", Validators.required ],
    //             "form_pswConfirm": ["", Validators.required ]
    //         })
    //     }
    }

    listUsers(){
        this.tenantUsers = [];

        this.sortField = "username";

        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        if(this.userfilter != ""){
            request.params["name"] = this.userfilter;
        }

        this.http.get("/v3/users", request).subscribe((res) => {
            res.json().users.map((item, index) => {
                let user = {};
                user["enabled"] = item.enabled;
                user["username"] = item.name;
                user["userid"] = item.id;
                user["defaultTenant"] = item.default_project_id;
                user["description"] = item.description;
                this.tenantUsers.push(user);
            });
        });
    }

    userStatus(userid, isEnabled){
        let msg = isEnabled == true ? "Are you sure you want to disable this user?" : "Are you sure you want to enable this user?";
        let status = isEnabled ? false : true;

        this.confirmationService.confirm({
            message: msg,
            header: "Confirm",
            icon: "fa fa-question-circle",
            accept: ()=>{
                let request: any = { user:{} };
                request.user = {
                    "enabled": status
                }
                this.http.patch("/v3/users/"+ userid, request).subscribe((res) => {
                    this.listUsers();
                });
                
            },
            reject:()=>{}
        })
    }

    deleteUsers(users){
        let arr=[];
        if(_.isArray(users)){
            users.forEach((item,index)=> {
                arr.push(item.userid);
            })
        }else{
            arr.push(users);
        }

        this.confirmationService.confirm({
            message: "Are you sure you want to delete users?",
            header: "Confirm",
            icon: "fa fa-question-circle",
            accept: ()=>{
                arr.forEach((item,index)=> {
                    let request: any = {};
                    this.http.delete("/v3/users/"+ item, request).subscribe((res) => {
                        if(index == arr.length-1){
                            this.listUsers();
                        }
                    });
                })
                
            },
            reject:()=>{}
        })

    }

    onRowExpand(evt) {
        this.isUserDetailFinished = false;
        this.detailUserInfo = evt.data.userid;
    }
}
