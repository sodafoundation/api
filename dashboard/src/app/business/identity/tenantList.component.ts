import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener, ViewChildren  } from '@angular/core';
import { Http } from '@angular/http';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem, ConfirmationService } from '../../components/common/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl} from '@angular/forms';


@Component({
    selector: 'tenant-list',
    templateUrl: 'tenantList.html',
    styleUrls: [
        'dialogcss.css'
    ],
    providers: [ConfirmationService],
    animations: []
})
export class TenantListComponent implements OnInit {
    tenants = [];
    isDetailFinished = false;
    createTenantDisplay = false;
    isEditTenant = false;

    sortField: string;
    currentTenant: string;
    popTitle: string;

    tenantFormGroup;
    projectID: string;

    constructor(
        private http: Http,
        private confirmationService: ConfirmationService,
        // private I18N: I18NService,
        // private router: Router,
        private fb: FormBuilder
    ) {
        this.tenantFormGroup = this.fb.group({
            "form_name": ["", Validators.required ],
            "form_description":["", Validators.maxLength(200) ],
        })
    }

    errorMessage = {
        "form_name": { required: "Username is required."},
        "form_description": { maxlength: "Max. length is 200."}
    };

    ngOnInit() {
        this.listTenants();

    }

    listTenants() {
        this.sortField = "name";

        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        this.http.get("/v3/projects", request).subscribe((res) => {
            this.tenants = res.json().projects;
        });
    }

    showCreateTenant(tenant?) {
        this.createTenantDisplay = true;
        if(tenant){
            this.isEditTenant = true;
            this.popTitle = "Modify";

            this.currentTenant = tenant.id;

            this.tenantFormGroup = this.fb.group({
                "form_name": [tenant.name, Validators.required ],
                "form_description":[tenant.description, Validators.maxLength(200) ],
            })

        }else{
            this.isEditTenant = false;
            this.popTitle = "Create";

            this.tenantFormGroup = this.fb.group({
                "form_name": ["", Validators.required ],
                "form_description":["", Validators.maxLength(200) ],
            })
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
                
                // create group
                let request: any = { group:{} };
                request.group = {
                    "domain_id": "default",
                    "name": "group_"+ tenantid
                }
                this.http.post("/v3/groups/", request).subscribe((groupRES) => {
                    let groupid = groupRES.json().group.id;

                    // get roles
                    let request: any = { params:{} };
                    this.http.get("/v3/roles", request).subscribe((roleRES) => {
                        let roleid;
                        roleRES.json().roles.forEach((item, index) => {
                            if(item.name == "Member"){
                                roleid = item.id;
                            }
                        })

                        // Assign role to group on project
                        let reqRole: any = { };
                        this.http.put("/v3/projects/"+ tenantid +"/groups/"+ groupid +"/roles/"+ roleid, reqRole).subscribe(() => {
                            this.createTenantDisplay = false;
                            this.listTenants();
                        })
                    });
                })

                
            });
        }else{

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

    onRowExpand(evt) {
        this.isDetailFinished = false;

        this.projectID = evt.data.id;
        // this.http.get("/v3/role_assignments?scope.project.id="+ evt.data.id).subscribe((res)=>{
        //     let arr = res.json().role_assignments;
        //     let newarr = [];
        //     let roles=[];
        //     let groups=[];

        //     // get roles
        //     let reqRole: any = { params:{} };
        //     this.http.get("/v3/roles", reqRole).subscribe((roleRES) => {
        //         roleRES.json().roles.forEach((item, index) => {
        //             if(item.name == "Member"){ // more role can be expand
        //                 let roleJson = {};
        //                 roleJson["id"] = item.id;
        //                 roleJson["name"] = item.name;
        //                 roles.push(roleJson);
        //             }
        //         })

        //         roles.forEach((item, index)=>{
        //             arr.forEach(ele => {
        //                 if(ele.role.id == item.id){
        //                     ele.role["name"] = item.name;
        //                     newarr.push(ele);
        //                 }
        //             });
        //         })

        //         newarr.forEach((item, index) => {
        //             if(item.group){
        //                 let groupJson = {};
        //                 groupJson["groupid"] = item.group.id;
        //                 groupJson["grouprole"] = item.role
        //                 groups.push(groupJson);
        //             }
        //         });

        //         this.userGroups = groups;
        //     })
        // })
    }

    label: object = {
        tenantNameLabel: 'Name',
        descriptionLabel: 'Description',
    }

}






