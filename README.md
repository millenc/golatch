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
if response, err := latch.Status("AccountID", false, false); err == nil {
	//Handle response
}
```
Upon success, this method will return a `response` struct of type `LatchStatusResponse`. If the second parameter is `true` (nootp), then the One-time password information will not be included in the response. If the third parameter is `true` (silent) Latch will not send a push notification to the user alerting of the access if the account's latch is on (this requires a SILVER, GOLD or PLATINUM subscription). You can easily access the information for the root application/operation:

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
if response, err := latch.OperationStatus("AccountID", "MyOperationID", false, false); err == nil {
    //Handle response
}
```
If the third parameter is `true` (nootp), then no One-time password information will be included in the response. If the fourth parameter is `true` (silent) Latch will not send a push notification to the user alerting of the access if the operation's latch is on (this requires a SILVER, GOLD or PLATINUM subscription). The response is of the same type and contains the same information as the one returned by the `Status()` method.

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

Get all the operations:

``` go
if response, err := latch.ShowOperation(""); err == nil {
	for id, operation := range response.Operations() {
		name := operation.Name
		two_factor := operation.TwoFactor
		lock_on_request := operation.LockOnRequest
		operations := operation.Operations
	}
}
```

Get single operation (using it's operation ID):

``` go
if response, err := latch.ShowOperation(operationId); err == nil {
	id, operation := response.FirstOperation()
	
	name := operation.Name
	two_factor := operation.TwoFactor
	lock_on_request := operation.LockOnRequest
	operations := operation.Operations
}
```

### History

*NOTE*: This method require a GOLD or PLATINUM subscription in order to work. If you don't have any of these types of subscriptions you will get an error.

You can get the history of an account (with it's `accountId`) between two dates (`from` and `to`) with the `History()` method:

``` go
if response, err := latch.History(accountId, from, to); err == nil {
	application := response.Application()
	lastSeen := response.LastSeen()
	clientVersion := response.ClientVersion()
	historyCount := response.HistoryCount()
	history := response.History()	
}
```
once you have a response, you can use the following methods to get the information contained in it:

* `Application()`: returns a struct of type `LatchApplication` with information about the application. This struct has the following fields:
	* `Status`: status of the application (on/off).
	* `PairedOn`: when the account was paired to the application.
	* `Name`: name of the application.
	* `Description`: description of the application.
	* `ImageURL`: URL of the application's image.
	* `ContactPhone`: contact phone of the application's administrator.
	* `ContactEmail`: contact email of the application's administrator.
	* `TwoFactor`: two factor setting.
	* `LockOnRequest`: lock on request setting.
	* `Operations`: array of LatchOperation structs containing information about the operations defined for the application.
* `LastSeen()`: last time there was user activity for this account.
* `ClientVersion()`: Contains information about the platforms and versions used by the client. The returned value is an array of structs `LatchClientVersion` where each value has the following fields:
	* `Platform`: Name of the platform ("Android" for example).
	* `App`: Version.
* `HistoryCount()`: number of history entries in the response.
* `History()`: history entries. Array of structs of type `LatchHistoryEntry`:
	* `Time`: time of this action.
	* `Action`: action (get,USER_UPDATE or DEVELOPER_UPDATE)
	* `What`: parameter that was affected by the action.
	* `Was`: previous value of the parameter.
	* `Value`: value of the parameter.
	* `Name`: name of the application or operation.
	* `UserAgent`: user agent of the user that performed the action.
	* `IP`: ip of the user that performed the action.

##User API Usage

Starting with API version 1.0 there's a User API that you can use to manage applications and get information about your subscription. The usage is pretty similar to the application API described in the previous section. The main diference is that instead of using the Application ID you have to use your User ID. Please note that all the functions described in this section require a GOLD or PLATINUM subscription in order to work.

To start using this API, you have to create an instance of the `LatchUser` struct, that you will you use later to call the appropriate methods:

``` go
import "golatch"

// ...
latch := golatch.NewLatchUser("MyUserID", "MySecretKey")
``` 

### Show applications

To get information about all the applications that you have defined in your account you can use the `ShowApplications()` function:

``` go
if response, err := latch.ShowApplications(); err == nil {
	applications := response.Applications()

	for applicationId, application := range applications {
		fmt.Println(applicationId)
		fmt.Println(application.Name)
		fmt.Println(application.Secret)
		fmt.Println(application.TwoFactor)
		fmt.Println(application.LockOnRequest)
		fmt.Println(application.ContactPhone)
		fmt.Println(application.ContactEmail)
		fmt.Println(application.ImageURL)
	}
}
```

This call will return a struct of type `LatchShowApplicationsResponse`. You can use the `Applications()` method of this struct to get the applications information as a map, where the key is the Application ID and the value is a struct of type `LatchApplicationInfo` with the following fields:

* `Name`: Name of the application.
* `Secret`: Secret key of the application.
* `TwoFactor`: Two factor authentication.
* `LockOnRequest`: Lock on request setting.
* `ContactPhone`: Contact phone.
* `ContactEmail`: Contact email.
* `ImageURL`: Image URL
* `Operations`: Map of `LatchOperation` structs with the application's operations.

### Add application

You can create a new application programatically using the `AddApplication()` method. You need to provide a `LatchApplicationInfo` struct with the application information like in the following example:

``` go
applicationInfo := &golatch.LatchApplicationInfo{
	Name:          "My Application Name",
	ContactEmail:  "my@contact_email.com",
	ContactPhone:  "111222333",
	TwoFactor:     golatch.DISABLED, //optional, can also be golatch.MANDATORY, golatch.OPT_IN or golatch.NOT_SET (empty)
	LockOnRequest: golatch.OPT_IN,   //optional, same values as TwoFactor
}

if response, err := latch.AddApplication(applicationInfo); err == nil {
	fmt.Println(response.AppID())
	fmt.Println(response.Secret())
}
```

If everything goes well, you can use the `AppID()` and `Secret()` methods on the response struct to get the Application ID and Secret Key of the new application.

### Update application

You can update the information of an existing application using the `UpdateApplication()` method. You must specify the Application's ID (first parameter) and the information you want to change (second parameter, which must be a struct of type `LatchApplicationInfo`). For example to change the name of the application you could this:

``` go
applicationInfo := &golatch.LatchApplicationInfo{
	Name: "Modified Application Name",
}

if err := latch.UpdateApplication("ApplicationID here", applicationInfo); err != nil {
	fmt.Println(err)
}
```

Because of the way the API works, only non empty values will be modified. For example if you don't fill the ContactPhone field and that field has a value, it will retain that value after the update.

### Delete an application

You can delete an existing application using the `DeleteApplication()` function:

```go
if err := latch.DeleteApplication("ApplicationID"); err == nil {
	//Application was deleted successfully!
}
```

### Subscription

You can get your current subscription information using the `Subscription()` method:

```go
if response, err := latch.Subscription(); err == nil {
	ID := response.ID()
	applications := response.Applications()
	users := response.Users()
	operations := response.Operations()
}
```

This method will return a struct of type `LatchSubscriptionResponse` with the following methods (that you can use to get information about your subscription):

* `ID()`: Gets the ID of your account/subscription type (for example "vip" for a VIP account).
* `Applications()`: Returns a struct of type `LatchSubscriptionUsage` with the following fields:
	* `InUse`: Number of applications currently being used.
	* `Limit`: Max number of applications that you can create (-1 for no limit).
* `Users()`: Returns a struct of type `LatchSubscriptionUsage` with the current number of users (InUse) and max number of users allowed (Limit).
* `Operations()`: Returns a map of `LatchSubscriptionUsage` keyed by application name that contains the current number of operations (InUse) and the max number of operations for each application (Limit).

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