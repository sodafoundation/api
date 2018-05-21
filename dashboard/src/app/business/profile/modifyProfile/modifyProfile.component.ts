import { Router,ActivatedRoute} from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

import { ProfileService, PoolService} from './../profile.service';

@Component({
    templateUrl: './modifyProfile.component.html',
    styleUrls: [

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
    totalFreeCapacity = 0;
    option;
    data;
    cars;
    cols: any[];
    formGroup;
    errorMessage;
    pools;
    totalCapacity = 0;
    // profileId;
    constructor(
        // private I18N: I18NService,
        // private router: Router
        private ActivatedRoute: ActivatedRoute,
        private ProfileService: ProfileService,
        private poolService:PoolService
    ) { }
    ngOnInit() {
        let profileId;
        this.ActivatedRoute.params.subscribe((params) => profileId = params.profileId);

        this.ProfileService.getProfileById(profileId).subscribe((res) => {
            // return res.json();
            // this.profiles = res.json();
            this.data = res.json();
            this.getPools();
        });
        this.items = [
            { label: 'Profile', url: '/profile' },
            { label: 'Profile detail', url: '/modifyProfile' }
        ];
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
        // this.data = {
        //     "id": "5d8c3732-a248-50ed-bebc-539a6ffd25c1",
        //     "name": "Gold",
        //     "protocol": "FC",
        //     "type": "Thin",
        //     "policys": [
        //         "Qos",
        //         "Snapshot",
        //         "Replication"
        //     ],
        //     "description": "provide gold storage service",
        //     "extras": {
        //         "key1": "value1",
        //         "key2": {
        //             "subKey1": "subValue1",
        //             "subKey2": "subValue2"
        //         },
        //         "key3": "value3"
        //     }
        // };

    }
    getPools() {
        this.poolService.getPools().subscribe((res) => {
            this.pools = res.json();
            this.totalFreeCapacity = this.getSumCapacity(this.pools, 'free');
            this.totalCapacity = this.getSumCapacity(this.pools, 'total');
            this.chartDatas = {
                labels: ['Total Capacity', 'Used Capacity'],
                datasets: [
                    {
                        data: [this.totalCapacity, this.totalCapacity-this.totalFreeCapacity],
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
        });
    }

    getSumCapacity(pools, FreeOrTotal) {
        let SumCapacity: number = 0;
        let arrLength = pools.length;
        for (let i = 0; i < arrLength; i++) {
            if(this.data.extras.protocol.toLowerCase() == pools[i].extras.ioConnectivity.accessProtocol &&  this.data.storageType == pools[i].extras.dataStorage.provisioningPolicy){
                if (FreeOrTotal === 'free') {
                    SumCapacity += pools[i].freeCapacity;
                } else {
                    SumCapacity += pools[i].totalCapacity;
                }
            }
        }
        return SumCapacity;
    }
}
