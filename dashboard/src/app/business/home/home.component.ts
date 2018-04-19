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
    option;
    constructor(
        private http: Http
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        let arr = [3, 10, 11, 17, 21, 23, 4];
        // let arr = [1,4,3,2];
        // this.arraySortUpdate(arr);
        this.items = [
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Tenants"
            },
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Users"
            },
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Block Storages"
            },
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Storage Pools"
            },
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Volumes"
            },
            {
                countNum: arr[0] || 　0,
                // imgName: "u288.png",
                label: "Volume Snapshots"
            },
            {
                countNum: arr[0] || 　0,
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

    }

    showData() {
        this.http.get("/v1beta/ef305038-cd12-4f3b-90bd-0612f83e14ee/profiles").subscribe((res) => {
            console.log(res.json().data);
        });
    }
}
