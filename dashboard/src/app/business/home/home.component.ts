import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Http } from '@angular/http';
import { ParamStorService } from 'app/shared/api';
import {Observable} from "rxjs/Rx";

@Component({
    templateUrl: './home.component.html',
    styleUrls: [
        './home.component.scss'
    ]
})
export class HomeComponent implements OnInit {
    items = [];
    chartDatas;
    chartDatasbar;
    option;
    lineData_nums;
    lineData_capacity;
    showAdminStatis = true;
    tenants =[];
    constructor(
        private http: Http,
        private paramStor: ParamStorService
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        if(this.paramStor.CURRENT_USER().split("|")[0] == "admin"){
            this.showAdminStatis = true;
        }else{
            this.showAdminStatis = false;
        }

        let arr = [4, 2, 2, 10, 3, 0, 4];
        // let arr = [1,4,3,2];
        // this.arraySortUpdate(arr);
        this.items = [
            {
                countNum: 0,
                // imgName: "u288.png",
                label: "Tenants"
            },
            {
                countNum:0,
                // imgName: "u288.png",
                label: "Users"
            },
            {
                countNum: arr[2] || 　0,
                // imgName: "u288.png",
                label: "Block Storages"
            },
            {
                countNum: 0,
                // imgName: "u288.png",
                label: "Storage Pools"
            },
            {
                countNum: 0,
                // imgName: "u288.png",
                label: "Volumes"
            },
            {
                countNum: 0,
                // imgName: "u288.png",
                label: "Volume Snapshots"
            },
            {
                countNum: arr[6] || 　0,
                // imgName: "u288.png",
                label: "Volume Replications"
            },
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Cross-Region Replications"
            },
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Cross-Region Migrations"
            }
        ];
        this.getCountData();

        this.chartDatas = {
            labels: ['Unused Capacity', 'Used Capacity'],
            datasets: [
                {
                    label: 'high_capacity',
                    data: [(1000 - 300), 300],//已使用300，总容量1000
                    backgroundColor: [
                        "rgba(224, 224, 224, .5)",
                        "#438bd3"
                    ]
                    // hoverBackgroundColor: [
                    //     "#FF6384",
                    //     "#36A2EB",
                    //     "#FFCE56"
                    // ]
                }]
        };
        this.chartDatasbar = {
            labels: ['high_capacity'],
            datasets: [
                {
                    label: 'Total Capacity',
                    backgroundColor: '#42A5F5',
                    borderColor: '#1E88E5',
                    data: [65]
                },
                {
                    label: 'Free Capacity',
                    backgroundColor: '#9CCC65',
                    borderColor: '#7CB342',
                    data: [28]
                }]
        };
        this.option = {
            cutoutPercentage: 80,
            // rotation: (0.5 * Math.PI),
            // circumference: (Math.PI),
            title: {
                display: false,
                text: 'My Title',
                fontSize: 12
            },
            legend: {
                labels: {
                    boxWidth: 12
                },
                display: true,
                position: 'right',
                fontSize: 12
            }
        };

        this.lineData_capacity = {
            labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
            datasets: [
                {
                    label: 'Capacity(GB)',
                    data: [10, 11, 20, 160, 156, 195, 200],
                    fill: false,
                    borderColor: '#4bc0c0'
                }
            ]
        }

        this.lineData_nums = {
            labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
            datasets: [
                {
                    label: 'Volumes',
                    data: [10, 23, 40, 38, 86, 107, 190],
                    fill: false,
                    borderColor: '#565656'
                }
            ]
        }

    }

    showData() {
        this.http.get("/v1beta/{project_id}/profiles").subscribe((res) => {
            console.log(res.json().data);
        });
    }
    listTenants() {
        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        this.http.get("/v3/projects", request).subscribe((res) => {

            this.items[0].countNum = res.json().projects.length;
            this.tenants = res.json().projects;
            this.tenants.forEach((item)=>{
                this.getAllvolumes(item.id);
                this.getAllSnapshots(item.id);
                this.getAllpools(item.id);
            });
        });
    }
    listUsers(){
        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }
        this.http.get("/v3/users", request).subscribe((res) => {
            this.items[1].countNum = res.json().users.length;
        });
    }
    getAllvolumes(projectId){
        let url = 'v1beta/'+projectId+'/block/volumes';
        this.http.get(url).subscribe((res)=>{
            this.items[4].countNum = this.items[4].countNum + res.json().length;
        });
    }
    getAllSnapshots(projectId){
        let url = 'v1beta/'+projectId+'/block/snapshots';
        this.http.get(url).subscribe((res)=>{
            this.items[5].countNum = this.items[5].countNum + res.json().length;
        });
    }
    getAllpools(projectId){
        let url = 'v1beta/'+projectId+'/pools';
        this.http.get(url).subscribe((res)=>{
            this.items[3].countNum = this.items[3].countNum + res.json().length;
        });
    }
    getCountData(){
        this.listTenants();
        this.listUsers();
    }
}
