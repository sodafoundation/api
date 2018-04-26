import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StoragePoolsTableComponent } from './storage-pools-table.component';

describe('StoragePoolsTableComponent', () => {
  let component: StoragePoolsTableComponent;
  let fixture: ComponentFixture<StoragePoolsTableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StoragePoolsTableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StoragePoolsTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
