import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ReplicationListComponent } from './replication-list.component';

describe('ReplicationListComponent', () => {
  let component: ReplicationListComponent;
  let fixture: ComponentFixture<ReplicationListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ReplicationListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ReplicationListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
