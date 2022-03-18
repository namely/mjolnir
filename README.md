# Mjolnir
 
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/30f14d8bd9864619bf404699d92682f3)](https://app.codacy.com/gh/namely/mjolnir/dashboard)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/30f14d8bd9864619bf404699d92682f3)](https://app.codacy.com/gh/namely/mjolnir/dashboard)
[![Build Status](https://travis-ci.org/namely/mjolnir.svg?branch=master)](https://travis-ci.org/namely/mjolnir)
[![Go Report Card](https://goreportcard.com/badge/namely/mjolnir)](https://goreportcard.com/report/namely/mjolnir)

Mjolnir is Namely's common Go packages. It contains packages for things like health checks, database utils, logging, etc.

### Logging

The Logger returns an interceptor that will set a *logrus.Entry on the context.

#### Errors

The logger makes calls to your gRPC endpoints assuming they return a protobuf response or an error.

    out, err := handler(ctx, req)

When building gRPC Services, sometimes you may want return user friendly protobuf error responses,
 and other times you may want to return generic internal errors (e.g. when your DB write fails).
 
 If your proto Error is set up like:

    message Error {
      string key = 1;
      string message = 2;
    }
 
 Then included in this repo is the ability to customize error responses to fit both conditions.
  Both KeyedErr and FieldErr can be passed around like standard errors because they implement the error interface.

 For example, say you have a standard Get User endpoint.
 In an error condition you decide to return a friendly "user not found" message if the user doesn't exist,
 or otherwise a generic internal server error (your client doesn't need to know about internal errors!).
          
 Then you can build your proto error response with something like: 
  
      var ErrUserNotFound = &KeyedErr{
          ErrorKey: "not_found",
          Message:  "user not found",
      }
      
      &proto.Error{
           Key:     ErrUserNotFound.Key(),
           Message: ErrUserNotFound.Error(),
       },
       
 And otherwise a generic response like:   
    
    var grpcError = data.NewFieldErr(
        &logrus.Fields{
            "error":  err,
        },
    )
