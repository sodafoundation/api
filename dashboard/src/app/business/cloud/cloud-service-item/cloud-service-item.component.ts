import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-cloud-service-item',
  templateUrl: './cloud-service-item.component.html',
  styleUrls: ['./cloud-service-item.component.scss']
})
export class CloudServiceItemComponent implements OnInit {

  @Input() cloudService;

  constructor() { }

  ngOnInit() {
  }

  deleteCloud(name){
    alert("delete "+name);
    //删除某个云
  }

}
