import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SnapshortListComponent } from './snapshort-list.component';

describe('SnapshortListComponent', () => {
  let component: SnapshortListComponent;
  let fixture: ComponentFixture<SnapshortListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SnapshortListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SnapshortListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
