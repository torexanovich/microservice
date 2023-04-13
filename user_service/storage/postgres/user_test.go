package postgres

import (
	"reflect"
	"testing"

	pb "gitlab.com/micro/user_service/genproto/user"
)

func TestUserRepo_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   pb.UserRequest
		want    pb.UserResponse
		wantErr bool
	}{
		{
			name: "succes",
			input: pb.UserRequest{
				FirstName: "Amirkhan",
				LastName:  "Kenesov",
				Email:     "torexanovich.l@gmail.com",
			},
			want: pb.UserResponse{
				Id:        "1",
				FirstName: "Amirkhan",
				LastName:  "Kenesov",
				Email:     "torexanovich.l@gmail.com",
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pgRepo.CreateUser(&tc.input)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.wantErr, err)
			}

			if !reflect.DeepEqual(tc.want, *got) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
		})
	}
}
