import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener, ViewChildren  } from '@angular/core';
import { Http } from '@angular/http';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem, ConfirmationService } from '../../components/common/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl} from '@angular/forms';

let _ = require("underscore");

@Component({
    selector: 'tenant-list',
    templateUrl: 'tenantList.html',
    styleUrls: [],
    providers: [ConfirmationService],
    animations: []
})
export class TenantListComponent implements OnInit {
    tenants = [];
    isDetailFinished = false;
    createTenantDisplay = false;
    isEditTenant = false;

    selectedTenants = [];

    sortField: string;
    currentTenant: string;
    popTitle: string;

    tenantFormGroup;
    projectID: string;
    projectName: string;

    validRule= {
        'name':'^[a-zA-Z]{1}([a-zA-Z0-9]|[_]){0,127}$'
    };

    constructor(
        private http: Http,
        private confirmationService: ConfirmationService,
        public I18N: I18NService,
        // private router: Router,
        private fb: FormBuilder
    ) {
        this.tenantFormGroup = this.fb.group({
            "form_name": ["", [Validators.required, Validators.pattern(this.validRule.name), this.ifTenantExisting(this.tenants)] ],
            "form_description":["", Validators.maxLength(200) ],
        })
    }

    errorMessage = {
        "form_name": { required: "Username is required.", pattern:"Beginning with a letter with a length of 1-128, it can contain letters / numbers / underlines.", ifTenantExisting:"Tenant is existing."},
        "form_description": { maxlength: "Max. length is 200."}
    };

    ngOnInit() {
        this.listTenants();

    }

    ifTenantExisting (param: any): ValidatorFn{
        return (c: AbstractControl): {[key:string]: boolean} | null => {
            let isExisting= false;
            this.tenants.forEach(element => {
                if(element.name == c.value){
                    isExisting = true;
                }
            })
            if(isExisting){
                return {'ifTenantExisting': true};
            }else{
                return null;
            }
        }
    }

    listTenants() {
        this.tenants=[];
        this.selectedTenants = [];

        this.sortField = "name";

        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        this.http.get("/v3/projects", request).subscribe((res) => {
            this.tenants = res.json().projects;
            this.tenants.forEach((item)=>{
                item["description"] = item.description == '' ? '--' : item.description;
                if(item.name == "admin" || item.name == "service"){
                    item["disabled"] = true;
                }
            })
        });
    }

    showCreateTenant(tenant?) {
        this.createTenantDisplay = true;

        //Reset form
        this.tenantFormGroup.reset();

        if(tenant){
            this.isEditTenant = true;
            this.popTitle = "Modify";

            this.currentTenant = tenant.id;

            this.tenantFormGroup.controls['form_name'].value = tenant.name;
            this.tenantFormGroup.controls['form_description'].value = tenant.description;

        }else{
            this.isEditTenant = false;
            this.popTitle = "Create";

        }
    }

    createTenant(){
        let request: any = { project:{} };
        request.project = {
            "domain_id": "default",
            "name": this.tenantFormGroup.value.form_name,
            "description": this.tenantFormGroup.value.form_description
        }
        
        if(this.tenantFormGroup.status == "VALID"){
            // create tenant
            this.http.post("/v3/projects", request).subscribe((res) => {
                let tenantid = res.json().project.id;
                

                // get roles
                let request: any = { params:{} };
                this.http.get("/v3/roles", request).subscribe((roleRES) => {
                    roleRES.json().roles.forEach((item, index) => {
                        if(item.name == "Member"){
                            // create group for role named [Member]
                            let request: any = { group:{} };
                            request.group = {
                                "domain_id": "default",
                                "name": "group_"+ tenantid + "_Member"
                            }
                            this.http.post("/v3/groups/", request).subscribe((groupRES) => {
                                let groupid = groupRES.json().group.id;

                                // Assign role to group on project
                                let reqRole: any = { };
                                this.http.put("/v3/projects/"+ tenantid +"/groups/"+ groupid +"/roles/"+  item.id, reqRole).subscribe(() => {
                                    this.createTenantDisplay = false;
                                    this.listTenants();
                                })
                            });
                        }
                    })

                })

                
            });
        }else{
            // validate
            for(let i in this.tenantFormGroup.controls){
                this.tenantFormGroup.controls[i].markAsTouched();
            }
        }
    }

    updateTenant(){
        let request: any = { project:{} };
        request.project = {
            "domain_id": "default",
            "name": this.tenantFormGroup.value.form_name,
            "description": this.tenantFormGroup.value.form_description
        }
        
        if(this.tenantFormGroup.status == "VALID"){
            this.http.patch("/v3/projects/"+ this.currentTenant, request).subscribe((res) => {
                this.createTenantDisplay = false;
                this.listTenants();
            });
        }
    }

    deleteTenant(tenants){
        let arr=[],msg;
        if(_.isArray(tenants)){
            tenants.forEach((item,index)=> {
                arr.push(item.id);
            })
            msg = "<div>Are you sure you want to delete the selected tenants?</div><h3>[ "+ tenants.length +" Tenants ]</h3>";
        }else{
            arr.push(tenants.id);
            msg = "<div>Are you sure you want to delete the tenant?</div><h3>[ "+ tenants.name +" ]</h3>"
        }
        
        this.confirmationService.confirm({
            message: msg,
            header: "Delete Tenant",
            acceptLabel: "Delete",
            isWarning: true,
            accept: ()=>{
                arr.forEach((ele)=> {
                    this.http.get("/v3/role_assignments?scope.project.id="+ ele).subscribe((res)=>{
                        res.json().role_assignments.forEach((item, index) => {
                            if(item.group){
                                let request: any = {};
                                this.http.delete("/v3/groups/"+ item.group.id, request).subscribe();
                            }
                        });

                        let request: any = {};
                        this.http.delete("/v3/projects/"+ ele, request).subscribe((r) => {
                            this.listTenants();
                        })
                    })
                })
            },
            reject:()=>{}
        })

    }

    onRowExpand(evt) {
        this.isDetailFinished = false;
        this.projectID = evt.data.id;
        this.projectName = evt.data.name;
    }

    tablePaginate() {
        this.selectedTenants = [];
    }

    label: object = {
        tenantNameLabel: 'Name',
        descriptionLabel: 'Description',
    }

}






