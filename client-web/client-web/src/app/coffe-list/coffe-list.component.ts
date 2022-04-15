import { AfterViewInit, ViewChild, Component, OnInit } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { Observable, Subject } from "rxjs";
import { HttpClient } from "@angular/common/http";
import { debounceTime } from "rxjs/operators";
import { distinctUntilChanged } from "rxjs/operators";
import { switchMap } from "rxjs/operators";


export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  SKU: string;
}

export interface HttpRequest {
  products: Product[];
}

@Component({
  selector: 'app-coffe-list',
  templateUrl: './coffe-list.component.html',
  styleUrls: ['./coffe-list.component.css']
})
export class CoffeListComponent implements OnInit, AfterViewInit {
  displayedColumns: string[] = ['id', 'name', 'description', 'price', 'sku'];
  @ViewChild(MatPaginator) paginator!: MatPaginator;
  products$!: Observable<any>;
  productsDataSource!: MatTableDataSource<Product[]>;
  characterDatabase = new HttpDatabase(this.httpClient);
  searchTerm$ = new Subject<string>();


  constructor(private httpClient: HttpClient,) {
    this.characterDatabase
      .search(this.searchTerm$)
      .subscribe((response: any) => {
        this.productsDataSource = new MatTableDataSource(response);
        this.productsDataSource.paginator = this.paginator;
        this.products$ = this.productsDataSource.connect();
        console.log("> Response", response);
      })
  }

  ngAfterViewInit() {
    this.paginator.page.subscribe(() => {
      this.characterDatabase
        .getProducts("", "", this.paginator.pageIndex)
        .subscribe((response: any) => {
          this.productsDataSource = new MatTableDataSource(response);
          //this.resultsLength = response.info.count;
          // this.characterDataSource.paginator = this.paginator;
          this.products$ = this.productsDataSource.connect();
        });
    });
  }

  ngOnInit(): void {
    this.characterDatabase.getProducts().subscribe((response: any) => {
      this.productsDataSource = new MatTableDataSource(response);
      //this.resultsLength = response.info.count;
      this.productsDataSource.paginator = this.paginator;
      this.products$ = this.productsDataSource.connect();
    });
    //this.changeDetectorRef.detectChanges();
  }

  ngOnDestroy() {
    if (this.productsDataSource) {
      this.productsDataSource.disconnect();
    }
  }
}

export class HttpDatabase {
  constructor(private _httpClient: HttpClient) { }

  search(terms: Observable<string>) {
    return terms
      .pipe(debounceTime(400),
        distinctUntilChanged(),
        switchMap(term => this.getProducts(term)));
  }

  getProducts(
    name: string = "",
    status: string = "",
    page: number = 0
  ): Observable<HttpRequest> {
    const href = "//localhost:9090/products";

    return this._httpClient.get<HttpRequest>(href);
  }
}