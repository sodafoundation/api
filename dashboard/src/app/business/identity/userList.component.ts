import { Component, OnInit, AfterViewChecked, AfterContentChecked, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener, ViewChildren } from '@angular/core';
import { Http } from '@angular/http';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem, ConfirmationService } from '../../components/common/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl} from '@angular/forms';
import { retry } from 'rxjs/operators';

let _ = require("underscore");

@Component({
    selector: 'user-list',
    templateUrl: 'userList.html',
    styleUrls: [],
    providers: [ConfirmationService],
    animations: []
})
export class UserListComponent implements OnInit, AfterViewChecked {
    tenantUsers = [];
    createUserDisplay = false;
    isUserDetailFinished = false;
    isEditUser = false;
    myFormGroup;

    selectedUsers = [];

    username: string;
    currentUser;

    detailUserInfo: string;
    popTitle: String;

    sortField: string;

    validRule= {
        'name':'^[a-zA-Z]{1}([a-zA-Z0-9]|[_]){2,127}$'
    };

    newPassword = "";

    constructor(
        private http: Http,
        private confirmationService: ConfirmationService,
        // private I18N: I18NService,
        // private router: Router,
        private fb: FormBuilder
    ) { 
        

        this.myFormGroup = this.fb.group({
            "form_username": ["", [Validators.required, Validators.pattern(this.validRule.name)] ],
            "form_description":["", Validators.maxLength(200) ],
            "form_tenant": [""],
            "form_isModifyPsw": [""],
            "form_psw": ["", [Validators.required, Validators.minLength(8), this.regPassword]],
            "form_pswConfirm": ["", [Validators.required, this.regConfirmPassword(this.newPassword) ] ]
        })
    }

    errorMessage = {
        "form_username": { required: "Username is required.", pattern:"Beginning with a letter with a length of 3-20, it can contain letters / numbers / underlines."},
        "form_description": { maxlength: "Max. length is 200."},
        "form_tenant": { required: "Tenant is required."},
        "form_psw": { required: "Password is required.", minlength: "At least two kinds of letters / numbers / special characters, min. length is 8.", regPassword:"At least two kinds of letters / numbers / special characters, min. length is 8." },
        "form_pswConfirm": { required: "Password is required.", regConfirmPassword: "Two inputted password inconsistencies." }
    };
    
    label:object = {
        userNameLabel:'Username',
        passwordLabel:'Password',
        descriptionLabel:'Description',
        confirmPasswordLabel:'Confirm Password',
        roleLabel:'Role',
        tenantLabel:'Tenant'
    }
    
    regPassword(c:AbstractControl):{[key:string]:boolean} | null {
        let reg1 = /.*[a-zA-Z]+.*/;
        let reg2 = /.*[0-9]+.*/;
        let reg3 = /.*[\ \`\~\!\@\#\$\%\^\&\*\(\)\-\_\=\+\\\|\[\{\}\]\;\:\'\"\,\<\.\>\/\?]+.*/;
        if( !reg1.test(c.value) && !reg2.test(c.value) ){
            return {'regPassword': true};
        }
        if( !reg1.test(c.value) && !reg3.test(c.value) ){
            return {'regPassword': true};
        }
        if( !reg2.test(c.value) && !reg3.test(c.value) ){
            return {'regPassword': true};
        }
        return null;
    }

    regConfirmPassword (param: any): ValidatorFn{
        return (c: AbstractControl): {[key:string]: boolean} | null => {
            if(c.value != this.newPassword){
                return {'regConfirmPassword': true};
            }
            return null;
        }
    }

    showUserForm(user?): void{
        if(user){
            this.isEditUser = true;
            this.popTitle = "Modify";

            this.username = user.username;
            this.currentUser = user;

            this.createUserDisplay = true;
            
            this.myFormGroup = this.fb.group({
                "form_description":[user.description, Validators.maxLength(200) ],
                "form_isModifyPsw": [false],
                "form_psw": ["",  [Validators.required, Validators.minLength(8), this.regPassword] ],
                "form_pswConfirm": ["",  [Validators.required, this.regConfirmPassword(this.newPassword) ]]
            })
            

        }else{
            this.isEditUser = false;
            this.popTitle = "Create";

            this.createUserDisplay = true;

            this.myFormGroup = this.fb.group({
                "form_username": ["",  [Validators.required, Validators.pattern(this.validRule.name)] ],
                "form_description":["", Validators.maxLength(200) ],
                "form_isModifyPsw": [true],
                "form_psw": ["", [Validators.required, Validators.minLength(8), this.regPassword] ],
                "form_pswConfirm": ["", [Validators.required, this.regConfirmPassword(this.newPassword) ] ]
            })
        }
    }

    createUser(){
        let request: any = { user:{} };
        request.user = {
            "domain_id": "default",
            "name": this.myFormGroup.value.form_username,
            "description": this.myFormGroup.value.form_description,
            "password": this.myFormGroup.value.form_psw
        }
        
        if(this.myFormGroup.status == "VALID"){
            this.http.post("/v3/users", request).subscribe((res) => {
                this.createUserDisplay = false;
                this.listUsers();
            });
        }else{

        }
    }

    updateUser(){
        let request: any = { user:{} };
        request.user = {
            "description": this.myFormGroup.value.form_description
        }
        if(this.myFormGroup.value.form_isModifyPsw==true){
            request.user["password"] = this.myFormGroup.value.form_psw;

            if(this.myFormGroup.status == "VALID"){
                this.http.patch("/v3/users/"+ this.currentUser.userid, request).subscribe((res) => {
                    this.createUserDisplay = false;
                    this.listUsers();
                });
            }
        }else{
            if(this.myFormGroup.controls['form_description'].valid == true){
                this.http.patch("/v3/users/"+ this.currentUser.userid, request).subscribe((res) => {
                    this.createUserDisplay = false;
                    this.listUsers();
                });
            }
        }
        
        
    }
    
    // getRoles(){
    //     let request: any = { params:{} };
    //     this.http.get("/v3/roles", request).subscribe((res) => {
    //         res.json().roles.forEach((item, index) => {
    //             if(item.name == "Member"){
    //                 this.userRole = item.id;
    //             }
    //         })
    //     });
    // }

    // getTenants(){
    //     this.tenantLists = [];

    //     let request: any = { params:{} };
    //     request.params = {
    //         "domain_id": "default"
    //     }

    //     this.http.get("/v3/projects", request).subscribe((res) => {
    //         res.json().projects.map((item, index) => {
    //             let tenant = {};
    //             tenant["label"] = item.name;
    //             tenant["value"] = item.id;
    //             this.tenantLists.push(tenant);
    //         });
    //     });
    // }


    ngOnInit() {
        this.listUsers();
        
    }
    
    ngAfterViewChecked(){
        this.newPassword = this.myFormGroup.value.form_psw;
    }

    listUsers(){
        this.tenantUsers = [];

        this.sortField = "username";

        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        this.http.get("/v3/users", request).subscribe((res) => {
            res.json().users.map((item, index) => {
                let user = {};
                user["enabled"] = item.enabled;
                user["username"] = item.name;
                user["userid"] = item.id;
                user["defaultTenant"] = item.default_project_id;
                user["description"] = !item.description ? '--' : item.description=='' ? '--' : item.description;
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
