package v1

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Philipid3s/go-rpc/pkg/api/v1"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// TestContactServiceServerRead Test function Read using go-sqlmock to simulate a true SQL server connection
// https://github.com/DATA-DOG/go-sqlmock
func TestContactServiceServerRead(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewContactServiceServer(db)

	type args struct {
		ctx context.Context
		req *v1.ReadRequest
	}
	tests := []struct {
		name    string
		s       v1.ContactServiceServer
		args    args
		mock    func()
		want    *v1.ReadResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "address"}).
					AddRow(1, "julien", "regnault", "Singapore")
				mock.ExpectQuery("SELECT (.+) FROM Contacts").WithArgs(1).WillReturnRows(rows)
			},
			want: &v1.ReadResponse{
				Api: "v1",
				Contact: &v1.Contact{
					Id:        1,
					Firstname: "title",
					Lastname:  "description",
					Address:   "address",
				},
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "SELECT failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM Contacts").WithArgs(1).
					WillReturnError(errors.New("SELECT failed"))
			},
			wantErr: true,
		},
		{
			name: "Not found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "address"})
				mock.ExpectQuery("SELECT (.+) FROM Contacts").WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Read(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("contactServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("contactServiceServer.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestContactServiceServerReadAll Test function ReadAll using go-sqlmock to simulate a true SQL server connection
// https://github.com/DATA-DOG/go-sqlmock
func TestContactServiceServerReadAll(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewContactServiceServer(db)

	type args struct {
		ctx context.Context
		req *v1.ReadAllRequest
	}
	tests := []struct {
		name    string
		s       v1.ContactServiceServer
		args    args
		mock    func()
		want    *v1.ReadAllResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1",
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "address"}).
					AddRow(1, "julien", "regnault", "Singapore").
					AddRow(2, "marc", "jeandau", "Geneva")
				mock.ExpectQuery("SELECT (.+) FROM Contacts").WillReturnRows(rows)
			},
			want: &v1.ReadAllResponse{
				Api: "v1",
				Contact: []*v1.Contact{
					{
						Id:        1,
						Firstname: "julien",
						Lastname:  "regnault",
						Address:   "Singapore",
					},
					{
						Id:        2,
						Firstname: "marc",
						Lastname:  "jeandau",
						Address:   "Geneva",
					},
				},
			},
		},
		{
			name: "Empty",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1",
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "address"})
				mock.ExpectQuery("SELECT (.+) FROM Contacts").WillReturnRows(rows)
			},
			want: &v1.ReadAllResponse{
				Api:     "v1",
				Contact: []*v1.Contact{},
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1",
				},
			},
			mock:    func() {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.ReadAll(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("contactServiceServer.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("contactServiceServer.ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
