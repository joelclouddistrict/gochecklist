package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/joelclouddistrict/gochecklist/log"
	svc "github.com/joelclouddistrict/gochecklist/services"
	impl "github.com/joelclouddistrict/gochecklist/services/implementation"
	"github.com/joelclouddistrict/gochecklist/store"
	"github.com/joelclouddistrict/gochecklist/store/gopg"
	"golang.org/x/net/context"

	gw "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jllopis/getconf"
	"google.golang.org/grpc"

	"github.com/cockroachdb/cmux"
)

type MicroServer struct {
	cmuxServer   cmux.CMux
	grpcListener net.Listener
	GrpcServer   *grpc.Server
	GrpcGwMux    *gw.ServeMux
	httpListener net.Listener
	httpServer   *http.Server

	config *getconf.GetConf

	defaultStore store.Storer
}

// New crea un nuevo MicroServer a partir de las opciones de configuración proporcionadas.
// Devuelve una referencia a MicroServer. La función siempre devuelve una estructura.
// Si se produce un error en su creación, se devuelve &MicroServer{}
func New(c *getconf.GetConf) *MicroServer {
	fmt.Printf("\n%v\n", c)
	// TODO (jllopis): permitir a través de una variable especificar si ha de conectarse con
	// un Store o no.
	// Conexión con el store
	ds, err := gopg.New()
	if err != nil {
		log.Err(fmt.Sprintf("cannot create default store: %v", err))
	}
	// Dialing store
	p, err := strconv.Atoi(c.GetString("store-port"))
	if err != nil {
		log.Err(fmt.Sprintf("Unmarshal store port error: %v", err))
		p = 5432
	}

	if err := ds.Dial(store.Options{
		"host":     c.GetString("store-host"),
		"port":     p,
		"user":     c.GetString("store-user"),
		"password": c.GetString("store-pass"),
		"dbname":   c.GetString("store-name"),
		"sslmode":  "disable",
	}); err != nil {
		log.Err(fmt.Sprintf("cannot create default store: %v", err))
	}
	// Enable query logger if mode is "dev"
	mode := c.GetString("mode")
	if mode == "dev" {
		ds.EnableQueryLogger()
	}

	// Crear el MicroServer
	ms := &MicroServer{
		config:       c,
		defaultStore: ds,
	}

	// Create the main listener.
	l, err := net.Listen("tcp", ":"+ms.config.GetString("port"))
	if err != nil {
		log.Err(err.Error())
		os.Exit(1)
	}

	// Crear el muxer cmux
	ms.cmuxServer = cmux.New(l)

	// Match connections in order:
	// First grpc, and otherwise HTTP.
	ms.grpcListener = ms.cmuxServer.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// Any significa que no hay coincidencia previa
	// En nuestro caso, no es grpc así que debe ser http.
	ms.httpListener = ms.cmuxServer.Match(cmux.Any())

	// Create your protocol servers.
	ms.GrpcServer = grpc.NewServer()
	ms.GrpcGwMux = gw.NewServeMux()

	httpmux := http.NewServeMux()
	httpmux.Handle("/", logh(ms.GrpcGwMux))

	ms.httpServer = &http.Server{
		Handler: httpmux,
		// ErrorLog: logger, // do not log user error
	}

	return ms
}

// Serve inica los servicios grpc y http para escuchar peticiones
func (ms *MicroServer) Serve() error {
	// Use the muxed listeners for your servers.
	go ms.GrpcServer.Serve(ms.grpcListener)
	go ms.httpServer.Serve(ms.httpListener)
	// Start serving!
	return ms.cmuxServer.Serve()

}

func (ms *MicroServer) Register() []error {

	errors := []error{}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	svc.RegisterTodoServiceServer(ms.GrpcServer, &impl.TodoServiceServer{Store: ms.defaultStore})
	if err := svc.RegisterTodoServiceHandlerFromEndpoint(context.Background(), ms.GrpcGwMux, ms.config.GetString("gwaddr")+":"+ms.config.GetString("port"), opts); err != nil {
		errors = append(errors, err)
	}

	// Devolver los posibles errores producidos en la inicialización de los servicios
	return errors
}

func logh(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		// Registrar llamada REST
		handler.ServeHTTP(w, r)
	})
}
