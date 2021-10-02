# Public Error (Go)

This is a transparent error wrapping library using a custom error type which
allows us to inject a "public" error message at any point in the error cause
chain without breaking support for errors.Is() or errors.Unwrap().

The most common use-case is server responses for failed operations in which
you want to return a different message to the client than the one you log
with sensitive application details.

Example:

```go

// Pretend we failed to find the user in our database query because we tried 
// to commit a rolled-back transaction (or some other bug)
err := sql.ErrTxDone

// this error might be wrapped by the model
err = fmt.Errorf("user: insert: %w", err)

// At some point in the call stack, wrap the error with the message you want the client to see
err = publicerror.Wrap(err, "Sorry, we were unable to create your account", http.StatusNotFound)

// Perhaps it's wrapped several more times as well

...

// then the error finally makes it back to the http handler that needs to inform the client

// 1) Log our error like normal so our devs can see the problem (ignores the public message above)
log.Println(err)

// 2) Show the client a "safe" error or http.StatusInternalServerError
http.Error(w, publicerror.Message(err), publicerror.StatusCode(err))
```

The public message (somewhere in the chain of errors) is for the user, the rest 
are for the logs. The idea is that you can still put plenty of debugging 
information in the full error logs - but show a nice, safe message to the 
client.
