package v1

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Philipid3s/go-rpc/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

type contactServiceServer struct {
	db *sql.DB
}

// NewContactServiceServer creates Contacts service
func NewContactServiceServer(db *sql.DB) v1.ContactServiceServer {
	return &contactServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *contactServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *contactServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Read Contact task
func (s *contactServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query Contact by ID
	rows, err := c.QueryContext(ctx, "SELECT id, firstname, lastname, address FROM Contacts WHERE id=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Contacts-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Contacts-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Contact with ID='%d' is not found",
			req.Id))
	}

	// get Contact data
	var td v1.Contact
	if err := rows.Scan(&td.Id, &td.Firstname, &td.Lastname, &td.Address); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Contact row-> "+err.Error())
	}

	return &v1.ReadResponse{
		Api:     apiVersion,
		Contact: &td,
	}, nil

}

// Read all Contacts
func (s *contactServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get Contacts list
	rows, err := c.QueryContext(ctx, "SELECT id, firstname, lastname, address FROM Contacts")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Contacts-> "+err.Error())
	}
	defer rows.Close()

	list := []*v1.Contact{}
	for rows.Next() {
		td := new(v1.Contact)
		if err := rows.Scan(&td.Id, &td.Firstname, &td.Lastname, &td.Address); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from Contact row-> "+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Contact-> "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:     apiVersion,
		Contact: list,
	}, nil
}
