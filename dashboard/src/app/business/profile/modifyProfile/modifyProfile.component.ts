import { Router,ActivatedRoute} from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

@Component({
    templateUrl: './modifyProfile.component.html',
    styleUrls: [
        "./modifyProfile.component.scss"
    ],
    animations: [
        trigger('overlayState', [
            state('hidden', style({
                opacity: 0
            })),
            state('visible', style({
                opacity: 1
            })),
            transition('visible => hidden', animate('400ms ease-in')),
            transition('hidden => visible', animate('400ms ease-out'))
        ]),

        trigger('notificationTopbar', [
            state('hidden', style({
                height: '0',
                opacity: 0
            })),
            state('visible', style({
                height: '*',
                opacity: 1
            })),
            transition('visible => hidden', animate('400ms ease-in')),
            transition('hidden => visible', animate('400ms ease-out'))
        ])
    ]
})
export class modifyProfileComponent implements OnInit {
    items;
    chartDatas;
    totalFreeCapacity;
    option;
    data;
    cars;
    cols: any[];
    formGroup;
    errorMessage;
    pools;
    profile;
    constructor(
        // private I18N: I18NService,
        // private router: Router
        private ActivatedRoute: ActivatedRoute
    ) { }
    ngOnInit() {

        this.ActivatedRoute.params.subscribe((params) => console.log(params));
        // this.profile = params.username
        console.log
        this.items = [
            { label: 'Profile', url: '/profile' },
            { label: 'Profile detail', url: '/modifyProfile' }
        ];
        this.chartDatas = {
            labels: ['Total Capacity', 'Used Capacity'],
            datasets: [
                {
                    data: [(1000 - 300), 300],//未使用容量放前面
                    backgroundColor: [
                        "rgba(224, 224, 224, 1)",
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
            // rotation: (-0.2 * Math.PI),
            title: {
                display: false,
                text: 'My Title',
                fontSize: 12
            },
            legend: {
                display: true,
                labels:{
                    boxWidth:12
                },
                position: 'bottom',
                fontSize: 12
            }
        };
        this.data = {
            "id": "5d8c3732-a248-50ed-bebc-539a6ffd25c1",
            "name": "Gold",
            "protocol": "FC",
            "type": "Thin",
            "policys": [
                "Qos",
                "Snapshot",
                "Replication"
            ],
            "description": "provide gold storage service",
            "extras": {
                "key1": "value1",
                "key2": {
                    "subKey1": "subValue1",
                    "subKey2": "subValue2"
                },
                "key3": "value3"
            }
        };

        this.cols = [
            { field: 'name', header: 'Name' },
            { field: 'freeCapacity', header: 'FreeCapacity' },
            { field: 'totalCapacity', header: 'TotalCapacity' },
            { field: 'dockId', header: 'Disk' },
            { field: 'storageType', header: 'StorageType' }
        ];

        this.pools =[
            {
              "id": "string",
              "createdAt": "2018-04-11T08:11:27.335Z",
              "updatedAt": "2018-04-11T08:11:27.335Z",
              "name": "string1",
              "storageType": "string",
              "description": "string",
              "status": "string",
              "availabilityZone": "string",
              "totalCapacity": 0,
              "freeCapacity": 1,
              "dockId": "string",
              "extras": {
                "additionalProp1": {},
                "additionalProp2": {},
                "additionalProp3": {}
              }
            },
            {
                "id": "string",
                "createdAt": "2018-04-11T08:11:27.335Z",
                "updatedAt": "2018-04-11T08:11:27.335Z",
                "name": "string2",
                "storageType": "string",
                "description": "string",
                "status": "string",
                "availabilityZone": "string",
                "totalCapacity": 0,
                "freeCapacity": 8,
                "dockId": "string",
                "extras": {
                  "additionalProp1": {},
                  "additionalProp2": {},
                  "additionalProp3": {}
                }
              },
              {
                "id": "string",
                "createdAt": "2018-04-11T08:11:27.335Z",
                "updatedAt": "2018-04-11T08:11:27.335Z",
                "name": "string3",
                "storageType": "string",
                "description": "string",
                "status": "string",
                "availabilityZone": "string",
                "totalCapacity": 0,
                "freeCapacity": 10,
                "dockId": "string",
                "extras": {
                  "additionalProp1": {},
                  "additionalProp2": {},
                  "additionalProp3": {}
                }
              }
          ]
          this.totalFreeCapacity = this.getSumCapacity(this.pools,'free');
    }

    getSumCapacity(pools,FreeOrTotal){
        let SumCapacity:number = 0;
        let arrLength = pools.length;
        for(let i=0;i<arrLength;i++){
            if(FreeOrTotal==='free'){
                SumCapacity += pools[i].freeCapacity;
            }else{
                SumCapacity += pools[i].totalCapacity;
            }
        }
        return SumCapacity;
    }
}
