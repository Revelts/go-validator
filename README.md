# GoValidator


GoValidator is a Package validator implements value validations for structs and nested structs fields based on tags.


## Tags
----
- ##### [`field` ] value : required
- ##### [`fieldType`] value : number / text
- ##### [`min`] value : (string number)
- ##### [`max`] value : (string number)
- ##### [`startswith`] value : (string)
- ##### [`endswith`] value : (string)
- ##### [`value_of`] value : (string)(can accept 1 or more string value seperated with ,(coma))
- ##### [`format`] value : email /alphabet/alphanumeric
- ##### [`match`] value : regex pattern 


Tag `min` / `max` will automatically validate based on the value type.
String : count the length
Int/numeric : validate the value

## Changelog
### 30-10-2018
- Add Handler for email format when its not required
- Add startswith and Endswith validation
### 04-11-2018
- Support slice struct
- Support Nested slice struct
### 05-11-2018
- Support float32 type for min & max value
- Support required value for slice data type
### 13-11-2018
- Add handler for type:number on string value
### 06-12-2018
- Add alphabet & alphanumeric format
- Allow space on alphabet & alphanumeric format
### 13-12-2018
- Change float32 type to float 64 type
### 04-09-2019
- add return json name if there's json tag
- trim spaces for required string
----
### 16-02-2021
- Handle Unexported Field
- Add regex validation
----
### 24-02-2021
- Handle int32 data type
----
### 01-06-2023
- Refactoring code structure
- Support to handle pointer struct or data type
- Support to handle interfacing struct type
- Support to handle slice of data type validation
- Change tags "type" -> "fieldType"
----
### Installation

```
go get github.com/vincentius93/govalidator
```
## Usage
---
#### Basic
`````
type User struct {
	FirstName      string       `field:"required"`
	LastName       string       `fieldType:"text"`
	Age            int          `fieldType:"number"`
	Email          string       `format:"email"`
}

func main(){
    MyStruct := User{FirstName:"John",LastName:"S"}
    err := validator.validate(MyStruct)
    if err != nil {
        fmt.Println(err)
    }
}
`````
#### Multiple Validation
`````
type User struct {
	FirstName      string       `field:"required" fieldType:"text" max:"10"`
	LastName       string       `field:"required" fieldType:"text" min:"1"`
	Age            int          `type:"number" max:"1"`
	Email          string       `format:"email"`
}

func main(){
    MyStruct := User{FirstName:"John",LastName:"S"}
    err := validator.validate(MyStruct)
    if err != nil {
        fmt.Println(err)
    }
}
`````
#### Nested Struct validation
`````
type User struct {
	FirstName      string       `field:"required" fieldType:"text" max:"10"`
	LastName       string       `field:"required" fieldType:"text" min:"1"`
	Age            int          `fieldType:"number" max:"1"`
	Email          string       `format:"email"`
}
type Orders struct {
	OrderId             int  `field:"required" max:"10"`
	UserData            User
}

func main(){
    MyStruct := User{FirstName:"John",LastName:"S"}
    MyOrder := Orders{OrderId:9930438,UserData:MyStruct}

    err := validator.validate(MyOrder)
    if err != nil {
        fmt.Println(err)
    }
}
`````
#### Slice Struct Validation
`````
type User struct {
	FirstName      string       `field:"required" fieldType:"text" max:"10"`
	LastName       string       `field:"required" fieldType:"text" min:"1"`
	Age            int          `type:"number" max:"1"`
	Email          string       `format:"email"`
}
type Orders struct {
	OrderId             int  `field:"required" max:"10"`
	UserData            User
}

func main(){
    var MyStruct []User
    MyStruct = append(MyStruct,User{FirstName:"John",LastName:"S"})
    err := validator.validate(Mystruct)
    if err != nil {
        fmt.Println(err)
    }
}
`````
#### Nested Slice Struct Validation
`````
type User struct {
	FirstName      string       `field:"required" fieldType:"text" max:"10"`
	LastName       string       `field:"required" fieldType:"text" min:"1"`
	Age            int          `type:"number" max:"1"`
	Email          string       `format:"email"`
	Orders         []Orders
}
type Orders struct {
	OrderId             int  `field:"required" max:"10"`
}

func main(){
    var myOrders []Orders
    myOrders = append(myOrders,Orders{OrderId:9930438})
    myOrders = append(myOrders,Orders{OrderId:9930432})
    MyStruct := User{FirstName:"John",LastName:"S",Orders:myOrders}
    err := validator.validate(Mystruct)
    if err != nil {
        fmt.Println(err)
    }
}
`````# go-validator
