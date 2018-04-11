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
        "./profile-card.css"
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
            labels: ['A', 'B'],
            datasets: [
                {
                    data: [(1000-300), 300],//未使用容量放前面
                    backgroundColor: [
                        "rgba(224, 224, 224, 1)",
                        "#FF6384"
                    ]
                    // hoverBackgroundColor: [
                    //     "#FF6384",
                    //     "#36A2EB",
                    //     "#FFCE56"
                    // ]
                }]
        };
        this.option = {
            title: {
                display: false,
                text: 'My Title',
                fontSize: 12
            },
            legend: {
                display: true,
                width: '5px',
                position: 'right',
                fontSize: 12
            }
        };
    }

}
