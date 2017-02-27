package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/opensds/opensds/api/openapi-spec/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../swagger.json

func configureFlags(api *operations.OpenSDSAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.OpenSDSAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.UrlformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	api.CreateShareHandler = operations.CreateShareHandlerFunc(func(params operations.CreateShareParams) middleware.Responder {
		return middleware.NotImplemented("operation .CreateShare has not yet been implemented")
	})
	api.CreateVolumeHandler = operations.CreateVolumeHandlerFunc(func(params operations.CreateVolumeParams) middleware.Responder {
		return middleware.NotImplemented("operation .CreateVolume has not yet been implemented")
	})
	api.DeleteShareHandler = operations.DeleteShareHandlerFunc(func(params operations.DeleteShareParams) middleware.Responder {
		return middleware.NotImplemented("operation .DeleteShare has not yet been implemented")
	})
	api.DeleteVolumeHandler = operations.DeleteVolumeHandlerFunc(func(params operations.DeleteVolumeParams) middleware.Responder {
		return middleware.NotImplemented("operation .DeleteVolume has not yet been implemented")
	})
	api.GetShareHandler = operations.GetShareHandlerFunc(func(params operations.GetShareParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetShare has not yet been implemented")
	})
	api.GetVersionv1Handler = operations.GetVersionv1HandlerFunc(func(params operations.GetVersionv1Params) middleware.Responder {
		return middleware.NotImplemented("operation .GetVersionv1 has not yet been implemented")
	})
	api.GetVolumeHandler = operations.GetVolumeHandlerFunc(func(params operations.GetVolumeParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetVolume has not yet been implemented")
	})
	api.ListShareResourcesHandler = operations.ListShareResourcesHandlerFunc(func(params operations.ListShareResourcesParams) middleware.Responder {
		return middleware.NotImplemented("operation .ListShareResources has not yet been implemented")
	})
	api.ListSharesHandler = operations.ListSharesHandlerFunc(func(params operations.ListSharesParams) middleware.Responder {
		return middleware.NotImplemented("operation .ListShares has not yet been implemented")
	})
	api.ListVersionsHandler = operations.ListVersionsHandlerFunc(func(params operations.ListVersionsParams) middleware.Responder {
		return middleware.NotImplemented("operation .ListVersions has not yet been implemented")
	})
	api.ListVolumeResourcesHandler = operations.ListVolumeResourcesHandlerFunc(func(params operations.ListVolumeResourcesParams) middleware.Responder {
		return middleware.NotImplemented("operation .ListVolumeResources has not yet been implemented")
	})
	api.ListVolumesHandler = operations.ListVolumesHandlerFunc(func(params operations.ListVolumesParams) middleware.Responder {
		return middleware.NotImplemented("operation .ListVolumes has not yet been implemented")
	})
	api.OperateVolumeHandler = operations.OperateVolumeHandlerFunc(func(params operations.OperateVolumeParams) middleware.Responder {
		return middleware.NotImplemented("operation .OperateVolume has not yet been implemented")
	})
	api.UpdateShareHandler = operations.UpdateShareHandlerFunc(func(params operations.UpdateShareParams) middleware.Responder {
		return middleware.NotImplemented("operation .UpdateShare has not yet been implemented")
	})
	api.UpdateVolumeHandler = operations.UpdateVolumeHandlerFunc(func(params operations.UpdateVolumeParams) middleware.Responder {
		return middleware.NotImplemented("operation .UpdateVolume has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
