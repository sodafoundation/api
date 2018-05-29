import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-suspension-frame',
  templateUrl: './suspension-frame.component.html',
  styleUrls: [
    
  ]
})
export class SuspensionFrameComponent implements OnInit {

  @Input() policy;

  constructor() { }

  ngOnInit() {
    // console.log(this.policy);
  }

}
