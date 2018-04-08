import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';

@Component({
    templateUrl: './home.component.html',
    styleUrls: []
})
export class HomeComponent implements OnInit{

    constructor(
        // private I18N: I18NService,
        // private router: Router
    ){}
    
    ngOnInit() {
        let arr = [3,10,11,17,21,23,4];
        // let arr = [1,4,3,2];
        // this.arraySortUpdate(arr);

      
    }
    
    // arraySort(array){
    //     for(let i = 1; i<array.length; i++){
    //         for(let j = i-1; j>=0; j--){
    //             if(array[j+1] < array[j]){
    //                 let tempArr = [array[j+1], array[j]];
    //                 array.splice(j, 2, tempArr[0],tempArr[1]);
    //             }else{
    //                 break;
    //             }
    //         }
    //     }
    //     console.log(array);
    //     return array;
    // }

    // arraySortUpdate(array){
    //     for(let i = 1; i<array.length; i++){
    //         var temmArr = array[i];
    //         var j;
    //         for(j=i; j>0 && temmArr < array[j-1]; j--){
    //             array[j] = array[j-1];
    //         }
    //         array[j] = temmArr;
    //     }
    //     console.log(array);
    //     return array;
    // }
}
