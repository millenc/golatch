#golatch [![Build Status](https://travis-ci.org/millenc/golatch.svg?branch=master)](https://travis-ci.org/millenc/golatch)

**golatch** is a package that easily lets you integrate [Latch](https://latch.elevenpaths.com/) in your [Go](https://golang.org/) applications. If you don't know what Latch is or how to use it please visit their [official site](https://latch.elevenpaths.com/).

##Installation


You can download the package manually and put it in your `/src` folder or use `go get`:

``` bash
$ go get github.com/millenc/golatch 
```
You can also use `go get -u` to update the package. 

##Usage

First you need to create the Latch struct that you will use to call all the operations of the API:

``` go
import "golatch"

// ...
latch := golatch.NewLatch("MyAppID", "MySecretKey")
``` 
where "MyAppID" and "MySecretKey" are your application's ID and secret key respectively. You can find this information in the developer's area of the official website once you create your application.

###Pairing

You can use the `Pair()` method providing the pairing token supplied by the user:

``` go
response, err := latch.Pair("MyToken")
```
If everything goes well, this call will return a response containing the Account ID that has been paired:

``` go
if err == nil {
	account_id := response.AccountId()
	//Store the account id
}
```
You must store this account ID together with the user information to use it in future API calls.

###Unpairing

Call the `Unpair()` method and pass the Account ID you want to unpair as argument:

``` go
if err := latch.Unpair("AccountID"); err != nil {
	//Handle error
}
```

If no error is returned, the account has been unpaired successfully.

###Locking/unlocking

*NOTE*: These methods require a GOLD or PLATINUM subscription in order to work. If you don't have any of these types of subscriptions you will get an error.

You can change the status of an account (lock/unlock) using the `Lock()`/`Unlock()` methods:

``` go
err := latch.Lock("AccountID")
err := latch.Unlock("AccountID")
```
If you have created operations for an application you can lock and unlock them using the `LockOperation()` and `UnlockOperation()` methods:

``` go
err := latch.LockOperation("AccountID", "MyOperationID")
err := latch.UnlockOperation("AccountID", "MyOperationID")
```

###Status

To get the current status of an account, you can use the `Status()` method:

``` go
if response, err := latch.Status("AccountID", false); err == nil {
	//Handle response
}
```
Upon success, this method will return a `response` struct of type `LatchStatusResponse`. If the second parameter is `true`, then no One-time password information will be included in the response. You can easily access the information for the root application/operation:

**`Status()`**: Gets the account status:

``` go
status := response.Status()
if status == golatch.LATCH_STATUS_ON {
	//Latch is ON for this account, don't let the user log in (for example)
}
```
As you can see in the previous snippet, to determine the status of an account it's recommended that you use the `golatch.LATCH_STATUS_ON` (ON status) and `golatch.LATCH_STATUS_OFF` (OFF status) constants.

**`TwoFactor()`**: Gets the two factor authentication information:

``` go
two_factor LatchTwoFactor = reponse.TwoFactor()
if two_factor.Token != "" {
	token     := two_factor.Token
	generated := two_factor.Generated
	//Handle two factor authentication
}
```
Please note that this method will always return a `LatchTwoFactor` struct, even if the response didn't include this information. Hence you should always test for zero values before proceeding (`two_factor.Token != ""`).

**`Operations()`**: Gets the child operations of the application:

``` go
operations := reponse.Operations()
```
An application (or operation) can have child operations. This method will return the status for these nested operations as a map of `LatchOperationStatus` indexed by operation ID. Each of these `LatchOperationStatus` structs have the following information:

``` go
status     := operation_status.Status
two_factor := operation_status.TwoFactor
operations := operation_status.Operations
```
The information contained in these fields is the same that has been previously discussed.

###Operation Status

**IMPORTANT**: If you have created operations for your application you should not query the status of the application, but rather the status of individual operations. To do so, you can use the `OperationStatus()` method:

``` go
if response, err := latch.OperationStatus("AccountID", "MyOperationID", false); err == nil {
    //Handle response
}
```
If the third parameter is `true`, then no One-time password information will be included in the response. The response is of the same type and contains the same information as the one returned by the `Status()` method.

###Managing operations

You can create/edit/delete operations directly from your application:

**Create operation**:

``` go
if response, err := latch.AddOperation(parentId, name, twoFactor, lockOnRequest); err == nil {
	//Get the ID of the newly created operation
	operationId := response.OperationId()
}
```
where:

* `parentId`: Id of this operation's parent (application or another operation).
* `name`: Name of the operation.
* `twoFactor`: Use of two factor authentication. One of these values:
	* OPT_IN (optional): `golatch.OPT_IN`
	* MANDATORY (mandatory): `golatch.MANDATORY`
	* DISABLED (disabled): `golatch.DISABLED`
* `lockOnRequest`: Takes the same values as twoFactor. 

**Modify operation**:

``` go
err := latch.UpdateOperation(operationId, name, twoFactor, lockOnRequest)
```
where:

* `operationId`: Id of the operation you want to modify.
* `name`: Name of the operation.
* `twoFactor`: This is optional, use the value `golatch.NOT_SET` if you want to leave the existing value.
* `lockOnRequest`: This is optional, use the value `golatch.NOT_SET` if you want to leave the existing value.

**Delete operation**:

``` go
err := DeleteOperation(operationId)
```

**Get operation information**:

``` go
if response, err := latch.ShowOperation(operationId); err == nil {
	operation := response.Operation()
	
	name := operation.Name
	two_factor := operation.TwoFactor
	lock_on_request := operation.LockOnRequest
	operations := operation.Operations
}
```

##Tests

You can run unit tests for this package using:

``` bash
$ go test golatch
```

##Documentation

Documentation is provided by *godoc* as usual:

``` bash
$ godoc golatch
```
You can also view the godoc's HTML documentation starting a web server (`godoc -http=:6060`) and navigating to the URL `http://localhost:6060/pkg/golatch/`.

Aditionally you can browse the documentation online [here](http://godoc.org/github.com/millenc/golatch).

##Contributions

All contributions are welcome! This is my first Go package ever, so if you're an experienced Go developer and want to share some advice or suggest any improvement to the code it will be greatly appreciated.