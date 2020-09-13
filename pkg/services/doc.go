// Package services defines an abstract service provider.
// Services are simple interfaces consumer may book or
// unbook. On demand these services can be executed
// concurrently. The implementor of a service has to
// take care how needed information is retrieved or
// passed back.
package services
