package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	pb "gitlab.com/go/grpc/genproto/product"
)

type Product struct {
	Name       string
	Categoryid int
	Typeid     int
}
type Store struct {
	Id        int64
	Name      string
	Addresses []Address
}
type Address struct {
	Id       int
	District string
	Street   string
}

type ProductResp struct {
	ID       int64
	Name     string
	Category string
	Type     string
}

type ProductInfo struct {
	ID       int64
	Name     string
	Category string
	Type     string
}
func DelelteProduct(index int64) (error){
	connStr := "user=postgres password=compos1995 dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_,err=db.Exec(`delete from products where id=$1`,index)
	if err!=nil{
		fmt.Println("error while deleting products",err)
		return err
	}
	return nil
}


func CreateStores(n int64, stores []*pb.StoreReq) ([]Store, error) {
	connStr := "user=postgres password=compos1995 dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var storesResp []Store
	for _, store := range stores {
		var s Store
		err = db.QueryRow(`insert into storages(name) values($1) returning id,name`, store.Name).Scan(&s.Id, &s.Name)
		if err != nil {
			fmt.Println("error while inserting ", err)
			return []Store{}, err
		}
		_, err = db.Exec(`insert into product_storages (product_id,storage_id) values($1,$2)`, n, s.Id)
		if err != nil {
			fmt.Println("error while inserting", err)
			return []Store{}, err
		}
		var addresses []Address
		for _, addr := range store.Addresses {
			var addressResp Address
			err = db.QueryRow(`insert into addresses(district,street) values($1,$2) returning id,district,street`, addr.District, addr.Street).Scan(&addressResp.Id, &addressResp.District, &addressResp.Street)
			if err != nil {
				fmt.Println("error inserting addresses", err)
				return []Store{}, err
			}
			_, err = db.Exec(`insert into storage_addresses(storage_id,address_id) values($1,$2)`, s.Id, addressResp.Id)
			if err != nil {
				fmt.Println("error while storage_addresses", err)
				return []Store{}, err
			}
			addresses = append(addresses, addressResp)
		}
		s.Addresses = addresses
		storesResp = append(storesResp, s)
	}
	return storesResp, nil
}

func GetProductInfo(n int64) (*ProductInfo, error) {
	connStr := "user=postgres password=compos1995 dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var p ProductInfo
	err = db.QueryRow(`select p.id,p.name,c.name,t.name FROM products p
	INNER JOIN categories c ON c.id=p.category_id
	INNER JOIN types t ON t.id=p.type_id
	WHERE p.id=$1`, n).Scan(&p.ID, &p.Name, &p.Category, &p.Type)
	if err != nil {
		fmt.Println("error while selecting", err)
		return &ProductInfo{}, err
	}
	return &p, nil
}

func GetProducts() ([]ProductInfo, error) {
	connStr := "user=postgres password=compos1995 dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query(`select p.id,p.name,c.name,t.name FROM products p
	INNER JOIN categories c ON c.id=p.category_id
	INNER JOIN types t ON t.id=p.type_id`)
	if err != nil {
		fmt.Println("error while selecting from products", err)
		return []ProductInfo{}, err
	}
	var productsResp []ProductInfo
	for rows.Next() {
		var productResp ProductInfo
		err = rows.Scan(&productResp.ID, &productResp.Name, &productResp.Category, &productResp.Type)
		if err != nil {
			fmt.Println("error while selecting all products", err)
			return []ProductInfo{}, err
		}
		productsResp = append(productsResp, productResp)
	}
	return productsResp, nil
}

func UpdateProduct(req *pb.ProductReq) (*ProductResp, error) {
	connStr := "user=postgres password=compos1995 dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	var productResp ProductResp
	defer db.Close()
	err = db.QueryRow(`update products
	SET name=$1,category_id=$2,type_id=$3
	WHERE id=$4 returning id,name,category_id,type_id`, req.Name, req.Categoryid, req.Typeid, req.Id).Scan(&productResp.ID, &productResp.Name, &productResp.Category, &productResp.Type)
	if err != nil {
		fmt.Println("ERROR while updating products", err)
		return &ProductResp{}, err
	}
	return &productResp, nil

}

func CreateProduct(product *Product) (*ProductResp, error) {
	connStr := "user=postgres password=compos1995 dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var p ProductResp
	err = db.QueryRow(`INSERT INTO products(name,category_id,type_id)
	values ($1,$2,$3)
	RETURNING id, name`, product.Name, product.Categoryid, product.Typeid).Scan(&p.ID, &p.Name)
	if err != nil {
		fmt.Println("error while inserting products", err)
		return &ProductResp{}, err
	}
	err = db.QueryRow(`select c.name,t.name from products p inner join categories c ON c.id=p.category_id inner join types t on t.id=p.type_id where p.id=$1`, p.ID).Scan(&p.Category, &p.Type)
	if err != nil {
		fmt.Println("error while selecting from products", err)
		return &ProductResp{}, err
	}
	return &p, nil
}
