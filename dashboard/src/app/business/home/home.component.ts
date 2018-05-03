import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Http } from '@angular/http';

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
    role;
    lineData_nums;
    lineData_capacity;
    showAdminStatis = true;
    constructor(
        private http: Http
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.role = localStorage['opensds-current-user'].split("|")[0];
        if(this.role == "admin"){
            this.showAdminStatis = true;
        }else{
            this.showAdminStatis = false;
        }

        let arr = [4, 2, 2, 10, 3, 0, 4];
        // let arr = [1,4,3,2];
        // this.arraySortUpdate(arr);
        this.items = [
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Tenants"
            },
            {
                countNum: arr[1] || 　0,
                // imgName: "u288.png",
                label: "Users"
            },
            {
                countNum: arr[2] || 　0,
                // imgName: "u288.png",
                label: "Block Storages"
            },
            {
                countNum: arr[3] || 　0,
                // imgName: "u288.png",
                label: "Storage Pools"
            },
            {
                countNum: arr[4] || 　0,
                // imgName: "u288.png",
                label: "Volumes"
            },
            {
                countNum: arr[5] || 　0,
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
        this.http.get("/v1beta/ef305038-cd12-4f3b-90bd-0612f83e14ee/profiles").subscribe((res) => {
            console.log(res.json().data);
        });
    }
}
