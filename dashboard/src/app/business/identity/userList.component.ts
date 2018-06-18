import { Component, OnInit, AfterViewChecked, AfterContentChecked, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener, ViewChildren } from '@angular/core';
import { Http } from '@angular/http';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate, group } from '@angular/animations';
import { MenuItem, ConfirmationService } from '../../components/common/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl} from '@angular/forms';
import { retry } from 'rxjs/operators';
import { forEach } from '@angular/router/src/utils/collection';

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

    userRole: string;
    tenantLists = [];

    username: string;
    currentUser;

    detailUserInfo: string;
    popTitle: string;

    sortField: string;

    validRule= {
        'name':'^[a-zA-Z]{1}([a-zA-Z0-9]|[_]){0,127}$'
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
            "form_username": ["", {validators:[Validators.required, Validators.pattern(this.validRule.name), this.ifUserExisting(this.tenantUsers)], updateOn:'change'}  ],
            "form_description":["", Validators.maxLength(200) ],
            "form_tenant": [""],
            "form_isModifyPsw": [true],
            "form_psw": ["", {validators: [Validators.required, Validators.minLength(8), this.regPassword], updateOn:'blur'} ],
            "form_pswConfirm": ["", [Validators.required, this.regConfirmPassword(this.newPassword)] ]
        })
    }

    errorMessage = {
        "form_username": { required: "Username is required.", pattern:"Beginning with a letter with a length of 1-128, it can contain letters / numbers / underlines.", ifUserExisting:"User is existing."},
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

    ifUserExisting (param: any): ValidatorFn{
        return (c: AbstractControl): {[key:string]: boolean} | null => {
            let isExisting= false;
            this.tenantUsers.forEach(element => {
                if(element.username == c.value){
                    isExisting = true;
                }
            })
            if(isExisting){
                return {'ifUserExisting': true};
            }else{
                return null;
            }
        }
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
        this.getRoles();
        this.getTenants();
        this.createUserDisplay = true;

        //Reset form
        this.myFormGroup.reset();

        if(user){
            this.isEditUser = true;
            this.popTitle = "Modify";

            this.username = user.username;
            this.currentUser = user;

            this.myFormGroup.controls['form_description'].value = user.description;
            this.myFormGroup.controls['form_isModifyPsw'].value = false;

            this.myFormGroup.controls['form_username'].clearValidators();
            this.myFormGroup.controls['form_psw'].clearValidators();
            this.myFormGroup.controls['form_pswConfirm'].clearValidators();


        }else{
            this.isEditUser = false;
            this.popTitle = "Create";

            this.myFormGroup.controls['form_isModifyPsw'].value = true;
            this.myFormGroup.controls['form_username'].setValidators({validators:[Validators.required, Validators.pattern(this.validRule.name), this.ifUserExisting(this.tenantUsers)], updateOn:'change'});
            this.myFormGroup.controls['form_psw'].setValidators([Validators.required, Validators.minLength(8), this.regPassword]);
            this.myFormGroup.controls['form_pswConfirm'].setValidators([Validators.required, this.regConfirmPassword(this.newPassword)] );
        }

        // Update form validate status
        for(let i in this.myFormGroup.controls){
            this.myFormGroup.controls[i].updateValueAndValidity();
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
            this.http.post("/v3/users", request).subscribe( (res) => {
                let tenants = this.myFormGroup.value.form_tenant;
                
                if(!tenants){
                    this.createUserDisplay = false;
                    this.listUsers();
                    return;
                }

                tenants.forEach( (element, i) => {
                    this.http.get("/v3/role_assignments?scope.project.id="+ element).subscribe((ass_res)=>{
                        let groupid;

                        ass_res.json().role_assignments.map((element)=>{
                            if(element.group){
                                return groupid = element.group.id;
                            }
                        })

                        let request: any = {};
                        this.http.put("/v3/groups/"+ groupid + "/users/"+ res.json().user.id, request).subscribe();
                    })
                    if(i == (tenants.length-1)){
                        this.createUserDisplay = false;
                        this.listUsers();
                    }
                })
            });
        }else{
            // validate
            for(let i in this.myFormGroup.controls){
                this.myFormGroup.controls[i].markAsTouched();
            }
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
            }else{
                // validate
                for(let i in this.myFormGroup.controls){
                    this.myFormGroup.controls[i].markAsTouched();
                }
            }
        }else{

            if(this.myFormGroup.status == "VALID"){
                this.http.patch("/v3/users/"+ this.currentUser.userid, request).subscribe((res) => {
                    this.createUserDisplay = false;
                    this.listUsers();
                });
            }else{
                // validate
                for(let i in this.myFormGroup.controls){
                    this.myFormGroup.controls[i].markAsTouched();
                }
            }
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
                if(item.name != "admin" && item.name != "service"){
                    tenant["label"] = item.name;
                    tenant["value"] = item.id;
                    this.tenantLists.push(tenant);
                }
            });
        });
    }


    ngOnInit() {
        this.listUsers();

    }

    ngAfterViewChecked(){
        this.newPassword = this.myFormGroup.value.form_psw;

    }

    ModifyPswChecked(checked){
        if(checked){
            this.myFormGroup.controls['form_psw'].setValidators([Validators.required, Validators.minLength(8), this.regPassword]);
            this.myFormGroup.controls['form_pswConfirm'].setValidators([Validators.required, this.regConfirmPassword(this.newPassword)] );

        }else{
            this.myFormGroup.controls['form_psw'].clearValidators();
            this.myFormGroup.controls['form_pswConfirm'].clearValidators();

        }

        // Update form validate status
        for(let i in this.myFormGroup.controls){
            this.myFormGroup.controls[i].updateValueAndValidity();
        }
    }

    listUsers(){
        this.tenantUsers = [];
        this.selectedUsers = [];

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

                if(item.name == "admin" || item.name == "opensds"){
                    user["disabled"] = true;
                }
                this.tenantUsers.push(user);
            });
        });
    }

    userStatus(user){
        let msg = user.enabled == true ? "<div>Are you sure you want to disable the user?</div><h3>[ "+ user.username +" ]</h3>" : "<div>Are you sure you want to enable the user?</div><h3>[ "+ user.username +" ]</h3>";
        let status = user.enabled ? false : true;

        this.confirmationService.confirm({
            message: msg,
            header: user.enabled ? 'Disable User' : 'Enable User',
            acceptLabel: user.enabled ? 'Disable' : 'Enable',
            accept: ()=>{
                let request: any = { user:{} };
                request.user = {
                    "enabled": status
                }
                this.http.patch("/v3/users/"+ user.userid, request).subscribe((res) => {
                    this.listUsers();
                });

            },
            reject:()=>{}
        })
    }

    deleteUsers(users){
        let arr=[], msg;
        if(_.isArray(users)){
            users.forEach((item,index)=> {
                arr.push(item.userid);
            })
            msg = "<div>Are you sure you want to delete the selected users?</div><h3>[ "+ users.length +" Users ]</h3>";
        }else{
            arr.push(users.userid);
            msg = "<div>Are you sure you want to delete the user?</div><h3>[ "+ users.username +" ]</h3>";
        }

        this.confirmationService.confirm({
            message: msg,
            header: "Delete User",
            acceptLabel: "Delete",
            isWarning: true,
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

    tablePaginate() {
        this.selectedUsers = [];
    }
}
