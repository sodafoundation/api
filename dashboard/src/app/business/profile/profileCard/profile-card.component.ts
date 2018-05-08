import { Router } from '@angular/router';
import { Component, OnInit, Input } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

import { ButtonModule } from './../../../components/common/api';

// import {CardModule} from 'primeng/card';

@Component({
    selector: 'profile-card',
    templateUrl: './profile-card.component.html',
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
export class ProfileCardComponent implements OnInit {

    @Input() data;
    
    chartDatas: any;
    constructor(
        // private I18N: I18NService,
        // private router: Router
    ) { }
    option = {};
    ngOnInit() {
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
                display: false,
                width: '5px',
                position: 'right',
                fontSize: 12
            }
        };
    }

    index;
    isHover;

    showSuspensionFrame(event){
        if(event.type === 'mouseenter'){
            this.isHover = true;
        }else if(event.type === 'mouseleave'){
            this.isHover = false;
        }
        let arrLength = event.target.parentNode.children.length;
        for(let i=0; i<arrLength; i++) {
            if(event.target.parentNode.children[i] === event.target){
                this.index = i;
            }
        }
    }
}
