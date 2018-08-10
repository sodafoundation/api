import { Component, Input, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Http } from '@angular/http';
import { I18NService, Utils } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { ParamStorService } from 'app/shared/api';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem, ConfirmationService} from '../../../components/common/api';

@Component({
    selector: 'tenant-detail',
    templateUrl: 'tenantDetail.html',
    styleUrls: ['tenantDetail.scss'],
    providers: [ConfirmationService],
    animations: []
})
export class TenantDetailComponent implements OnInit {
    @Input() projectID;
    @Input() projectName;
    @Input() isDetailFinished: Boolean;
    addUserDisplay: boolean=false;
    userfilter: string="";
    projectGroups = [];
    users = [];
    popSelectedUsers;
    allUsers;

    statistics_volumeSnapshots: string;
    statistics_volumes: string;
    statistics_capacity: string;

    constructor(
        private http: Http,
        private confirmationService: ConfirmationService,
        private paramStor: ParamStorService,
        public I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.listProjectGroup();
        this.projectResources();
    }

    projectResources(){
        this.http.get("/v1beta/"+ this.projectID +"/block/volumes").subscribe((res)=>{
            let originCapacity = 0;
            res.json().map(ele => {
                originCapacity += ele.size;
            })
            this.statistics_volumes = res.json().length;
            this.statistics_capacity = Utils.getDisplayGBCapacity(originCapacity);
        })

        this.http.get("/v1beta/"+ this.projectID +"/block/snapshots").subscribe((res)=>{
            this.statistics_volumeSnapshots = res.json().length;
        })
    }

    listProjectGroup(){
        this.http.get("/v3/role_assignments?scope.project.id="+ this.projectID).subscribe((res)=>{
            let arr = res.json().role_assignments;
            let newarr = [];
            let roles=[];

            // get roles
            let reqRole: any = { params:{} };
            this.http.get("/v3/roles", reqRole).subscribe((roleRES) => {
                let currentRoleName = this.projectName == "admin" ? "admin" : (this.projectName == "service" ? "service" : "Member");

                roleRES.json().roles.forEach((item, index) => {
                    if(item.name == currentRoleName){ // more role can be expand
                        let roleJson = {};
                        roleJson["id"] = item.id;
                        roleJson["name"] = item.name;
                        roles.push(roleJson);
                    }
                })

                roles.forEach((item, index)=>{
                    arr.forEach(ele => {
                        if(ele.role.id == item.id){
                            ele.role["name"] = item.name;
                            newarr.push(ele);
                        }
                    });
                })

                newarr.forEach((item, index) => {
                    if(item.group){
                        let groupJson = {};
                        groupJson["groupid"] = item.group.id;
                        groupJson["grouprole"] = item.role
                        this.projectGroups.push(groupJson);
                    }
                });

                this.listUsers();
            })
        })
    }

    showAddUsers(){
        this.addUserDisplay = true;
        this.listAllUsers();
    }

    addUsers(){
        
        let group_id;
        this.projectGroups.forEach((item)=>{
            if(item.grouprole.name == "Member"){
                group_id = item.groupid;
            }
        })

        this.popSelectedUsers.forEach((user, i) => {
            let request: any = {};
            this.http.put("/v3/groups/"+ group_id +"/users/"+ user.id, request).subscribe((r) => {
                if(i == (this.popSelectedUsers.length-1)){
                    this.listUsers();
                    this.addUserDisplay = false;
                }
            })
        });

    }

    listAllUsers(){
        this.popSelectedUsers = [];
        this.allUsers = [];
        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        if(this.userfilter != ""){
            request.params["name"] = this.userfilter;
        }

        this.http.get("/v3/users", request).subscribe((res) => {
            let newarr = [];
            res.json().users.map((item, index) => {
                item["description"] = item.description == '' ? '--' : item.description;
                newarr.push(item);
            });

            //Filter added users
            if(this.users.length > 0){
                this.users.forEach((addedUser) => {
                    this.allUsers = newarr.filter((user, idx, arr)=>{
                        return (user.name != 'admin' && user.name != 'opensds'  && user.name != addedUser.name);
                    })
                })
            }else{
                this.allUsers = newarr.filter((user, idx, arr)=>{
                    return (user.name != 'admin' && user.name != 'opensds');
                })
            }

        });
    }

    listUsers(){
        this.users = [];
        this.isDetailFinished = false;
        this.projectGroups.forEach((item, index)=>{
            let request: any = { params:{} };
            this.http.get("/v3/groups/"+ item.groupid +"/users", request).subscribe((userRES)=>{
                userRES.json().users.forEach((ele) => {
                    ele["role"] = item.grouprole.name;
                    ele["description"] = ele.description == '' ? '--' : ele.description;
                    this.users.push(ele);
                });

                this.isDetailFinished = true;
            })
        })

    }

    removeUser(user){
        let group_id;
        this.projectGroups.forEach((item)=>{
            if(item.grouprole.name == user.role){
                group_id = item.groupid;
            }
        })

        let msg = "<div>Are you sure you want to remove the user?</div><h3>[ "+ user.name +" ]</h3>"
        this.confirmationService.confirm({
            message: msg,
            header: "Remove user",
            acceptLabel: "Remove",
            accept: ()=>{
                let request: any = {};
                this.http.delete("/v3/groups/"+ group_id +"/users/"+ user.id, request).subscribe((r) => {
                    this.listUsers();
                })
            },
            reject:()=>{}
        })
    }

}
